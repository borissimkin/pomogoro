package router

import tea "github.com/charmbracelet/bubbletea"

type RouteKey string

type RouteValue tea.Model

type Route struct {
	Key   RouteKey
	Value RouteValue
}

func NewRoute(key RouteKey, value RouteValue) Route {
	return Route{key, value}
}

type Router struct {
	Routes       map[RouteKey]Route
	currentRoute RouteKey
}

func NewRouter() Router {
	return Router{
		Routes: make(map[RouteKey]Route),
	}
}

func (r *Router) SetRoutes(routes []Route) {
	if len(routes) == 0 {
		panic("routes should not be empty")
	}

	for _, route := range routes {
		r.Routes[route.Key] = route
	}

	r.currentRoute = routes[0].Key
}

func (r *Router) CurrentRoute() Route {
	return r.Routes[r.currentRoute]
}

func (r *Router) To(key RouteKey) (tea.Model, tea.Cmd) {
	r.currentRoute = key

	route := r.Routes[key]

	return route.Value, route.Value.Init()
}
