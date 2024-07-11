package main

import (
	"log"
	"net/http"
	"path"
)

type RouterGroup struct {
	basePath string
	handles  []HandlerFunc

	engine *Engine
	parent *RouterGroup
}

func (g *RouterGroup) Group(basePath string, handles ...HandlerFunc) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		basePath: g.calculateAbsolutePath(basePath),
		handles:  g.combineHandlers(handles...),
		parent:   g,
		engine:   engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, relativePath string, handlerFunc ...HandlerFunc) {
	absolutePath := g.calculateAbsolutePath(relativePath)
	handlers := g.combineHandlers(handlerFunc...)
	log.Printf("\nRoute %4s - %s\n", method, absolutePath)
	g.engine.router.addRoute(method, absolutePath, handlers...)
}

func (g *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(g.basePath, relativePath)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func (g *RouterGroup) GET(pattern string, handlerFunc HandlerFunc) {
	g.addRoute(http.MethodGet, pattern, handlerFunc)
}

func (g *RouterGroup) POST(pattern string, handlerFunc HandlerFunc) {
	g.addRoute(http.MethodPost, pattern, handlerFunc)
}

func (g *RouterGroup) Use(handles ...HandlerFunc) {
	g.handles = append(g.handles, handles...)
}

func (g *RouterGroup) combineHandlers(handles ...HandlerFunc) []HandlerFunc {
	size := len(g.handles) + len(handles)
	ret := make([]HandlerFunc, size)
	copy(ret, g.handles)
	copy(ret[len(g.handles):], handles)

	return ret
}
