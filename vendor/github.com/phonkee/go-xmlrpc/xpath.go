/*
Xpath function helpers
*/
package xmlrpc

import (
	"strconv"

	"encoding/base64"
	"strings"

	"github.com/beevik/etree"
)

/*
XPathValueGetBytes Returns []byte from value
*/
func XPathValueGetBytes(element *etree.Element, name string) (result []byte, err error) {
	tmp := element.FindElement("base64")

	if tmp == nil {
		err = Errorf(400, "Not found %v")
		return
	}

	result, err = base64.StdEncoding.DecodeString(tmp.Text())

	return
}

/*
XPathValueGetFloat64 Returns float64 from value
*/
func XPathValueGetFloat64(element *etree.Element, name string) (result float64, err error) {
	tmp := element.FindElement("double")

	if tmp == nil {
		err = Errorf(400, "Not found %v")
		return
	}

	result, err = strconv.ParseFloat(tmp.Text(), 64)
	return
}

/*
XPathValueGetFloat32 Returns float64 from value
*/
func XPathValueGetFloat32(element *etree.Element, name string) (result float32, err error) {
	tmp := element.FindElement("double")

	if tmp == nil {
		err = Errorf(400, "Not found %v")
		return
	}

	var f64 float64

	if f64, err = strconv.ParseFloat(tmp.Text(), 32); err == nil {
		result = float32(f64)
	}

	return
}

/*
XPathValueGetInt Returns int from value
*/
func XPathValueGetInt(element *etree.Element, name string) (result int, err error) {

	var tmp *etree.Element

	if tmp = element.FindElement("int"); tmp == nil {
		tmp = element.FindElement("i4")
	}

	if tmp == nil {
		err = Errorf(400, "Not found %v")
		return
	}

	result, err = strconv.Atoi(tmp.Text())

	return
}

/*
XPathValueGetInt64 Returns int64 from value
*/
func XPathValueGetInt64(element *etree.Element, name string) (result int64, err error) {
	var i int

	if i, err = XPathValueGetInt(element, name); err != nil {
		return
	}

	result = int64(i)

	return
}

/*
XPathValueGetInt64 Returns int64 from value
*/
func XPathValueGetInt32(element *etree.Element, name string) (result int32, err error) {
	var i int

	if i, err = XPathValueGetInt(element, name); err != nil {
		return
	}

	result = int32(i)

	return
}

/*
XPathValueGetString Returns bool from value
*/
func XPathValueGetString(element *etree.Element, name string) (result string, err error) {
	tmp := element.FindElement("string")
	if tmp == nil {
		err = Errorf(400, "not found %v", name)
	}

	result = tmp.Text()

	return
}

/*
XPathValueGetBool Returns bool from value
*/
func XPathValueGetBool(element *etree.Element, name string) (result bool, err error) {
	tmp := element.FindElement("boolean")
	if tmp == nil {
		err = Errorf(400, "not found %v", name)
	}

	result = strings.TrimSpace(tmp.Text()) == "1"

	return
}
