package gluon

import "github.com/justinas/alice"

func (r *router) Use(middleware ...alice.Constructor) {
	r.middlewareStack = r.middlewareStack.Append(middleware...)
}
