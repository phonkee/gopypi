package core

import (
	"net/http"

	"strconv"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/phonkee/go-response"
	"github.com/phonkee/go-classy"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/phonkee/go-metadata"
)

/*
FeatureAPIViewSet provides rest endpoints for features
*/
type FeatureAPIViewSet struct {
	classy.SlugViewSet

	// config instance to access managers, etc..
	Config Config
}

/*
List returns list of all features
*/
func (f *FeatureAPIViewSet) List(w http.ResponseWriter, r *http.Request) response.Response {
	features := []Feature{}

	if err := f.Config.Manager().Feature().List(&features, FFOrderBy("id ASC")).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return response.Error(err)
		}
		err = nil
	}

	return response.OK().SliceResult(features)
}

/*
Retrieve retrieves single feature from database
*/
func (f *FeatureAPIViewSet) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {
	slug := mux.Vars(r)["slug"]

	feature := Feature{}
	if err := f.Config.Manager().Feature().Get(&feature, FFWhere("id = ?", slug)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
		return response.Error(err)
	}

	return response.Result(feature)
}

/*
Update updates given feature (just value)
*/
func (f *FeatureAPIViewSet) Update(w http.ResponseWriter, r *http.Request) response.Response {

	serializer := FeatureSerializer{}

	if err := Bind(r, &serializer); err != nil {
		return response.BadRequest().Error(err)
	}

	slug := mux.Vars(r)["slug"]

	feature := Feature{}
	if err := f.Config.Manager().Feature().Get(&feature, FFWhere("id = ?", slug)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
		return response.Error(err)
	}

	// update value
	feature.Value = serializer.Value
	if err := f.Config.DB().Save(&feature).Error; err != nil {
		return response.Error(err)
	}

	return response.OK()
}

/*
InfoAPIView serves GET request and returns information about gopypi server
*/
type InfoAPIView struct {
	classy.GenericView

	// store config instance
	Config Config
}

/*
GET method returns information about gopypi such as version, features, system info
*/
func (i *InfoAPIView) GET(w http.ResponseWriter, r *http.Request) response.Response {

	features := []Feature{}

	if err := i.Config.Manager().Feature().List(&features, FFOrderBy("id ASC")).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return response.Error(err)
		}
	}

	return response.New().Result(map[string]interface{}{
		"version":  VERSION,
		"features": features,
		"system":   stats_api.GetStats(),
	})
}

/*
LicenseAPIViewSet provides rest apis for handling of license model
*/
type LicenseAPIViewSet struct {
	classy.ViewSet

	// configuration
	Config Config
}

/*
Metadata for list endpoints
*/
func (l *LicenseAPIViewSet) MetadataList(w http.ResponseWriter, r *http.Request) response.Response {

	// enable debug
	metadata.Debug()

	md := metadata.New().Name("License list endpoint")
	md.Action(metadata.ACTION_RETRIEVE).Field("result").From([]License{})

	return response.OK().Result(md)
}

/*
List all licenses
*/
func (l *LicenseAPIViewSet) List(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err error
	)
	list := []License{}
	if err = l.Config.DB().Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return response.Error(err)
	}

	return response.OK().SliceResult(list)
}

/*
Retrieve single license from database
*/
func (l *LicenseAPIViewSet) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {
	var (
		err error
	)
	license := License{}
	if err = l.Config.DB().First(&license, "id = ?", mux.Vars(r)["pk"]).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		} else {
			return response.Error(err)
		}
	}

	return response.OK().Result(license)
}

/*
Update updates license with serializer data
*/
func (l *LicenseAPIViewSet) Update(w http.ResponseWriter, r *http.Request) response.Response {
	serializer := LicenseUpdateSerializer{}

	var (
		err error
		pk  int
	)

	// get primary key
	if pk, err = strconv.Atoi(mux.Vars(r)["pk"]); err != nil {
		return response.BadRequest().Error(err)
	}

	license := License{}
	if err = l.Config.DB().First(license, "id = ?", pk).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		} else {
			return response.Error(err)
		}
	}

	if vr := serializer.Validate(l.Config, pk); !vr.IsValid() {
		return response.BadRequest().Error(vr)
	}

	// update license with data from serializer
	serializer.UpdateLicense(&license)

	if err = l.Config.DB().Save(&license).Error; err != nil {
		return response.Error(err)
	}

	return response.OK().Result(license)

}

