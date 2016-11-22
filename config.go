package core

import (
	"github.com/uber-go/zap"

	"os"
	"path"

	"github.com/phonkee/gopypi/templates"

	"net/url"

	"bytes"

	"io/ioutil"

	"fmt"

	gbht "github.com/arschles/go-bindata-html-template"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	toml "github.com/pelletier/go-toml"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

/*
Config interface
*/
type Config interface {

	// Core returns CoreConfig
	Core() CoreConfig

	// getter for database connection
	DB() *gorm.DB

	// logger instance
	Logger() zap.Logger

	// validates configuration
	Validate() error

	// parse html template
	ParseTemplate(name, filename string) (*gbht.Template, error)

	// parse multiple templates
	ParseTemplateFiles(name string, filenames ...string) (*gbht.Template, error)

	// render template with given data
	RenderTemplate(data interface{}, name, filename string) (string, error)

	// render templates with given data
	RenderTemplateFiles(data interface{}, name string, filenames ...string) (string, error)

	// return router
	Router() *mux.Router

	// DownloadStats return config for download stats
	DownloadStats() DownloadStatsConfig

	// packages returns configuration for packages
	Packages() PackagesConfig

	// Manager returns interface that supplies multiple db managers
	Manager(tx ...*gorm.DB) ManagerConfig
}

type CoreConfig interface {
	// Listen returns listen address with port
	Listen() string

	// Host returns host when defined
	Host() string

	// SecretKey returns secret key for hashing and crypto
	SecretKey() string
}

type DownloadStatsConfig interface {

	// Returns how many weeks we should store weekly statistics
	ArchiveWeekly() int

	// Returns how many months we should store monthly statistics
	ArchiveMonthly() int
}

type ManagerConfig interface {
	// ClassifierManager returns new ClassifierManager instance
	Classifier(tx ...*gorm.DB) *ClassifierManager

	// DownloadStatsManager returns new DownloadStatsManager instance
	DownloadStats(tx ...*gorm.DB) *DownloadStatsManager

	// FeatureManager
	Feature(tx ...*gorm.DB) *FeatureManager

	// LicenseManager returns new LicenseManager instance
	License(tx ...*gorm.DB) *LicenseManager

	// PackageManager returns new PackageManager instance
	Package(tx ...*gorm.DB) *PackageManager

	// PackageVersionManager returns new PackageVersionManager instance
	PackageVersion(tx ...*gorm.DB) *PackageVersionManager

	// PackageVersionFileManager
	PackageVersionFile(tx ...*gorm.DB) *PackageVersionFileManager

	// PlatformManager returns new PlatformManager instance
	Platform(tx ...*gorm.DB) *PlatformManager

	// UserManager returns UserManager instance to query user data
	User(tx ...*gorm.DB) *UserManager
}

type PackagesConfig interface {
	// Directory returns full directory containint all packages
	Directory() string

	// JoinDirectory joins given parts with directory
	JoinDirectory(part ...string) string
}

/*
NewConfigFromFilename returns configuration from filename
*/
func NewConfigFromFilename(filename string) (result Config, err error) {
	var content []byte

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	var tree *toml.TomlTree
	if tree, err = toml.LoadReader(bytes.NewReader(content)); err != nil {
		return
	}

	result, err = NewConfig(tree)

	return
}

/*
Return config from string
*/
func NewConfig(tree *toml.TomlTree) (result Config, err error) {
	driver := tree.GetDefault("database.driver", DEFAULT_DB_DRIVER).(string)

	// check if database driver exists
	if !StringListContains(AVAILABLE_DB_DRIVERS, driver) {
		return nil, ErrUnknownDBDriver
	}

	dsn := tree.GetDefault("database.dsn", "gopypi:gopypy@/gopypi").(string)
	packagesDir := tree.GetDefault("packages.directory", ".packages").(string)
	secret := tree.GetDefault("core.secret_key", "").(string)
	listen := tree.GetDefault("core.listen", "0.0.0.0:9700").(string)
	host := tree.GetDefault("core.host", fmt.Sprintf("http://%v", listen)).(string)

	if !path.IsAbs(packagesDir) {
		cwd, _ := os.Getwd()
		packagesDir = path.Join(cwd, packagesDir)
	}

	var db *gorm.DB
	if db, err = gorm.Open(driver, dsn); err != nil {
		return
	}

	// configure gorm
	db.SingularTable(true)

	// setup database
	setupDB(db)

	// setup logging
	//db.LogMode(true)

	router := mux.NewRouter().StrictSlash(true)

	dsc := &downloadStatsConfig{
		archiveWeekly:  tree.GetDefault("download_stats.archive_weekly", 4).(int),
		archiveMonthly: tree.GetDefault("download_stats.archive_monthly", 4).(int),
	}

	// config implementation
	result = &config{
		db:          db,
		dsc:         dsc,
		host:        host,
		listen:      listen,
		logger:      zap.New(zap.NewTextEncoder(), zap.DebugLevel),
		packagesDir: packagesDir,
		secret:      secret,
		router:      router,
		tplasset:    templates.Asset,
		funcmap: gbht.FuncMap{
			// add url reverse functionality
			"reverse": func(name string, pairs ...string) (result string) {
				var (
					err error
					url *url.URL
				)
				if url, err = router.Get(name).URL(pairs...); err != nil {
					return ""
				}
				return url.String()
			},
		},
	}

	return
}

/*
Config implementation
*/
type config struct {
	logger      zap.Logger
	db          *gorm.DB
	funcmap     gbht.FuncMap
	packagesDir string
	listen      string
	host        string
	router      *mux.Router
	secret      string
	tplasset    func(name string) ([]byte, error)
	dsc         *downloadStatsConfig
}

func (c *config) Core() CoreConfig {
	return &coreconfig{
		config: c,
	}
}

/*
Logger returns zap.Logger instance
*/
func (c *config) DB() *gorm.DB {
	return c.db
}

/*
Logger returns zap.Logger instance
*/
func (c *config) Logger() zap.Logger {
	return c.logger
}

func (c *config) ParseTemplateFiles(name string, filenames ...string) (*gbht.Template, error) {
	tpl := gbht.New(name, c.tplasset)
	tpl.Funcs(c.funcmap)
	return tpl.ParseFiles(filenames...)
}

func (c *config) ParseTemplate(name string, filename string) (*gbht.Template, error) {
	tpl := gbht.New(name, c.tplasset)
	tpl.Funcs(c.funcmap)
	return tpl.Parse(filename)
}

func (c *config) RenderTemplateFiles(data interface{}, name string, filenames ...string) (result string, err error) {
	var t *gbht.Template

	if t, err = c.ParseTemplateFiles(name, filenames...); err != nil {
		return
	}

	var b bytes.Buffer
	if err = t.Execute(&b, data); err != nil {
		return
	}

	result = b.String()
	return
}

/*
RenderTemplate renders template with given data
*/
func (c *config) RenderTemplate(data interface{}, name string, filename string) (result string, err error) {
	var t *gbht.Template

	if t, err = c.ParseTemplate(name, filename); err != nil {
		return
	}

	var b bytes.Buffer
	if err = t.Execute(&b, data); err != nil {
		return
	}

	result = b.String()
	return
}

/*
Router returns instantiated router instance
*/
func (c *config) Router() *mux.Router {
	return c.router
}

/*
SecretKey returns secret key for hashing
*/
func (c *config) SecretKey() string {
	return c.secret
}

/*
Validate validates configuration
*/
func (c *config) Validate() (err error) {
	return err
}

/*
Packages returns PackagesConfig
*/
func (c *config) Packages() PackagesConfig {
	return &packagesconfig{
		config: c,
	}
}

/*
DownloadStats returns download stats configuration
*/
func (c *config) DownloadStats() DownloadStatsConfig {
	return c.dsc
}

func (c *config) Manager(tx ...*gorm.DB) ManagerConfig {
	db := c.DB()
	if len(tx) > 0 {
		db = tx[0]
	}
	return &managerconfig{
		config: c,
		db:     db,
	}
}

type managerconfig struct {
	config *config
	db     *gorm.DB
}

/*
getDB returns db connection
*/
func (m *managerconfig) getDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 {
		return tx[0]
	} else {
		if m.db != nil {
			return m.db
		}
		return m.config.DB()
	}
}

