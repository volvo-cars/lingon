package bla

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestRouter(t *testing.T) {

	handle := func(req *Request) (*Reply, error) {
		// tu.AssertEqual(t, http.MethodGet, req.Method)

		// JSON marshal the request and return it as the body.
		bData, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("marshal: %w", err)
		}
		return &Reply{
			Body: bData,
		}, nil
	}
	// router := NewRouter()
	// router.Get("/", handle)
	// router.Get("/foo", handle)
	// router.Get("/foo/bar", handle)
	// router.Get("/foo/{param}", handle)
	// router.Get("/foo/*", handle)

	type test struct {
		request   *Request
		pattern   string
		expMatch  bool
		expParams map[string]string
	}
	tests := []test{
		{
			request:   methodAndPathRequest(http.MethodGet, "/"),
			pattern:   "/",
			expMatch:  true,
			expParams: nil,
		},
		{
			request:   methodAndPathRequest(http.MethodGet, "/foo"),
			pattern:   "/foo",
			expMatch:  true,
			expParams: nil,
		},
		{
			request:   methodAndPathRequest(http.MethodGet, "/foo/bar"),
			pattern:   "/foo/bar",
			expMatch:  true,
			expParams: nil,
		},
		{
			request:   methodAndPathRequest(http.MethodGet, "/foo/bar"),
			pattern:   "/foo/{param}",
			expMatch:  true,
			expParams: map[string]string{"param": "bar"},
		},
		{
			request:   methodAndPathRequest(http.MethodGet, "/wildcard"),
			pattern:   "/*",
			expMatch:  true,
			expParams: map[string]string{"*": "wildcard"},
		},
		{
			request:   methodAndPathRequest(http.MethodGet, "/wildcard/nested"),
			pattern:   "/wildcard/*",
			expMatch:  true,
			expParams: map[string]string{"*": "nested"},
		},
		{
			request: methodAndPathRequest(
				http.MethodGet,
				"/wildcard/nested/deep",
			),
			pattern:   "/wildcard/*",
			expMatch:  true,
			expParams: map[string]string{"*": "nested/deep"},
		},
		{
			request: methodAndPathRequest(
				http.MethodGet,
				"/wildcard/param/bla/nested",
			),
			pattern:   "/wildcard/param/{param}/*",
			expMatch:  true,
			expParams: map[string]string{"param": "bla", "*": "nested"},
		},
		{
			request: methodAndPathRequest(
				http.MethodGet,
				"/wildcard/param/bla/nested/deep",
			),
			pattern:   "/wildcard/param/{param}/*",
			expMatch:  true,
			expParams: map[string]string{"param": "bla", "*": "nested/deep"},
		},
		{
			request: methodAndPathRequest(
				http.MethodGet,
				"/foo",
			),
			pattern:   "/",
			expMatch:  false,
			expParams: nil,
		},
	}
	for _, test := range tests {
		t.Run(
			test.request.Method+":"+test.request.RequestURI,
			func(t *testing.T) {
				router := NewRouter()
				switch test.request.Method {
				case http.MethodGet:
					router.Get(test.pattern, handle)
				// case http.MethodPost:
				default:
					t.Fatalf("unsupported method %s", test.request.Method)
				}

				reply, err := router.Serve(test.request)
				if test.expMatch {
					tu.AssertNoError(t, err)
				} else {
					tu.AssertErrorMsg(t, err, ErrRouteNotFound.Error())
					return
				}

				var replyReq Request
				err = json.Unmarshal(reply.Body, &replyReq)
				tu.AssertNoError(t, err, "unmarshal reply")

				// Assert the params match
				tu.True(
					t,
					mapsEqual(test.expParams, replyReq.Params),
					fmt.Sprintf(
						"expected %v, got %v",
						test.expParams,
						replyReq.Params,
					),
				)
			})
	}
}

func methodAndPathRequest(method string, path string) *Request {
	return &Request{
		Method:     method,
		RequestURI: path,
	}
}

func mapsEqual(a map[string]string, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
