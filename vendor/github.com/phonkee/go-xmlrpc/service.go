/*
Protocol handles first level of marshalling/unmarshalling
*/
package xmlrpc

import "github.com/beevik/etree"

/*
Service interface must be satisfied when registering service
This interface is satisfied when go generate is called
*/
type Service interface {

	// Dispatch method dispatches xmlrpc call
	Dispatch(method string, root *etree.Element) (result *etree.Document, err error)

	// returns list of all rpc methods
	ListMethods() []string

	// check if service has method
	MethodExists(string) bool
}
