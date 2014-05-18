package routes

import (
	"log"
	"net/http"

	"github.com/jijeshmohan/mtgrid/middleware"
	"github.com/jijeshmohan/mtgrid/models"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func ListDevices(r render.Render, data middleware.Data) {
	devices, err := models.GetAllDevices()
	if err != nil {
		log.Println("error ", err)
		handleError(500, r)
		return
	}
	data["devices"] = devices
	r.HTML(200, "devices/index", data)
}

func NewDevice(r render.Render, data middleware.Data) {
	r.HTML(200, "devices/new", data)
}

func CreateDevice(r render.Render, err binding.Errors, device models.Device, data middleware.Data) {
	if err.Count() > 0 {
		data["errors"] = err
		log.Println(err)
		r.HTML(http.StatusBadRequest, "devices/new", data)
		return
	}
	if e := models.AddDevice(&device); e != nil {
		log.Println(err)
		data["errors"] = err
		log.Println(err)
		r.HTML(http.StatusBadRequest, "devices/new", data)
		return
	}
	r.Redirect("/devices")
}
