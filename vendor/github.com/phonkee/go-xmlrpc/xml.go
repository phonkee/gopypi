package xmlrpc

import (
	"strconv"

	"github.com/beevik/etree"
)

/*
XMLWriteError writes xml error
*/
func XMLWriteError(element *etree.Element, err error) {

	faultCode := 500

	if e, ok := err.(Error); ok {
		faultCode = e.Code()
	}

	faultStruct := element.CreateElement("struct")
	m1 := faultStruct.CreateElement("member")
	m1.CreateElement("name").SetText("faultCode")
	m1.CreateElement("value").SetText(strconv.Itoa(faultCode))

	m2 := faultStruct.CreateElement("member")
	m2.CreateElement("name").SetText("faultString")
	m2.CreateElement("value").SetText(err.Error())
}

/*
XMLWriteStringSlice writes array of string slice
 */
func XMLWriteStringSlice(element *etree.Element, stringSlice []string) {
	dataElement := element.CreateElement("array").CreateElement("data")
	for _, item := range stringSlice {
		dataElement.CreateElement("value").CreateElement("string").SetText(item)
	}
}
