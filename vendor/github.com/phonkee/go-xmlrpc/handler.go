package xmlrpc

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/beevik/etree"
)

var (
	ErrNotService               = errors.New("argument is not service")
	ErrServiceAlreadyRegistered = errors.New("service already registered")
	ErrServiceMustBePointer     = errors.New("registered service must be pointer")
)

/*
Handler interface
*/
type Handler interface {

	// AddService under given namespace
	AddService(service interface{}, name string) error

	// ListMethods returns all available xmlrpc methods
	ListMethods() []string

	// ServeHTTP satisfy http.Handler
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

/*
NewHandler returns xmlrpc handler
*/
func NewHandler() Handler {
	return &handler{
		services: map[string]Service{},
	}
}

type handler struct {
	services map[string]Service
}

/*
AddService adds service to handler, you ask why interface{} and not Service? Well after cast to Service we can write
some meaningful error message.
*/
func (h *handler) AddService(service interface{}, name string) error {

	// True whether it's really service
	if s, ok := service.(Service); !ok {
		return ErrNotService
	} else {

		// check if it's pointer receiver
		if reflect.TypeOf(service).Kind() != reflect.Ptr {
			return ErrServiceMustBePointer
		}

		// check if service already exists
		if _, ok := h.services[name]; ok {
			return ErrServiceAlreadyRegistered
		}
		h.services[name] = s
	}

	return nil
}

/*
ListMethods returns list of all available XML rpc methods
*/
func (h *handler) ListMethods() []string {
	result := []string{}
	for name, service := range h.services {
		for _, method := range service.ListMethods() {
			if name == "" {
				result = append(result, method)
			} else {
				result = append(result, name+"."+method)
			}
		}
	}
	return result
}

/*
ServeHTTP serves http
*/
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// check for POST method
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/xml")

	var (
		doc *etree.Document
		err error
	)

	resultDoc := etree.NewDocument()
	resultDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	// create new document
	doc = etree.NewDocument()

	if _, err = doc.ReadFrom(r.Body); err != nil {
		faultStruct := resultDoc.CreateElement("methodResponse").CreateElement("fault").CreateElement("value")
		XMLWriteError(faultStruct, errors.New("cannot parse body"))
		resultDoc.WriteTo(w)
		return
	}

	var method string

	if element := doc.FindElement("methodCall/methodName"); element == nil {
		faultStruct := resultDoc.CreateElement("methodResponse").CreateElement("fault").CreateElement("value")
		XMLWriteError(faultStruct, errors.New("methodName not found"))
		resultDoc.WriteTo(w)
		return
	} else {
		method = element.Text()
	}

	// list methods serve directly
	if method == "system.listMethods" {
		value := resultDoc.CreateElement("methodResponse").CreateElement("params").CreateElement("param").CreateElement("value")
		availMethods := h.ListMethods()
		availMethods = append(availMethods, "system.listMethods")
		XMLWriteStringSlice(value, availMethods)
		resultDoc.WriteTo(w)
		return
	}

	// now we need to split methods by dot make a lookup and perform
	splitted := strings.SplitN(method, ".", 2)
	service := ""

	if len(splitted) == 2 {
		service = splitted[0]
		method = splitted[1]
	}

	s, ok := h.services[service]
	if !ok {
		faultStruct := resultDoc.CreateElement("methodResponse").CreateElement("fault").CreateElement("value")
		XMLWriteError(faultStruct, errors.New("method not found"))
		resultDoc.WriteTo(w)
		return
	}

	if !s.MethodExists(method) {
		faultStruct := resultDoc.CreateElement("methodResponse").CreateElement("fault").CreateElement("value")
		XMLWriteError(faultStruct, errors.New("method not found"))
		resultDoc.WriteTo(w)
		return
	}

	el := doc.FindElement("methodCall/params")
	if el == nil {
		faultStruct := resultDoc.CreateElement("methodResponse").CreateElement("fault").CreateElement("value")
		XMLWriteError(faultStruct, errors.New("params not found"))
		resultDoc.WriteTo(w)
		return
	}

	var (
		res *etree.Document
	)

	// call dispatch and write result
	res, err = s.Dispatch(method, el)
	if err != nil {
		faultStruct := resultDoc.CreateElement("methodResponse").CreateElement("fault").CreateElement("value")
		XMLWriteError(faultStruct, err)
		resultDoc.WriteTo(w)
		return
	}

	res.WriteTo(w)

	return
}
