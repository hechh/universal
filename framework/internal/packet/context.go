package packet

import "universal/framework/define"

type Context struct {
	define.IHeader
	router define.IRouter
}

func NewContext(h define.IHeader, r define.IRouter) *Context {
	return &Context{IHeader: h, router: r}
}

func (d *Context) GetRouter() define.IRouter {
	return d.router
}