/*
LoginAPIView servers single POST method
*/
type LoginAPIView struct {
	classy.GenericView

	Config Config
}

/*
POST checks for username and password in JSON format and returns appropriate json response.
If succeeds, Authorization header is added with correct token
*/
func (l *LoginAPIView) POST(w http.ResponseWriter, r *http.Request) response.Response {
	ser := LoginSerializer{}

	// bind request to serializerif it fails raise error
	if err := Bind(r, &ser); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	// perform validate on serializer
	if err := (&ser).Validate(); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	user := User{}

	db := l.Config.DB()

	// find user by username
	if db.Where("username = ?", ser.Username).First(&user).RecordNotFound() {
		return response.NotFound().Error("user with given username and password not found")
	}

	// verify password
	if !l.Config.Manager().User().VerifyPassword(user, ser.Password) {
		return response.NotFound().Error("user with given username and password not found")
	}

	var (
		err   error
		token string
	)

	// create token
	if token, err = CreateToken(db, user, l.Config.Core().SecretKey(), TOKEN_EXPIRATION); err != nil {
		return response.New(http.StatusInternalServerError).Error(err)
	}

	return response.New().Header("Authorization", fmt.Sprintf("Bearer %s", token))
}

/*
MeAPIView gives information about currently logged in user

GET method returns information about user from token
POST method updates information

*/
type MeAPIView struct {
	classy.GenericView

	// store config
	Config Config
}

/*
GET returns information about user (from auth token)
*/
func (m *MeAPIView) GET(w http.ResponseWriter, r *http.Request) response.Response {
	var (
		err  error
		user User
	)

	// Get user from context.
	if user, err = ContextGetTokenUser(r.Context()); err != nil {
		return response.Error(err)
	}

	return response.Result(user)
}

/*
POST updates user information
*/
func (m *MeAPIView) POST(w http.ResponseWriter, r *http.Request) response.Response {

	s := UserProfileSerializer{}

	var (
		err error
	)
	if err = Bind(r, &s); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	var (
		user User
	)

	// get user
	if user, err = ContextGetTokenUser(r.Context()); err != nil {
		return response.Error(err)
	}

	// validate serializer
	vr := s.Validate(m.Config)
	if !vr.IsValid() {
		return response.New(http.StatusBadRequest).Error(vr)
	}

	user.FirstName = s.FirstName
	user.LastName = s.LastName
	user.Email = s.Email

	// try to save to database
	if err = m.Config.DB().Save(&user).Error; err != nil {
		return response.Error(err)
	}

	return response.Result(user)
}

/*
MeChangePasswordAPIView provides api to change password for currently logged user
 */
type MeChangePasswordAPIView struct {
	classy.GenericView

	// store config
	Config Config
}

/*
POST method updates password for currently logged user
 */
func (m *MeChangePasswordAPIView ) POST(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err  error
		user User
	)

	// Get user from context.
	if user, err = ContextGetTokenUser(r.Context()); err != nil {
		return response.Error(err)
	}

	s := UserChangePasswordSerializer{}
	if err = Bind(r, &s); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	s.User = user
	if vr := s.Validate(m.Config); !vr.IsValid() {
		return response.New(http.StatusBadRequest).Error(vr)
	}

	// change password
	if err = s.ChangePassword(m.Config, &user); err != nil {
		return response.Error(err)
	}

	// save user
	if err = m.Config.DB().Save(&user).Error; err != nil {
		return response.Error(err)
	}

	return response.OK()
}



/*
MyPackageAPIView gives information about packages for currently logged in user
*/
type MyPackageAPIView struct {
	classy.ListView

	// store config
	Config Config
}

/*
List returns list of packages for given logged user. User must be either author or maintainer of packege to be returned
in list.
 */
