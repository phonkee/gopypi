/*
Managers is set of manager structs that provide querying methods
*/
package core

import (
	"encoding/base64"
	"fmt"

	"strconv"
	"strings"

	"sort"

	"path"

	"time"

	"github.com/blang/semver"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
)

/*
Generic manager

Should have following methods:

* List
* Get

*/
type Manager struct {
	DB *gorm.DB
}

/*
Get calls Where method from given pack with filter funcs applied
*/
func (p *Manager) Get(target interface{}, filter ...FilterFunc) *gorm.DB {
	queryset := p.DB.Where(target)
	queryset = ApplyFilterFuncs(queryset, filter...)
	return queryset.First(target)
}

/*
List calls Where method from given pack with filter funcs applied
*/
func (p *Manager) List(target interface{}, filter ...FilterFunc) *gorm.DB {
	queryset := p.DB.Where(target)
	queryset = ApplyFilterFuncs(queryset, filter...)
	return queryset.Find(target)
}

/*
ClassifierManager database manager
*/
type ClassifierManager struct {
	DB *gorm.DB
}

/*
ListOrCreate returns list of classifiers from given string list
*/
func (c *ClassifierManager) ListOrCreate(result *[]Classifier, cd []string) (err error) {

	for _, cls := range cd {
		classifier := Classifier{}

		// find first classifier or create one
		if err = c.DB.FirstOrCreate(&classifier, Classifier{Name: c.NormalizeName(cls)}).Error; err != nil {
			return
		}

		*result = append(*result, classifier)
	}
	return
}

/*
NormalizeName normalizes classifier name
*/
func (c *ClassifierManager) NormalizeName(name string) string {
	splitted := strings.Split(name, "::")

	parts := []string{}
	for _, part := range splitted {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, " :: ")
}

/*
FeatureManager database manager for model Feature
*/
type FeatureManager struct {
	*Manager
}

/*
IsEnabledFeature returns whether given feature is available
*/
func (f *FeatureManager) IsEnabledFeature(feature string) (result bool, err error) {
	target := Feature{}

	if err := f.Manager.DB.First(&target, "id = ?", feature).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return target.Value, nil
}

/*
LicenseManager database manager
*/
type LicenseManager struct {
	DB *gorm.DB
}

/*
GetByCode returns license by code, if not exist create new one
*/
func (l *LicenseManager) GetByCode(code string) (license License, err error) {
	code = strings.TrimSpace(code)

	if code == "" {
		err = ErrLicenseNotFound
		return
	}

	license = License{
		Code: code,
	}

	l.DB.FirstOrCreate(&license, license).Attrs(license)

	return
}

/*
PackageManager database manager
*/
type PackageManager struct {
	DB *gorm.DB
}

/*
Get calls Where method from given pack with preloads
*/
func (p *PackageManager) Get(pack *Package, filter ...FilterFunc) *gorm.DB {
	queryset := p.DB.Where(pack)
	queryset = ApplyFilterFuncs(queryset, filter...)
	return queryset.First(pack)
}

/*
Get calls Where method from given pack with preloads
*/
func (p *PackageManager) List(packages *[]Package, filter ...FilterFunc) *gorm.DB {
	queryset := p.DB.Model(Package{})
	queryset = ApplyFilterFuncs(queryset, filter...)
	return queryset.Find(packages)
}

/*
check if user is maintainer
*/
func (p *PackageManager) IsMaintainer(pack *Package, user *User) bool {

	if p.DB.NewRecord(pack) || p.DB.NewRecord(user) {
		return false
	}

	target := Package{}

	if p.DB.Preload("Maintainers").First(&target, "id = ?", pack.ID).Error != nil {
		return false
	}

	for _, maintainer := range target.Maintainers {
		if maintainer.ID == user.ID {
			return true
		}
	}

	return false
}

/*
Item for ordering Package Versions
*/
type orderItem struct {
	ID     uint
	Semver semver.Version
}

