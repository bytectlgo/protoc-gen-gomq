package mqtt

import (
	"net/http"
	"path"

	xhttp "github.com/go-kratos/kratos/v2/transport/http"
)

// WalkRouteFunc is the type of the function called for each route visited by Walk.
type WalkRouteFunc func(RouteInfo) error

// RouteInfo is an MQTT route info.
type RouteInfo struct {
	Path   string
	Method string
}

// HandlerFunc defines a function to serve MQTT requests.
type HandlerFunc func(Context) error

// Router is an MQTT router.
type Router struct {
	prefix  string
	srv     *Server
	filters []xhttp.FilterFunc
}

func newRouter(prefix string, srv *Server, filters ...xhttp.FilterFunc) *Router {
	r := &Router{
		prefix:  prefix,
		srv:     srv,
		filters: filters,
	}
	return r
}

// Group returns a new router group.
func (r *Router) Group(prefix string, filters ...xhttp.FilterFunc) *Router {
	var newFilters []xhttp.FilterFunc
	newFilters = append(newFilters, r.filters...)
	newFilters = append(newFilters, filters...)
	return newRouter(path.Join(r.prefix, prefix), r.srv, newFilters...)
}

// Handle registers a new route with a matcher for the URL path and method.
func (r *Router) Handle(method, relativePath string, h HandlerFunc, filters ...xhttp.FilterFunc) {
	next := http.Handler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := &wrapper{router: r}
		ctx.Reset(res, req)
		if err := h(ctx); err != nil {
			r.srv.ene(res, req, err)
		}
	}))
	next = xhttp.FilterChain(filters...)(next)
	next = xhttp.FilterChain(r.filters...)(next)
	r.srv.router.Handle(path.Join(r.prefix, relativePath), next).Methods(method)
}

// POST registers a new POST route for a path with matching handler in the router.
func (r *Router) POST(path string, h HandlerFunc, m ...xhttp.FilterFunc) {
	r.Handle(http.MethodPost, path, h, m...)
}
