package routes

import (
	"fmt"

	"github.com/go-martini/martini"
	"github.com/jijeshmohan/mtgrid/models"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func InitRoutes(m *martini.ClassicMartini) {

	m.Group("/devices", func(r martini.Router) {
		r.Get("", ListDevices)
		r.Post("", binding.Form(models.Device{}), CreateDevice)
		r.Get("/new", NewDevice)
		// r.Get("/:id", ShowDevice)
	})

	m.Get("/sock", Socket)
	m.NotFound(func(r render.Render) {
		r.HTML(404, "status/404", "")
	})
}

func handleError(err int, r render.Render) {
	template := fmt.Sprintf("status/%d", err)
	r.HTML(err, template, "")
}
