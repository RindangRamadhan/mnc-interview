package example

import "github.com/RindangRamadhan/mnc-interview/internal/router"

var routes = []router.Route{
	{
		Methods:     []string{"GET"},
		Pattern:     "/example",
		HandlerFunc: ExampleHandlers.List,
	},
}

func init() {
	router.RouteApply(routes, "/api")
}
