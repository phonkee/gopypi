package core

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

var (
	replacer *strings.Replacer
)

func init() {
	replacer = strings.NewReplacer("_", "-", ".", "-")
}

type PackageVersionPathInfo struct {
	FullPath         string
	RelativePath     string
	RelativeFilename string
	FullFilename     string
}

/*
GetPackageRelativePath returns relative filename with added directory to packages dir
*/
func GetPackageVersionFilePath(packagedir, packagename, filename string) (pvi PackageVersionPathInfo, err error) {

	prefix := ""
	if len(packagename) > 2 {
		prefix = packagename[:2]
	} else {
		prefix = packagename[:1]
	}

	pvi = PackageVersionPathInfo{
		RelativeFilename: path.Join(prefix, packagename, filename),
		RelativePath:     path.Join(prefix, packagename),
		FullPath:         path.Join(packagedir, prefix, packagename),
		FullFilename:     path.Join(packagedir, prefix, packagename, filename),
	}

	//hash := MD5(filename)
	//fmt.Sprintf("%v/%v", hash[:2], filename)

	return
}

/*
Bind json unmarshals request body to target
*/
func Bind(r *http.Request, target interface{}) (err error) {

	body := make([]byte, 0)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	defer r.Body.Close()

	if err = json.Unmarshal(body, target); err != nil {
		return
	}

	return
}

/*
GetUsernamePassword extracts username and password from Authorization header
*/
func getUsernamePassword(r *http.Request) (username, password string, err error) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(s) != 2 {
		err = ErrInvalidAuthHeader
		return
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		err = ErrInvalidAuthHeader
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		err = ErrInvalidAuthHeader
		return
	}

	return pair[0], pair[1], nil
}

/*
GenerateSalt generates salt for hashing passwords
*/
func GenerateSalt(length int) (result []byte) {
	result = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, result); err != nil {
		panic(err)
	}
	return
}

/*
MD5 shortcut to create md5 hash
*/
func MD5(input string) string {
	h := md5.New()
	io.WriteString(h, input)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*
IsEnabledOption returns whether last varargs option is enabled
*/
func IsEnabledOption(option []bool) bool {
	if len(option) == 0 {
		return false
	}

	return option[0]
}

/*
NormalizePackageName normalizes package name by PEP 503
*/
func NormalizePackageName(name string) string {
	return replacer.Replace(name)
}
