package router

import (
	"core/router/types"
	"core/server/dispatch"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Router struct {
	router *router.Router
}

type Group struct {
	group *router.Group
}

func New() *Router {
	return &Router{
		router: router.New(),
	}
}

func (r *Router) Group(path string) *Group {
	return &Group{
		group: r.router.Group(path),
	}
}

func (r *Router) Mutable(v bool) {
	r.router.Mutable(v)
}

func (r *Router) List() map[string][]string {
	return r.router.List()
}

// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handler)
func (r *Router) GET(path string, handler types.RequestHandler) {
	r.router.GET(path, dispatch.Dispatch(handler))
}

// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handler)
func (r *Router) HEAD(path string, handler types.RequestHandler) {
	r.router.HEAD(path, dispatch.Dispatch(handler))
}

// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handler)
func (r *Router) POST(path string, handler types.RequestHandler) {
	r.router.POST(path, dispatch.Dispatch(handler))
}

// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handler)
func (r *Router) PUT(path string, handler types.RequestHandler) {
	r.router.PUT(path, dispatch.Dispatch(handler))
}

// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handler)
func (r *Router) PATCH(path string, handler types.RequestHandler) {
	r.router.PATCH(path, dispatch.Dispatch(handler))
}

// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handler)
func (r *Router) DELETE(path string, handler types.RequestHandler) {
	r.router.DELETE(path, dispatch.Dispatch(handler))
}

// CONNECT is a shortcut for router.Handle(fasthttp.MethodConnect, path, handler)
func (r *Router) CONNECT(path string, handler types.RequestHandler) {
	r.router.CONNECT(path, dispatch.Dispatch(handler))
}

// OPTIONS is a shortcut for router.Handle(fasthttp.MethodOptions, path, handler)
func (r *Router) OPTIONS(path string, handler types.RequestHandler) {
	r.router.OPTIONS(path, dispatch.Dispatch(handler))
}

// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handler)
func (r *Router) TRACE(path string, handler types.RequestHandler) {
	r.router.TRACE(path, dispatch.Dispatch(handler))
}

// ANY is a shortcut for router.Handle(router.MethodWild, path, handler)
//
// WARNING: Use only for routes where the request method is not important
func (r *Router) ANY(path string, handler types.RequestHandler) {
	r.router.ANY(path, dispatch.Dispatch(handler))
}

// ServeFiles serves files from the given file system root.
// The path must end with "/{filepath:*}", files are then served from the local
// path /defined/root/dir/{filepath:*}.
// For example if root is "/etc" and {filepath:*} is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a fasthttp.FSHandler is used, therefore fasthttp.NotFound is used instead
// Use:
//
//	router.ServeFiles("/src/{filepath:*}", "./")
func (r *Router) ServeFiles(path string, rootPath string) {
	r.router.ServeFiles(path, rootPath)
}

// ServeFilesCustom serves files from the given file system settings.
// The path must end with "/{filepath:*}", files are then served from the local
// path /defined/root/dir/{filepath:*}.
// For example if root is "/etc" and {filepath:*} is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a fasthttp.FSHandler is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// Use:
//
//	router.ServeFilesCustom("/src/{filepath:*}", *customFS)
func (r *Router) ServeFilesCustom(path string, fs *fasthttp.FS) {
	r.router.ServeFilesCustom(path, fs)
}

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handler types.RequestHandler) {
	r.router.Handle(method, path, dispatch.Dispatch(handler))
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handler function.
// Otherwise the second return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string, ctx *fasthttp.RequestCtx) (fasthttp.RequestHandler, bool) {
	return r.router.Lookup(method, path, ctx)
}

// Handler makes the router implement the http.Handler interface.
func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	r.router.Handler(ctx)
}
