package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var Router *mux.Router

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type InsertRoutersT struct {
	Name    string
	Methods string
	Pattern string
}

var DataRouters []InsertRoutersT

func init() {
	Router = mux.NewRouter()
	Router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
}

func MiddlewareInit() {
	// Router.Use(middleware.MaiMiddleware.JWT(Router))
	// Router.Use(middleware.MaiMiddleware.Limit(Router))
}

func RouteApply(routes []Route, pathPrefix string) {
	for _, route := range routes {
		handler := route.HandlerFunc
		/*
			if os.Getenv("GO_ENV") == "development" {
				handler = middleware.MaiMiddleware.Logger(os.Stderr, route.HandlerFunc)
			} else {
				handler = route.HandlerFunc
			}*/

		if pathPrefix != "" {
			Router.
				PathPrefix(pathPrefix).
				Methods(
					route.Methods...,
				).Path(
				route.Pattern,
			). /*Name(
					route.Name,
				).*/Handler(
					handler,
				)
		} else {
			Router.
				Methods(
					route.Methods...,
				).Path(
				route.Pattern,
			). /*Name(
					route.Name,
				).*/Handler(
					handler,
				)
		}
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	if os.Getenv("GO_ENV") == "development" {
		log.Println("Error 404, at URL:", r.RequestURI)
	}
}