func (m *MyPackageAPIView) List(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err  error
		user User
	)

	// Get user from context.
	if user, err = ContextGetTokenUser(r.Context()); err != nil {
		return response.Error(err)
	}

	packages := []Package{}

	if err = m.Config.Manager().Package().List(&packages, FFPackagesFor(user), FFPreload("Author", "Maintainers", "Versions", "Versions.Author")).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return response.Error(err)
		}
		err = nil
	}

	return response.OK().SliceResult(packages)
}

/*
PackageAPIViewSet provides following methods

List - list packages
Retrieve - retrieve single package
*/
type PackageAPIViewSet struct {
	classy.ViewSet

	// store config
	Config Config
}

/*
List retrieves list of packages
*/
func (p *PackageAPIViewSet) List(w http.ResponseWriter, r *http.Request) response.Response {

	db := p.Config.DB()

	// don't forget to parse form
	r.ParseForm()

	paginator := CommonPaginator(r.Form)

	limit, offset := paginator.GetLimitOffset()

	packages := []Package{}
	queryset := db.
		Limit(limit).
		Offset(offset).
		Preload("Versions").
		Preload("Versions.Files").
		Preload("Versions.Files.Author").
		Preload("Author").
		Preload("Maintainers").
		Find(&packages)

	// find all packages
	if err := queryset.Error; err != nil {
		return response.Error(err)
	}

	// add DownloadURL to all files
	for i, pack := range packages {
		for j, version := range pack.Versions {
			for k, vfile := range version.Files {
				packages[i].Versions[j].Files[k].DownloadURL = p.Config.Manager().PackageVersionFile().GetDownloadURL(&vfile)
			}
		}
	}

	// write packages as response
	return response.SliceResult(packages).Data("paginator", paginator)
}

/*
Retrieve returns single package
*/
func (p *PackageAPIViewSet) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {
	pack := Package{
		ID: Atoui(mux.Vars(r)["pk"]),
	}

	// find single package
	preload := FFPreload("Author", "Versions", "Versions.Files", "Versions.Files.Author", "Versions.Author", "Maintainers")
	if p.Config.Manager().Package().Get(&pack, preload).RecordNotFound() {
		return response.New(http.StatusNotFound)
	}

	// add DownloadURL to all files
	for i, version := range pack.Versions {
		for j, vfile := range version.Files {
			pack.Versions[i].Files[j].DownloadURL = p.Config.Manager().PackageVersionFile().GetDownloadURL(&vfile)
		}
	}

	// write packages as response
	return response.Result(pack)
}

/*
StatsAPIView returns some statistic information for admin dashboard.
 */
type StatsAPIView struct {
	classy.GenericView

	Config Config
}

/*
Stats result
*/
type Stats struct {
	Packages    int `json:"packages"`
	ActiveUsers int `json:"active_users"`
	Licenses    int `json:"licenses"`
	Downloads   int `json:"downloads"`
}

/*
GET returns summary statistics about gopypi
 */
func (s *StatsAPIView) GET(w http.ResponseWriter, r *http.Request) response.Response {
	db := s.Config.DB()
	// create stats
	stats := Stats{}

	// add packages count
	db.Model(Package{}).Count(&(stats.Packages))

	// add count of all active users
	db.Model(User{}).Where("is_active = ?", true).Count(&(stats.ActiveUsers))

	// add count of opensource licenses in system
	db.Model(License{}).Count(&(stats.Licenses))

	// return count of all downloaded files
	s.Config.Manager().DownloadStats().GetCount(&(stats.Downloads))
	return response.New().Result(stats)
}

/*
StatsDownloadAllAPIView provides stats of package downloads
*/
type StatsDownloadAllAPIView struct {
	classy.ListView

	// config instance
	Config Config
}

/*
List returns download stats about all packages (sums)
*/
func (s *StatsDownloadAllAPIView) List(w http.ResponseWriter, r *http.Request) response.Response {

	var err error

	stats := map[string][]StatsDownloadItem{}

	// get all stats
	if err = s.Config.Manager().DownloadStats().GetAllStats(stats); err != nil {
		return response.Error(err)
	}

	return response.OK().Result(stats)
}