type versionOrder []orderItem

func (v versionOrder) Len() int           { return len(v) }
func (v versionOrder) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v versionOrder) Less(i, j int) bool { return v[i].Semver.LT(v[j].Semver) }

/*
UpdateVersionOrder updates order for all versions
*/
func (p *PackageManager) UpdateVersionOrder(pack Package) (err error) {

	// if not saved in database return
	if p.DB.NewRecord(pack) {
		return ErrPackageNotFound
	}

	o := versionOrder{}

	versions := []PackageVersion{}

	if err = p.DB.Find(&versions, "package_id = ?", pack.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			return
		}
	}

	var (
		errSemver error
		sv        semver.Version
	)

	for _, version := range versions {
		if sv, errSemver = semver.Make(version.Version); errSemver != nil {
			continue
		}

		o = append(o, orderItem{
			ID:     version.ID,
			Semver: sv,
		})
	}

	// sort them by semver
	sort.Sort(o)

	// update
	p.DB.Model(PackageVersion{}).Where("package_id = ?", pack.ID).Update("version_order", 0)

	for index, oi := range o {
		p.DB.Model(PackageVersion{}).Where("id = ?", oi.ID).Update("version_order", index+1)
	}

	return
}

/*
PackageVersionManager database manager
*/
type PackageVersionManager struct {
	*Manager
}

/*
PackageVersionFileManager database manager
*/
type PackageVersionFileManager struct {
	*Manager

	PackagesDir string
	Router      *mux.Router
}

/*
GetAbsoluteFilename returns full package filename
*/
func (p *PackageVersionFileManager) GetAbsoluteFilename(pvf *PackageVersionFile) string {
	return path.Join(p.PackagesDir, p.GetRelativeFilename(pvf))
}

/*
GetAbsoluteFilename returns full package filename
*/
func (p *PackageVersionFileManager) GetRelativeFilename(pvf *PackageVersionFile) string {
	return path.Join(pvf.RelativePath, pvf.Filename)
}

/*
GetDownloadURL returns full url for downloading package
*/
func (p *PackageVersionFileManager) GetDownloadURL(pvf *PackageVersionFile) string {
	url, _ := p.Router.Get("package_download").URL("filename", p.GetRelativeFilename(pvf))
	return url.String()
}

/*
PlatformManager database manager for model Platform
*/
type PlatformManager struct {
	DB *gorm.DB
}

/*
Get returns platform by set fields
*/
func (p *PlatformManager) Get(platform *Platform) *gorm.DB {
	return p.DB.FirstOrCreate(&platform, platform).Attrs(platform)
}

/*
List returns list of platforms
*/
func (p *PlatformManager) List(platforms *[]Platform, filter ...FilterFunc) *gorm.DB {
	db := ApplyFilterFuncs(p.DB, filter...)
	return db.Find(platforms)
}

/*
Return multiple platforms
*/
func (p *PlatformManager) ListOrCreate(target *[]Platform, platforms []string) error {

	for _, platform := range platforms {

		po := Platform{
			Name: platform,
		}

		if err := p.Get(&po).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			*target = append(*target, po)
		}
	}

	return nil
}

/*
DownloadStatsManager database manager for all download stats
*/
type DownloadStatsManager struct {
	DB     *gorm.DB
	Config Config
}

