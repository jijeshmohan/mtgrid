package main

import (
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jijeshmohan/mtgrid/config"
	"github.com/jijeshmohan/mtgrid/models"
	"github.com/martini-contrib/render"
)

var (
	m *martini.ClassicMartini
	c *config.Config
)

func startServer() {
	m = martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	if _, err := models.InitDb(c); err != nil {
		log.Fatalln("Error :", err)
	}

	log.Println("Starting server at ", c.Addr)
	http.ListenAndServe(c.Addr, m)
}

func main() {
	c = config.New()
	if err := c.LoadFile("app.conf"); err != nil {
		log.Fatalln("Error :", err)
	}
	if err := c.LoadEnv(); err != nil {
		log.Fatalln("Error :", err)
	}

	startServer()
}
