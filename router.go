package core

import (
	"errors"

	"github.com/justinas/alice"

	"net/http"
	"github.com/phonkee/go-classy"
	"github.com/phonkee/go-xmlrpc"
	"github.com/elazarl/go-bindata-assetfs"
)

/*
GetRouter instantiates router with all its registered routes
*/
func InitRouter(config Config) (chain alice.Chain, err error) {

	router := config.Router()

	// create base middlewares chain
	chain = alice.New(CommonMiddleware(config, router))

	// enable debug for all classy views (for now)
	classy.Debug()

	// prepare middleware for logged user with list permission
	listAuth := BasicAuthLoginRequired(config, DoNotBypass, func(user User) (err error) {
		if !user.CanList {
			return ErrUserCannotRetrievePackages
		}
		return
	})

	// prepare middleware for logged user with download permission
	downloadAuth := BasicAuthLoginRequired(config, DoNotBypass, func(user User) (err error) {
		if !user.CanDownload {
			return ErrUserCannotDownloadPackages
		}
		return
	})

	// packages routes
	classy.Path("/packages").Register(
		router,
		classy.New(&PackageDownloadView{Config: config}).Use(downloadAuth),
	)

	// register /simple package list path
	classy.Path("/simple").Use(listAuth).Register(
		router,
		classy.New(&PackageListView{Config: config}),
		classy.New(&PackageDetailView{Config: config}),
	)

	// Register homepage and post package view
	classy.Register(
		router,
		classy.New(
			TemplateView{
				Config:       config,
				TemplateName: "index.tpl.html",
				Data: map[string]interface{}{
					"version": VERSION,
				},
			},
		).Name("homepage"),

		// post endpoint has special middleware
		classy.New(&PostPackageView{Config: config}).Use(PostEndpointCheckMiddleware(config)),
	)

	// prepare token auth for admin
	adminAuth := TokenAuthLoginRequired(config, func(user User) (err error) {
		if !user.IsActive {
			return errors.New("user inactive")
		}
		if !user.IsAdmin {
			return errors.New("only user with admin access allowed")
		}
		return
	})

	// login api route (not secured by token auth)
	classy.Name("api:{name}").Path("/api/login").Register(
		router,

		classy.New(&LoginAPIView{Config: config}),
	)

	// api endpoints secured by token auth
	classy.Name("api:{name}").Path("/api").Use(adminAuth).Register(
		router,

		classy.New(&FeatureAPIViewSet{Config: config}).Path("/feature"),
		classy.New(&InfoAPIView{Config: config}).Path("/info"),
		classy.New(&LicenseAPIViewSet{Config: config}).Path("/license"),

		// me views - all about current logged user (by token)
		classy.Group(
			"/me",
			classy.New(&MeAPIView{Config: config}),
			classy.New(&MeChangePasswordAPIView{Config: config}).Path("/password"),
			classy.New(&MyPackageAPIView{Config: config}).
				Path("/package"),
		),

		// package views
		classy.Group(
			"/package",
			classy.New(&PackageAPIViewSet{Config: config}),
			classy.New(&PackageMaintainerAPIViewSet{Config: config}).
				Path("/{package_pk:[0-9]+}/maintainer/"),
		),

		// platform views
		classy.New(&PlatformAPIViewSet{Config: config}).Path("/platform"),

		// stat classy views
		classy.Group(
			"/stats",
			classy.New(&StatsAPIView{Config: config}).Path("/server"),
			classy.Group(
				"/download/package",
				classy.New(&StatsDownloadAllAPIView{Config: config}),
				classy.New(&StatsDownloadPackageAPIView{Config: config}),
				classy.New(&StatsDownloadPackageVersionAPIView{Config: config}).
					Path("/{package_pk:[0-9]+}/version"),
			),
		),

		classy.New(&UserAPIViewSet{Config: config}).Path("/user"),
	)

	// register rpc service
	xmlrpcHandler := xmlrpc.NewHandler()
	if e := xmlrpcHandler.AddService(&SearchService{Config: config}, ""); e != nil {
		println("error:", e.Error())
	}

	router.Handle("/RPC2", alice.New(listAuth).Then(xmlrpcHandler)).Methods("POST")

	/*
		admin static handler with fallback
	*/
	router.PathPrefix("/admin/").
		Handler(
			http.StripPrefix(
				"/admin/", http.FileServer(
					&FallbackFileSystem{
						original: &assetfs.AssetFS{
							Asset:     Asset,
							AssetDir:  AssetDir,
							AssetInfo: AssetInfo,
							Prefix:    "admin",
						},
						fallback: "index.html",
					})))
	// add alias
	router.Handle("/admin", http.RedirectHandler("/admin/", http.StatusMovedPermanently))
	return
}
