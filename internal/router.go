package internal

import "github.com/julienschmidt/httprouter"

func Router(routes []Route) *httprouter.Router {
	router := httprouter.New()

	for _, route := range routes {
		var handle httprouter.Handle = route.HandlerFunc

		router.Handle(route.Method, route.Path, handle)
	}
	return router
}