/*
StatsDownloadPackageAPIView provides stats of package downloads
*/
type StatsDownloadPackageAPIView struct {
	classy.DetailView

	// config instance
	Config Config
}

/*
List returns download stats about all packages (sum)
*/
func (s *StatsDownloadPackageAPIView) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err  error
		pack Package
	)

	// get package from url var
	if err = s.Config.Manager().Package().Get(&pack, FFID(Atoui(mux.Vars(r)["pk"]))).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		} else {
			return response.Error(err)
		}
	}

	// get all stats
	stats := map[string][]StatsDownloadItem{}
	if err = s.Config.Manager().DownloadStats().GetAllStats(stats, FFDownloadStatsPackage(&pack)); err != nil {
		return response.Error(err)
	}

	return response.OK().Result(stats)
}

/*
StatsDownloadPackageVersionAPIView provides stats of package version downloads
*/
type StatsDownloadPackageVersionAPIView struct {
	classy.DetailView

	// config instance
	Config Config
}

/*
List returns download stats about all package versions (sum)
*/
func (s *StatsDownloadPackageVersionAPIView) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err error
	)

	vars := mux.Vars(r)

	pv := PackageVersion{
		ID:        Atoui(vars["pk"]),
		PackageID: Atoui(vars["package_pk"]),
	}

	// get package from url var
	if err = s.Config.Manager().PackageVersion().Get(&pv, FFWhere(pv)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		} else {
			return response.Error(err)
		}
	}

	// get all stats
	stats := map[string][]StatsDownloadItem{}
	if err = s.Config.Manager().DownloadStats().GetAllStats(stats, FFDownloadStatsPackageVersion(&pv)); err != nil {
		return response.Error(err)
	}

	return response.OK().Result(stats)
}

/*
UserAPIViewSet handles basic crud on users
*/
type UserAPIViewSet struct {
	classy.ViewSet

	// store config
	Config Config
}

/*
List of users
*/
func (u *UserAPIViewSet) List(w http.ResponseWriter, r *http.Request) response.Response {
	db := u.Config.DB()

	// don't forget to parse form
	r.ParseForm()
	paginator := CommonPaginator(r.Form)

	users := []User{}

	// Limit queryset with applied filter
	queryset := LimitQueryset(NewUserListFilter(r).Apply(db), paginator).Find(&users)

	// find all users
	if err := queryset.Error; err != nil {
		return response.Error(err)
	}

	// set count
	CountQueryset(queryset, paginator)

	// write users as response
	return response.OK().SliceResult(users).Data("paginator", paginator)
}

/*
Retrieve single user
*/
func (u *UserAPIViewSet) Retrieve(w http.ResponseWriter, r *http.Request) {
	db := u.Config.DB()
	pk := mux.Vars(r)["pk"]

	// prepare blank user
	user := User{}

	// fetch user from database
	if db.Where("id = ?", pk).First(&user).RecordNotFound() {
		response.New(http.StatusNotFound).Write(w, r)
		return
	}

	// write users as response
	response.New().Result(user).Write(w, r)
}

/*
Create is called when new user is created
*/
func (u *UserAPIViewSet) Create(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err error
	)
	s := &UserAddSerializer{}
	if err = Bind(r, s); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	var (
		vr ValidationResult
	)

	// validate serializer
	if vr = s.Validate(u.Config); !vr.IsValid() {
		return response.New(http.StatusBadRequest).Error(vr)
	}

	// get user from serializer
	user := s.GetUser(u.Config)

	// insert user to database
	if err = u.Config.DB().Create(&user).Error; err != nil {
		return response.Error(err)
	}

	return response.Result(user)
}

