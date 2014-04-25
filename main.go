package main

import (
	"log"

	"github.com/go-martini/martini"
	"github.com/jijeshmohan/mtgrid/config"
	_ "github.com/jijeshmohan/mtgrid/models"
	"github.com/martini-contrib/render"
)

var (
	m *martini.ClassicMartini
	c *config.Config
)

func initServer() {
	m = martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
}

func startServer() {
	m.Run()
}

func main() {
	c := config.New()
	if err := c.LoadFile("app.conf"); err != nil {
		log.Fatalln("Error :", err)
	}
	if err := c.LoadEnv(); err != nil {
		log.Fatalln("Error :", err)
	}

	initServer()
	startServer()
}