/*
AddDownload adds download for given package version
*/
func (s *DownloadStatsManager) AddDownload(version *PackageVersion) (err error) {

	var enabled bool

	if enabled, err = s.Config.Manager().Feature().IsEnabledFeature(FEATURE_DOWNLOAD_STATS); err != nil {
		return
	}

	// check if download stats are enabled
	if !enabled {
		return
	}

	// this is now
	now := time.Now()

	// weekly stats are stored
	if s.Config.DownloadStats().ArchiveWeekly() > 0 {
		ds := DownloadStatsWeekly{}
		ds.PackageVersion = version
		ds.CreatedAt = TimeAlignWeek(now)

		attrs := DownloadStatsWeekly{}
		attrs.PackageVersion = version
		attrs.CreatedAt = TimeAlignWeek(now)

		if err = s.DB.Where("package_version_id = ? AND created_at = ?", version.ID, ds.CreatedAt).Attrs(attrs).FirstOrCreate(&ds).Error; err != nil {
			return
		}

		// update downloads for given download stats
		if err = s.DB.Model(&ds).Updates(map[string]interface{}{"downloads": gorm.Expr("downloads + 1")}).Error; err != nil {
			return
		}
	}

	// weekly stats are stored
	if s.Config.DownloadStats().ArchiveMonthly() > 0 {
		ds := DownloadStatsMonthly{}
		ds.PackageVersion = version
		ds.CreatedAt = TimeAlignMonth(now)

		attrs := DownloadStatsMonthly{}
		attrs.PackageVersion = version
		attrs.CreatedAt = TimeAlignMonth(now)

		if err = s.DB.Where("package_version_id = ? AND created_at = ?", version.ID, ds.CreatedAt).Attrs(attrs).FirstOrCreate(&ds).Error; err != nil {
			return
		}

		// update downloads for given download stats
		if err = s.DB.Model(&ds).Updates(map[string]interface{}{"downloads": gorm.Expr("downloads + 1")}).Error; err != nil {
			return
		}
	}

	// Yearly statistics are store everytime
	ds := DownloadStatsYearly{}
	ds.PackageVersion = version
	ds.CreatedAt = TimeAlignYear(now)

	attrs := DownloadStatsYearly{}
	attrs.PackageVersion = version
	attrs.CreatedAt = TimeAlignYear(now)

	if err = s.DB.Where("package_version_id = ? AND created_at = ?", version.ID, ds.CreatedAt).Attrs(attrs).FirstOrCreate(&ds).Error; err != nil {
		return
	}

	// update downloads for given download stats
	if err = s.DB.Model(&ds).Updates(map[string]interface{}{"downloads": gorm.Expr("downloads + 1")}).Error; err != nil {
		return
	}

	return
}

/*
AddDownloadFile adds download by PackageVersionFile
*/
func (d *DownloadStatsManager) AddDownloadFile(versionfile *PackageVersionFile) (err error) {

	pv := PackageVersion{}

	if err = d.DB.First(&pv, "id = ?", versionfile.PackageVersionID).Error; err != nil {
		return err
	}

	// add download by package version
	err = d.AddDownload(&pv)

	return
}

/*
GetSum returns count of all downloads
*/
func (d *DownloadStatsManager) GetCount(target interface{}) (err error) {
	err = d.DB.Table("download_stats_yearly").Select("sum(downloads) as total").Row().Scan(target)
	return
}

/*
Returns all download stats
*/
func (d *DownloadStatsManager) GetAllStats(target map[string][]StatsDownloadItem, filter ...FilterFunc) (err error) {

	targetYearly := []StatsDownloadItem{}

	// get yearly stats
	if err = d.GetStats(STATS_DOWNLOAD_ALL, &targetYearly, filter...).Error; err != nil {
		return
	}

	// get weekly stats
	targetWeekly := []StatsDownloadItem{}
	if err = d.GetStats(STATS_DOWNLOAD_WEEKLY, &targetWeekly, filter...).Error; err != nil {
		return
	}

	// get monthly stats
	targetMonthly := []StatsDownloadItem{}
	if err = d.GetStats(STATS_DOWNLOAD_MONTHLY, &targetMonthly, filter...).Error; err != nil {
		return
	}

	target["all"] = targetYearly
	target["weekly"] = targetWeekly
	target["monthly"] = targetMonthly

	return
}

