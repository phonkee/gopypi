package core

import (
	"net/http"

	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/uber-go/zap"
	"github.com/phonkee/go-response"
)

// DoNotBypass is function that doesn't bypasses basic auth
var DoNotBypass = func(r *http.Request) bool {
	return false
}

/*
PostEndpointCheckMiddleware custom middleware just for post package
*/
func PostEndpointCheckMiddleware(cfg Config) alice.Constructor {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// prepare error, we will need it
			var err error

			// first parse something
			if err = r.ParseMultipartForm(2 << 32); err != nil {
				response.New(http.StatusBadRequest).Error(err).Write(w, r)
				return
			}

			var action int
			if action, err = GetPostAction(r); err != nil {
				response.New(http.StatusBadRequest).Error(err).Write(w, r)
				return
			}

			// now check action
			switch action {
			case POST_PACKAGE_ACTION_SUBMIT, POST_PACKAGE_ACTION_VERIFY:
				response.New(http.StatusNotAcceptable).Write(w, r)
				return
			case POST_PACKAGE_ACTION_PASSWORD_RESET, POST_PACKAGE_ACTION_USER:
				// This will come in the future
				response.New(http.StatusNotImplemented).Write(w, r)
				return
			case POST_PACKAGE_ACTION_DOC_UPLOAD:
				response.New(http.StatusNotImplemented).Write(w, r)
				return
			case POST_PACKAGE_ACTION_REMOVE_PKG:
				response.New(http.StatusNotAcceptable).Write(w, r)
				return
			}

			// get auth handler
			authHandler := BasicAuthLoginRequired(cfg, DoNotBypass)(h)

			// call auth handler
			authHandler.ServeHTTP(w, r)
		})
	}
}

/*
LoginRequired checks if username provided correct auth

bypassAuth is function that when returns true, auth is bypassed (!WARNING!)
*/
func BasicAuthLoginRequired(cfg Config, bypassAuth func(r *http.Request) bool, permfuncs ...func(User) error) alice.Constructor {

	// for later usage
	db := cfg.DB()

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// try bypass function
			if bypassAuth(r) {
				h.ServeHTTP(w, r)
				return
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="Gopypi restricted access"`)

			username, password, err := getUsernamePassword(r)

			if err != nil {
				response.New(http.StatusUnauthorized).Error("Not authorized").Write(w, r)
				return
			}
			user := User{}

			if db.First(&user, "username = ?", username).RecordNotFound() {
				response.New(http.StatusForbidden).Write(w, r)
				return
			}

			if !cfg.Manager().User().VerifyPassword(user, password) {
				response.New(http.StatusForbidden).Write(w, r)
				return
			}

			// call permissions callback
			if len(permfuncs) > 0 {
				for _, permission := range permfuncs {
					if err = permission(user); err != nil {
						response.New(http.StatusUnauthorized).Error(err).Write(w, r)
						return
					}
				}
			}

			// add user to context
			*r = *r.WithContext(ContextSetTokenUser(r.Context(), user))

			// call chain
			h.ServeHTTP(w, r)

		})
	}
}

func CommonMiddleware(cfg Config, router *mux.Router) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			now := time.Now()

			var (
				match     mux.RouteMatch
				routeName string
			)

			routeExists := router.Match(r, &match)
			if routeExists && match.Route != nil && len(match.Route.GetName()) > 0 {
				routeName = match.Route.GetName()
			}

			// set route name to context
			*r = *r.WithContext(
				context.WithValue(r.Context(), CONTEXT_ROUTE_NAME, routeName),
			)

			// panic recovery support
			defer func() {
				if rec := recover(); rec != nil {
					cfg.Logger().Info("request",
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
						zap.Stringer("duration", time.Now().Sub(now)),
						zap.String("name", routeName),
						zap.String("error", fmt.Sprintf("%+v", rec)),
						zap.Int("status", http.StatusInternalServerError),
						zap.String("remote", r.RemoteAddr),
					)
					cfg.Logger().Debug("stack trace:",
						zap.String("stack", string(debug.Stack())),
					)
					// write json http internal server error
					response.New(http.StatusInternalServerError).Error(rec).Write(w, r)
				}
			}()

			// call chain
			h.ServeHTTP(w, r)

			// get status from headers
			status := strings.TrimSpace(w.Header().Get(response.STATUS_HEADER))
			if status != "" {
				w.Header().Del(response.STATUS_HEADER)
			} else {
				status = "404"
			}

			// log request
			cfg.Logger().Info("request",
				zap.String("path", r.URL.Path),
				zap.String("method", r.Method),
				zap.Stringer("duration", time.Now().Sub(now)),
				zap.String("name", routeName),
				zap.String("status", status),
				zap.String("remote", r.RemoteAddr),
			)
		})
	}
}

/*
TokenAuthLoginRequired is token auth verification. It reads `gopypi-token`from headers
*/
func TokenAuthLoginRequired(cfg Config, permissions ...func(User) error) alice.Constructor {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var (
				claims *TokenClaims
				err    error
			)

			// parse token and get claims
			if claims, err = ParseToken(r, cfg.Core().SecretKey()); err != nil {
				response.New(http.StatusUnauthorized).Error(err).Write(w, r)
				return
			}

			// prepare user to be read from database
			user := User{}

			if cfg.DB().First(&user, "id = ?", claims.UserID).RecordNotFound() {
				response.New(http.StatusBadRequest).Error(err).Write(w, r)
				return
			}

			// call permissions callback
			if len(permissions) > 0 {
				for _, permission := range permissions {
					if err = permission(user); err != nil {
						response.New(http.StatusUnauthorized).Error(err).Write(w, r)
						return
					}
				}
			}

			*r = *r.WithContext(ContextSetTokenUser(r.Context(), user))

			// call next
			h.ServeHTTP(w, r)
		})
	}
}
