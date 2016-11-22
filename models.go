package core

import (
	"time"

	"path/filepath"

	"github.com/jinzhu/gorm"
)

/*
Migrate runs AutoMigrate on all models.
*/
func Migrate(config Config) (err error) {
	db := config.DB().Debug()
	db.AutoMigrate(Package{}, PackageVersion{}, PackageVersionFile{})
	db.AutoMigrate(User{})
	db.AutoMigrate(Classifier{})
	db.AutoMigrate(License{})
	db.AutoMigrate(Platform{})
	db.AutoMigrate(DownloadStatsWeekly{}, DownloadStatsMonthly{}, DownloadStatsYearly{})
	db.AutoMigrate(Feature{})

	// create all features
	if err = createFeatures(db); err != nil {
		return
	}

	return
}

/*
Classifier model

We store just values for now, later we can do some sort of tree
*/
type Classifier struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Approved bool   `json:"approved"`
	Name     string `gorm:"unique_index" json:"name"`
}

/*
DownloadStats is common structure for all of them

CreatedAt is aligned on every struct that embeds DownloadStats with given aggregation level.
This assures that records are groupped correctly and no need to do grouping no database level.
*/
type DownloadStats struct {
	ID               uint            `gorm:"primary_key" json:"-"`
	PackageVersion   *PackageVersion `gorm:"ForeignKey:PackageVersionID" json:"version,omitempty"`
	PackageVersionID uint            `json:"-"`
	Downloads        int             `json:"downloads"`
	CreatedAt        time.Time       `json:"created_at"`
}

/*
DownloadStatsWeekly represents download stats aggregated by week
*/
type DownloadStatsWeekly struct {
	DownloadStats
}

/*
BeforeSave aligns CreatedAt correctly
*/
func (s *DownloadStatsWeekly) BeforeCreate() error {
	s.CreatedAt = TimeAlignWeek(time.Now())
	return nil
}

/*
DownloadStatsMonthly represents download stats aggregated by month
*/
type DownloadStatsMonthly struct {
	DownloadStats
}

/*
BeforeSave aligns CreatedAt correctly
*/
func (s *DownloadStatsMonthly) BeforeCreate() error {
	s.CreatedAt = TimeAlignMonth(time.Now())
	return nil
}

/*
DownloadStatsYearly represents download stats aggregated by year
*/
type DownloadStatsYearly struct {
	DownloadStats
}

/*
BeforeSave aligns CreatedAt correctly
*/
func (s *DownloadStatsYearly) BeforeCreate() error {
	s.CreatedAt = TimeAlignYear(time.Now())
	return nil
}

/*
Feature is model for enabling/disabling gopypi features
*/
type Feature struct {
	ID          string `gorm:"primary_key" json:"id"`
	Description string `json:"description"`
	Value       bool   `json:"value"`
}

/*
License model

Holds informations about available licenses.
*/
type License struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Approved bool   `json:"approved"`
	Code     string `json:"code"`
	Content  string `json:"content"`
	Name     string `json:"name"`
}

/*
Package model.
*/
type Package struct {
	ID          uint             `gorm:"primary_key" json:"id"`
	Name        string           `json:"name"`
	Versions    []PackageVersion `gorm:"ForeignKey:PackageID" json:"versions,omitempty"`
	Maintainers []User           `gorm:"many2many:package_maintainers;" json:"maintainers,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Author      *User            `gorm:"ForeignKey:AuthorID" json:"author,omitempty"`
	AuthorID    uint             `json:"-"`
}

/*
BeforeCreate sets CreatedAt
*/
func (p *Package) BeforeCreate() error {
	p.CreatedAt = gorm.NowFunc()
	return nil
}

/*
BeforeSave sets UpdatedAt
*/
func (p *Package) BeforeSave() error {
	p.UpdatedAt = gorm.NowFunc()
	return nil
}

/*
PackageVersion model that holds information about given package version
*/
type PackageVersion struct {
	ID           uint                 `gorm:"primary_key" json:"id"`
	PackageID    uint                 `json:"-"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Author       *User                `gorm:"ForeignKey:AuthorID" json:"author,omitempty"`
	AuthorID     uint                 `json:"-"`
	Comment      string               `json:"comment"`
	Description  string               `json:"description"`
	Summary      string               `json:"summary"`
	HomePage     string               `json:"home_page"`
	License      *License             `gorm:"ForeignKey:LicenseID" json:"license,omitempty"`
	LicenseID    uint                 `json:"-"`
	Version      string               `gorm:"index" json:"version"`
	VersionOrder int                  `gorm:"index" json:"version_order"`
	Files        []PackageVersionFile `gorm:"ForeignKey:PackageVersionID" json:"files,omitempty"`
	Classifiers  []Classifier         `gorm:"many2many:package_version_classifiers;" json:"classifiers,omitempty"`
}

/*
BeforeCreate sets CreatedAt
*/
func (p *PackageVersion) BeforeCreate() error {
	p.CreatedAt = gorm.NowFunc()
	return nil
}

/*
BeforeSave sets UpdatedAt
*/
func (p *PackageVersion) BeforeSave() error {
	p.UpdatedAt = gorm.NowFunc()
	return nil
}

/*
PackageVersionFile model
*/
type PackageVersionFile struct {
	ID               uint      `gorm:"primary_key" json:"id"`
	PackageVersionID uint      `json:"-"`
	Filename         string    `json:"filename"`
	RelativePath     string    `json:"relative_path"`
	MD5Digest        string    `gorm:"column:md5_digest" json:"md5_digest"`
	Author           *User     `gorm:"ForeignKey:AuthorID" json:"author,omitempty"`
	AuthorID         uint      `json:"-"`
	CreatedAt        time.Time `json:"created_at"`

	// this field is used to have pregenearated download url
	DownloadURL string `gorm:"-" json:"download_url,omitempty"`
}

/*
BeforeCreate sets CreatedAt
*/
func (p *PackageVersionFile) BeforeCreate() error {
	p.CreatedAt = gorm.NowFunc()
	return nil
}

/*
GenerateRelativePath generates random relative path
*/
func (p PackageVersionFile) GenerateRelativePath() string {
	random := GenerateSalt(32)
	hash := MD5(string(random))
	return filepath.Join(hash[:2], hash[:4], hash[4:])
}

/*
Platform model tracks all platforms such as: Linux, Darwin even more esoteric.
*/
type Platform struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
User model

To set/verify user password please use UserManager
*/
type User struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Username  string `gorm:"type:varchar(20);unique_index" json:"username"`
	Email     string `gorm:"type:varchar(100);index" json:"email"`
	Password  string `gorm:"type:varchar(256)" json:"-"`
	FirstName string `gorm:"type:varchar(40)" json:"first_name"`
	LastName  string `gorm:"type:varchar(40)" json:"last_name"`
	IsActive  bool   `json:"is_active"`
	IsAdmin   bool   `json:"is_admin"`

	// permissions
	CanList     bool `json:"can_list"`
	CanCreate   bool `json:"can_create"`
	CanDownload bool `json:"can_download"`
	CanUpdate   bool `json:"can_update"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
BeforeCreate sets CreatedAt
*/
func (u *User) BeforeCreate() error {
	u.CreatedAt = gorm.NowFunc()
	return nil
}

/*
BeforeSave sets UpdatedAt
*/
func (u *User) BeforeSave() error {
	u.UpdatedAt = gorm.NowFunc()
	return nil
}
