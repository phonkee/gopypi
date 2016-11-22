/*

Classy package

Support for class based views. Inspired by django class based views. Every structure can be view when it provides needed methods.
Classy has som pre built configured base views, so you don't need to create one. They are built with rest principles in mind.
Design of classy is that you can provide mapping from http methods to view methods. Classy then registers views to
gorilla mux router.
Classy uses reflection quite a lot, but the advantage is that it's used only during registration to router. During the
run of server it doesn't affects speed and performance

Classy has also support for "ViewSet" as we know them from django rest framework, so you can combine list/detail view
in single struct.

Example:

Let's create users list class based view. We will use predefined ListView.

	type UserListView struct {
		ListView
	}




*/
package classy
