package main

import (
	"log"
	"net/http"
	"path"
)

type RouterGroup struct {
	basePath    string
	middlewares []HandlerFunc

	engine *Engine
	parent *RouterGroup
}

func (g *RouterGroup) Group(basePath string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		basePath: g.calculateAbsolutePath(basePath),
		parent:   g,
		engine:   engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, relativePath string, handlerFunc HandlerFunc) {
	absolutePath := g.calculateAbsolutePath(relativePath)
	log.Printf("Route %4s - %s", method, absolutePath)
	g.engine.router.addRoute(method, absolutePath, handlerFunc)
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