func (d *DownloadStatsManager) GetStats(aggregation StatsAggregation, target *[]StatsDownloadItem, filter ...FilterFunc) *gorm.DB {

	var m interface{}

	switch aggregation {
	case STATS_DOWNLOAD_WEEKLY:
		m = DownloadStatsWeekly{}
	case STATS_DOWNLOAD_ALL:
		m = DownloadStatsYearly{}
	case STATS_DOWNLOAD_MONTHLY:
		m = DownloadStatsMonthly{}
	default:
		panic("unknown aggregation")
	}

	// create queryset
	db := d.DB.Model(m).
		Select("created_at, sum(downloads) as downloads").
		Group("created_at").
		Order("created_at ASC")

	// apply FilterFuncs
	db = ApplyFilterFuncs(db, filter...)

	return db.Scan(target)
}

/*
Cleanup deletes weekly and monthly stats from database, yearly stats will stay forever
 */
func (d *DownloadStatsManager ) Cleanup()  {

	weekly := time.Now().Add(-(time.Hour * time.Duration(24 * 7 * d.Config.DownloadStats().ArchiveWeekly())))
	monthly := time.Now().Add(-(time.Hour * time.Duration(24 * 7 * 4 * d.Config.DownloadStats().ArchiveMonthly())))

	wd := d.Config.DB().Delete(DownloadStatsWeekly{}, "created_at < ?", weekly).RowsAffected
	println("Deleted weekly download stats records:", wd)

	md := d.Config.DB().Delete(DownloadStatsMonthly{}, "created_at < ?", monthly).RowsAffected
	println("Deleted monthly download stats records:", md)
}

/*
UserManager groups functionality to query user model instances
*/
type UserManager struct {
	DB        *gorm.DB
	SecretKey string
}

/*
ExistsUsername returns whether user with given username exists in database
*/
func (u *UserManager) ExistsUsername(username string) bool {
	return !u.DB.First(&User{}, "username = ?", username).RecordNotFound()
}

/*
ExistsEmail returns whether user with given email exists in database
*/
func (u *UserManager) ExistsEmail(email string) bool {
	return !u.DB.First(&User{}, "email = ?", email).RecordNotFound()
}

/*
Get
*/
func (u *UserManager) Get(user *User) *gorm.DB {
	return u.DB.Where(user).First(user)
}

/*
SetPassword sets password for given user
*/
func (u *UserManager) SetPassword(user *User, password string) {
	salt := GenerateSalt(PASSWORD_SALT_BYTES)
	strsalt := base64.StdEncoding.EncodeToString(salt)

	var (
		err     error
		hash    []byte
		strhash string
	)

	finalSalt := append(salt, []byte(u.SecretKey)...)

	if hash, err = scrypt.Key([]byte(password), finalSalt, PASSWORD_ITERATIONS, 8, 1, PASSWORD_HASH_BYTES); err != nil {
		panic(err)
	}

	strhash = base64.StdEncoding.EncodeToString(hash)
	user.Password = fmt.Sprintf("scrypt$%v$%v$%v", PASSWORD_ITERATIONS, strsalt, strhash)
	return
}

/*
VerifyPassword verifies password for user
*/
func (u *UserManager) VerifyPassword(user User, password string) bool {
	// unusable password
	if strings.TrimSpace(user.Password) == "" {
		return false
	}

	splitted := strings.SplitN(user.Password, "$", 4)

	// invalid password
	if len(splitted) != 4 {
		return false
	}

	// algorithm currently not used (in the future we will have probably them)
	iterations, salt, hash := splitted[1], splitted[2], splitted[3]

	var (
		err     error
		iterint int
	)
	if iterint, err = strconv.Atoi(iterations); err != nil {
		return false
	}

	saltbyte, errsalt := base64.StdEncoding.DecodeString(salt)
	if errsalt != nil {
		return false
	}

	finalSalt := append(saltbyte, []byte(u.SecretKey)...)

	one, err := scrypt.Key([]byte(password), finalSalt, iterint, 8, 1, PASSWORD_HASH_BYTES)
	if err != nil {
		return false
	}

	dest := base64.StdEncoding.EncodeToString(one)
	return dest == hash
}
