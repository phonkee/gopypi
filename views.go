package core

import (
	"net/http"

	"io"
	"os"
	"path/filepath"

	"strings"

	"io/ioutil"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/phonkee/go-response"
	"github.com/phonkee/go-classy"
	"gopkg.in/h2non/filetype.v0"
)

/*
TemplateView can be directly instantiated when registering view.
Be careful that template view needs Config and TemplateName.

Example how ::

	classy.New(TemplateView{Config: cfg, TemplateName: "homepage.tpl.html"}).Register(router, "homepage")

*/
type TemplateView struct {
	classy.GenericView

	// configuration
	Config Config

	// Data for template
	Data map[string]interface{}

	// TemplateFilename
	TemplateName string
}

/*
GET handles return rendered template
*/
func (t TemplateView) GET(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{}
	if t.Data != nil {
		data = t.Data
	}

	var (
		err    error
		result string
	)

	if result, err = t.Config.RenderTemplate(data, t.TemplateName, t.TemplateName); err != nil {
		panic(err)
	}

	// Write html
	response.New().HTML(result).Write(w, r)
}

/*
PostPackageView handles all methods on package
*/
type PostPackageView struct {
	classy.GenericView

	Config Config
}

/*
Before every request we need to check permissions
*/
func (p *PostPackageView) Before(w http.ResponseWriter, r *http.Request) (resp response.Response) {
	if err := r.ParseMultipartForm(2 << 32); err != nil {
		resp = response.New(http.StatusBadRequest).Error(err)
	}
	return
}

/*
POST method for handling package
*/
func (p *PostPackageView) POST(w http.ResponseWriter, r *http.Request) response.Response {

	var (
		err error
	)

	resp := response.New(http.StatusInternalServerError)

	var action int

	// get action from request
	if action, err = GetPostAction(r); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	// call by action
	switch action {
	case POST_PACKAGE_ACTION_FILE_UPLOAD:
		return p.ActionFileUpload(r)
	case POST_PACKAGE_ACTION_SUBMIT:
		return response.New(http.StatusNotAcceptable)
	}

	// return response
	return resp
}

/*
ActionFileUpload Handles file upload which is called with

	`python setup.py upload`
*/
func (p *PostPackageView) ActionFileUpload(r *http.Request) response.Response {

	var (
		err  error
		pack Package
	)

	// get package by request
	if pack, err = GetPostedPackage(p.Config, r); err != nil {
		return response.New(http.StatusBadRequest).Error(err)
	}

	var (
		user User
	)

	// get user from context
	if user, err = ContextGetTokenUser(r.Context()); err != nil {
		return response.Error(err)
	}

	// if package is newly created, check permissions
	if p.Config.DB().NewRecord(pack) {

		// check if user can create new package
		if !user.CanCreate {
			return response.New(http.StatusForbidden)
		}

		// try to save package to database
		if err = p.Config.DB().Create(&pack).Error; err != nil {
			return response.Error(err)
		}
	} else {
		// check if user is maintainer or author
		if !(pack.AuthorID == user.ID || (p.Config.Manager().Package().IsMaintainer(&pack, &user) && user.CanUpdate)) {
			return response.New(http.StatusForbidden).Error("You are not maintainer")
		}
	}

	var (
		pv PackageVersion
	)

	// get package version
	if pv, err = GetPostedPackageVersion(p.Config, pack, r); err != nil {
		return response.Error(err)
	}

	// new package version / store it to database
	if p.Config.DB().NewRecord(pv) {
		pv.Author = &user
		if err = p.Config.DB().Create(&pv).Error; err != nil {
			return response.Error(err)
		}

		// udpate versions order
		if err := p.Config.Manager().Package().UpdateVersionOrder(pack); err != nil {
			return response.Error(err)
		}
	}

	var (
		f   string
		pvf PackageVersionFile
	)

	// get file
	if pvf, f, err = GetPostedPackageVersionFile(p.Config, pv, r); err != nil {
		return response.Error(err)
	}

	// file already exists in database
	if !p.Config.DB().NewRecord(pvf) {
		return response.New(http.StatusForbidden)
	}

	// assign author to file
	pvf.Author = &user

	abspath := filepath.Join(p.Config.Packages().Directory(), pvf.RelativePath)
	fullfilename := filepath.Join(abspath, pvf.Filename)

	if err = os.MkdirAll(abspath, 0777); err != nil {
		return response.Error(err)
	}

	var (
		sourcefile *os.File
		targetfile *os.File
	)

	// open source file
	if sourcefile, err = os.Open(f); err != nil {
		return response.Error(err)
	}

	// first close and then remove
	defer os.Remove(f)
	defer sourcefile.Close()

	if targetfile, err = os.Create(fullfilename); err != nil {
		return response.Error(err)
	}

	// cleanup
	defer targetfile.Close()

	// copy from file to target
	if _, err = io.Copy(targetfile, sourcefile); err != nil {
		return response.Error(err)
	}

	// flush contents
	if err = targetfile.Sync(); err != nil {
		return response.Error(err)
	}

	// create PackageVersionFile
	if err = p.Config.DB().Create(&pvf).Error; err != nil {
		return response.Error(err)
	}

	return response.OK()
}

