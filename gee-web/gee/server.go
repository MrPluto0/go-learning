package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type RouterGroup struct {
	prefix       string
	middlewares  []HandlerFunc // support middleware
	parent       *RouterGroup  // support nesting
	engine       *Engine       // all groups share a Engine instance
	htmlTemplate *template.Template
	funcMap      template.FuncMap
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplate = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix:      g.prefix + prefix,
		middlewares: make([]HandlerFunc, 0),
		parent:      g,
		engine:      g.engine,
	}
	g.engine.groups = append(g.engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Static File System
func (g *RouterGroup) createStaticServer(relativePath string, root string) HandlerFunc {
	fs := http.Dir(root)              // file system
	absoluteFs := http.FileServer(fs) // file server
	relativeFs := http.StripPrefix(relativePath, absoluteFs)
	return func(c *Context) {
		path := c.Param("filepath")
		if _, err := fs.Open(path); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		relativeFs.ServeHTTP(c.W, c.R)
	}
}

func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.createStaticServer(relativePath, root)
	pattern := path.Join(relativePath, "/*filepath")
	g.GET(pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodGet, pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute(http.MethodPost, pattern, handler)
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

func (e *Engine) Run(addr string) {
	log.Printf("Listening at %s\n", addr)
	http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)

	var middlewares []HandlerFunc
	for _, g := range e.groups {
		if strings.HasPrefix(c.Path, g.prefix) {
			middlewares = append(middlewares, g.middlewares...)
		}
	}
	c.handlers = middlewares
	c.engine = e

	e.router.handle(c)
}
