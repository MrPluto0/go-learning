package gee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type router struct {
	roots    map[string]*node       // Map[Method String] = Node
	handlers map[string]HandlerFunc // Map[Key] = Handler
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			if v[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler

	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	params := make(map[string]string, 0)
	searchParts := parsePattern(path)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	target := root.search(searchParts, 0)
	if target == nil {
		return nil, nil
	}

	parts := parsePattern(target.pattern)
	for i, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[i]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = "/" + strings.Join(searchParts[i:], "/")
		}
	}
	return target, params
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		// store to 'handlers', and call 'Next' so that all middlewares are called.
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
