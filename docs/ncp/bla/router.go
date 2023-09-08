package bla

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

type Handler func(req *Request) (*Reply, error)

func NewRouter() *Router {
	return &Router{
		Routes: []*Route{},
	}
}

type Router struct {
	Routes []*Route
}

var ErrRouteNotFound = errors.New("route not found")

func (r *Router) Serve(req *Request) (*Reply, error) {
	for _, route := range r.Routes {
		params, ok := route.Match(req)
		if !ok {
			continue
		}
		req.Params = params
		return route.Handler(req)
	}
	slog.Error("route not found", "method", req.Method, "path", req.RequestURI)
	return nil, ErrRouteNotFound
}

func (r *Router) Get(pattern string, handler Handler) {
	r.Routes = append(r.Routes, NewRoute(http.MethodGet, pattern, handler))
}

func NewRoute(method string, pattern string, handler Handler) *Route {
	// Remove the leading slash and split by forward slash
	steps := strings.Split(strings.TrimPrefix(pattern, "/"), "/")

	route := &Route{
		Handler: handler,
		Method:  method,
		Pattern: pattern,
	}

	for _, step := range steps {
		switch {
		case step == "*":
			// It is expected that the wildcar comes last in the path.
			route.Steps = append(route.Steps, Step{
				Wildcard: true,
			})
			route.Wildcard = true
			return route
		case len(step) > 2 && step[0] == '{' && step[len(step)-1] == '}':
			route.Steps = append(route.Steps, Step{
				Param: step[1 : len(step)-1],
			})
		default:
			route.Steps = append(route.Steps, Step{
				Value: step,
			})
		}
	}

	return route
}

type Route struct {
	// Handler is the function to call to handle this route.
	Handler Handler
	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	Method   string
	Pattern  string
	Wildcard bool

	Steps []Step
}

func (r Route) Match(req *Request) (map[string]string, bool) {
	// Check the HTTP method matches.
	if r.Method != req.Method {
		return nil, false
	}
	steps := strings.Split(strings.TrimPrefix(req.RequestURI, "/"), "/")

	// If the number of steps does not match, and there is no wildcard,
	// then there is no match.
	if len(r.Steps) != len(steps) &&
		!(r.Wildcard && len(r.Steps) < len(steps)) {
		return nil, false
	}

	params := make(map[string]string)
	for i, step := range r.Steps {
		switch {
		case step.Param != "":
			params[step.Param] = steps[i]
			continue
		case step.Wildcard:
			params["*"] = strings.Join(steps[i:], "/")
			return params, true
		case step.Value == steps[i]:
			continue
		}
		return nil, false
	}

	return params, true
}

type Step struct {
	Value    string
	Param    string
	Wildcard bool
}
