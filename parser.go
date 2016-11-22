/*
parser package handles parsing of submitted packages.
*/
package core

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

/*
Return action from form
*/
func GetPostAction(r *http.Request) (result int, err error) {
	var ok bool
	if result, ok = POST_PACKAGE_ACTIONS[r.Form.Get(":action")]; !ok {
		err = ErrPostPackageInvalidAction
	}
	return
}

/*
GetPackage parses request form and returns package
*/
func GetPostedPackage(cfg Config, r *http.Request) (pack Package, err error) {

	form := r.Form
	name := strings.TrimSpace(form.Get("name"))

	// if no name given error
	if name == "" {
		err = ErrPostPackageInvalidName
		return
	}

	// prepare package
	pack = Package{}

	var (
		user User
	)

	if cfg.DB().Preload("Author").First(&pack, "name = ?", name).RecordNotFound() {
		pack.Name = name
		// get user from context
		if user, err = ContextGetTokenUser(r.Context()); err != nil {
			return
		}

		pack.Author = &user
	}

	return
}

/*
Return package version
*/
func GetPostedPackageVersion(cfg Config, pack Package, r *http.Request) (pv PackageVersion, err error) {

	pv = PackageVersion{}

	version := strings.TrimSpace(r.Form.Get("version"))

	// check version
	if version == "" {
		err = ErrPostPackageInvalidVersion
		return
	}

	// check if package version exists in database
	if cfg.DB().First(&pv, "package_id = ? AND version = ?", pack.ID, version).RecordNotFound() {
		pv.Version = version
		pv.Comment = strings.TrimSpace(r.Form.Get("comment"))
		pv.Description = strings.TrimSpace(r.Form.Get("description"))
		pv.Summary = strings.TrimSpace(r.Form.Get("summary"))
		pv.Version = strings.TrimSpace(r.Form.Get("version"))
		pv.HomePage = strings.TrimSpace(r.Form.Get("home_page"))

		// assign package
		pv.PackageID = pack.ID

		var (
			classifiers []string
			ok          bool
		)

		// get classifiers
		if classifiers, ok = r.Form["classifiers"]; !ok {
			classifiers = []string{}
		}

		c := make([]Classifier, 0, len(classifiers))

		// get and assign classifiers
		if err = cfg.Manager().Classifier().ListOrCreate(&c, classifiers); err != nil {
			return
		}

		pv.Classifiers = c

		// get license
		if license, errLicense := cfg.Manager().License().GetByCode(strings.TrimSpace(r.Form.Get("license"))); errLicense != nil {
			if errLicense != ErrLicenseNotFound {
				err = errLicense
				return
			}
		} else {
			pv.License = &license
		}
	}

	return
}

/*
Returns PackageVersionFile along with tempfile
*/
func GetPostedPackageVersionFile(config Config, pv PackageVersion, r *http.Request) (result PackageVersionFile, filename string, err error) {
	var (
		f      *os.File
		file   multipart.File
		header *multipart.FileHeader
	)

	// parse file content
	file, header, err = r.FormFile("content")
	if err != nil {
		return
	}

	// create temporary file
	if f, err = ioutil.TempFile("", header.Filename); err != nil {
		return
	}

	// get filename
	filename = f.Name()

	var (
		content []byte
	)

	// read file content
	if content, err = ioutil.ReadAll(file); err != nil {
		return
	}

	// write contents to file
	if _, err = f.Write(content); err != nil {
		return
	}

	defer f.Close()

	// close file
	result = PackageVersionFile{}

	// check if packaged version file exists / if not provide one
	if config.DB().Where("filename = ? AND package_version_id = ?", header.Filename, pv.ID).First(&result).RecordNotFound() {
		result.Filename = header.Filename
		result.MD5Digest = MD5(string(content))
		result.RelativePath = result.GenerateRelativePath()
		result.PackageVersionID = pv.ID
		return
	}

	return
}
