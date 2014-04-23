package main

import (
	"github.com/go-martini/martini"
	_ "github.com/jijeshmohan/mtgrid/config"
	_ "github.com/jijeshmohan/mtgrid/models"
	"github.com/martini-contrib/render"
)

var (
	m *martini.ClassicMartini
)

func init() {
	m = martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
}

func main() {
	m.Run()
}