/*
Classifier returns ClassifierManager instance
*/
func (m *managerconfig) Classifier(tx ...*gorm.DB) *ClassifierManager {
	return &ClassifierManager{DB: m.getDB(tx...)}
}

/*
Feature returns FeatureManager instance
*/
func (m *managerconfig) Feature(tx ...*gorm.DB) *FeatureManager {
	return &FeatureManager{
		&Manager{
			DB: m.getDB(tx...),
		},
	}
}

/*
License returns LicenseManager instance
*/
func (m *managerconfig) License(tx ...*gorm.DB) *LicenseManager {
	return &LicenseManager{DB: m.getDB(tx...)}
}

/*
Package returns PackageManager instance
*/
func (m *managerconfig) Package(tx ...*gorm.DB) *PackageManager {
	return &PackageManager{DB: m.getDB(tx...)}
}

/*
PackageVersion returns PackageVersionManager
*/
func (m *managerconfig) PackageVersion(tx ...*gorm.DB) *PackageVersionManager {
	return &PackageVersionManager{
		&Manager{
			DB: m.getDB(tx...),
		},
	}
}

/*
PackageVersionFile returns PackageVersionFileManager
*/
func (m *managerconfig) PackageVersionFile(tx ...*gorm.DB) *PackageVersionFileManager {
	return &PackageVersionFileManager{
		Manager: &Manager{
			DB: m.getDB(tx...),
		},
		PackagesDir: m.config.packagesDir,
		Router:      m.config.router,
	}
}

/*
Platform return new PlatformManager instance
*/
func (m *managerconfig) Platform(tx ...*gorm.DB) *PlatformManager {
	return &PlatformManager{DB: m.getDB(tx...)}
}

/*
DownloadStats returns new DownloadStatsManager instance
*/
func (m *managerconfig) DownloadStats(tx ...*gorm.DB) *DownloadStatsManager {
	return &DownloadStatsManager{
		DB:     m.getDB(tx...),
		Config: m.config,
	}
}

/*
User returns UserManager instance
*/
func (m *managerconfig) User(tx ...*gorm.DB) *UserManager {
	return &UserManager{
		DB:        m.getDB(tx...),
		SecretKey: m.config.SecretKey(),
	}
}

/*
downloadStatsConfig
*/
type downloadStatsConfig struct {
	archiveWeekly  int
	archiveMonthly int
}

func (d *downloadStatsConfig) ArchiveWeekly() int {
	return d.archiveWeekly
}

func (d *downloadStatsConfig) ArchiveMonthly() int {
	return d.archiveMonthly
}

type coreconfig struct {
	config *config
}

func (c *coreconfig) Listen() string {
	return c.config.listen
}

func (c *coreconfig) Host() string {
	return c.config.host
}

func (c *coreconfig) SecretKey() string {
	return c.config.secret
}

type packagesconfig struct {
	config *config
}

func (p *packagesconfig) Directory() string {
	return p.config.packagesDir
}

func (p *packagesconfig) JoinDirectory(parts ...string) string {
	all := []string{}
	all = append(all, p.config.packagesDir)
	all = append(all, parts...)
	return path.Join(all...)
}