/*
Update is called when POST is called to update user
*/
func (u *UserAPIViewSet) Update(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err error
		pk  int
	)
	serializer := UserUpdateSerializer{}

	// bind request to serializer
	if err = Bind(r, &serializer); err != nil {
		return response.BadRequest()
	}

	// convert primary key to int
	if pk, err = strconv.Atoi(mux.Vars(r)["pk"]); err != nil {
		return response.BadRequest()
	}

	// prepare user
	user := User{}
	if u.Config.DB().First(&user, "id = ?", pk).RecordNotFound() {
		return response.NotFound()
	}

	// add primary key to serializer so Validate method can use it
	serializer.ID = uint(pk)

	if vr := serializer.Validate(u.Config); !vr.IsValid() {
		return response.BadRequest().Error(vr)
	}

	// Update user
	serializer.UpdateUser(u.Config, &user)

	// save user
	if err = u.Config.DB().Save(&user).Error; err != nil {
		return response.Error(err)
	}

	return response.OK().Result(user)
}

/*
PackageMaintainerAPIViewSet provides rest endpoints for maintainers of given package
*/
type PackageMaintainerAPIViewSet struct {
	classy.ViewSet

	// config instance
	Config Config
}

/*
GetPackage returns package from request
*/
func (p *PackageMaintainerAPIViewSet) GetPackage(r *http.Request) (result Package, err error) {
	result = Package{}
	err = p.Config.DB().First(&result, "id = ?", mux.Vars(r)["package_pk"]).Error
	return
}

/*
List method lists all maintainers for given package
*/
func (p *PackageMaintainerAPIViewSet) List(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err  error
		pack Package
	)
	if pack, err = p.GetPackage(r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
	}

	var (
		maintainers []User
	)

	// list all maintainers
	if err = p.Config.DB().Model(pack).Association("Maintainers").Find(&maintainers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return response.Error(err)
		}

		err = nil
	}

	return response.OK().SliceResult(maintainers)
}

/*
Update association means that when it's not present, it's created,
Update method doesn't handle request body.
*/
func (p *PackageMaintainerAPIViewSet) Update(w http.ResponseWriter, r *http.Request) response.Response {
	var (
		err  error
		pack Package
	)
	if pack, err = p.GetPackage(r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
	}

	var user User

	// try to find user by pk first
	if p.Config.DB().First(&user, "id = ?", mux.Vars(r)["pk"]).RecordNotFound() {
		return response.NotFound()
	}

	var maintainers []User

	if err = p.Config.DB().Model(&pack).Association("Maintainers").Find(&maintainers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return response.Error(err)
		}
		err = nil
	}

	// iterate over maintainers to find whether maintainer is already there
	for _, maintainer := range maintainers {
		if maintainer.ID == user.ID {
			return response.OK().Result(user)
		}
	}

	// maintainer not found, create one
	if err = p.Config.DB().Model(&pack).Association("Maintainers").Append(user).Error; err != nil {
		return response.Error(err)
	}

	return response.OK().Result(user)
}

/*
Delete removes association of maintainer to given package
*/
func (p *PackageMaintainerAPIViewSet) Delete(w http.ResponseWriter, r *http.Request) response.Response {
	var (
		err  error
		pack Package
	)
	if pack, err = p.GetPackage(r); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
	}

	var user User

	// try to find user by pk first
	if err = p.Config.DB().First(&user, "id = ?", mux.Vars(r)["pk"]).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
		return response.Error(err)
	}

	// try to delete association
	if err = p.Config.DB().Model(&pack).Association("Maintainers").Delete(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		}
		return response.Error(err)
	}

	return response.OK()
}

/*
PlatformAPIViewSet provides rest endpoints for platform (RU)
*/
type PlatformAPIViewSet struct {
	classy.ViewSet

	// configuration
	Config Config
}

/*
List returns all platform stored in database
*/
func (p *PlatformAPIViewSet) List(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err error
	)

	list := []Platform{}

	if err = p.Config.Manager().Platform().List(&list).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return response.Error(err)
		}
	}

	return response.OK().SliceResult(list)
}

/*
Retrieve returns single platform
*/
func (p *PlatformAPIViewSet) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {

	// prepare platform to find
	platform := Platform{
		ID: Atoui(mux.Vars(r)["pk"]),
	}

	// get platform
	if err := p.Config.Manager().Platform().Get(&platform).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFound()
		} else {
			return response.Error(err)
		}
	}

	return response.OK().Result(platform)
}