/*
PackageListView returns list of packages
*/
type PackageListView struct {
	classy.ListView

	// config instance
	Config Config
}

/*
List (http GET) returns list of all packages

if `format` url query is set to json, json response will be returned
*/
func (p *PackageListView) List(rw http.ResponseWriter, r *http.Request) response.Response {
	var (
		err  error
		list []Package
	)

	// list all packages
	if err = p.Config.DB().Order("name").Find(&list).Error; err != nil {
		return response.Error(err)
	}

	// handle ?format=json
	if r.URL.Query().Get("format") == "json" {
		return response.OK().SliceResult(list)
	}

	data := map[string]interface{}{
		"Packages": list,
	}

	var rendered string

	if rendered, err = p.Config.RenderTemplate(data, "index", "package_list.tpl.html"); err != nil {
		return response.Error(err)
	}

	return response.OK().HTML(rendered)
}

/*
PackageDetailView returns detail of package
*/
type PackageDetailView struct {
	classy.SlugDetailView

	// config instance
	Config Config
}

/*
Retrieve is GET method
*/
func (p *PackageDetailView) Retrieve(w http.ResponseWriter, r *http.Request) response.Response {
	slug := mux.Vars(r)["slug"]

	pack := &Package{}

	if p.Config.DB().Where("name = ?", slug).Preload("Versions").Preload("Versions.License").First(pack).RecordNotFound() {
		return response.New(http.StatusNotFound)
	}

	return response.OK().Result(pack)

}

/*
PackageDownloadView serves download
*/
type PackageDownloadView struct {
	classy.BaseView

	Config Config
}

/*
Routes returns list of routes with predefined method maps
*/
func (p *PackageDownloadView) Routes() (result map[string]classy.Mapping) {
	result = map[string]classy.Mapping{
		"/{filename:.+}": classy.NewMapping(
			[]string{"GET", "Download"},
		),
	}
	return
}

/*
Download returns content of requested file.

Aside of that, when download_stats feature is enabled, stats will be recorded to database.
 */
func (p *PackageDownloadView) Download(w http.ResponseWriter, r *http.Request) response.Response {

	filename := mux.Vars(r)["filename"]
	splitted := strings.Split(filename, "/")

	if len(splitted) == 1 {
		return response.BadRequest()
	}

	final := splitted[len(splitted)-1]
	path := strings.Join(splitted[:len(splitted)-1], "/")

	pvf := PackageVersionFile{}

	// get file
	if p.Config.DB().Where("filename = ? AND relative_path = ?", final, path).First(&pvf).RecordNotFound() {
		return response.NotFound()
	}

	// return file
	absfilename := p.Config.Manager().PackageVersionFile().GetAbsoluteFilename(&pvf)

	// prepare response
	result := response.OK()

	if typ, err := filetype.MatchFile(absfilename); err != nil {
		return response.Error(err)
	} else {
		result.Header("Content-Type", typ.MIME.Value)
	}

	result.Header("Content-Disposition", "attachment; filename="+final)

	// read file contents
	if body, err := ioutil.ReadFile(absfilename); err != nil {
		return response.Error(err)
	} else {
		result.Header("Content-Length", strconv.Itoa(len(body)))
		result.Body(body)
	}

	// Add download
	p.Config.Manager().DownloadStats().AddDownloadFile(&pvf)

	return result
}
