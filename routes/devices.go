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
	if err.Len() > 0 {
		data["errors"] = err
		r.HTML(http.StatusBadRequest, "devices/new", data)
		return
	}
	if e := models.AddDevice(&device); e != nil {
		data["errors"] = updateError(err, e)
		r.HTML(http.StatusBadRequest, "devices/new", data)
		return
	}
	r.Redirect("/devices")
}

func updateError(err binding.Errors, e error) binding.Errors {
	newError := binding.Error{
		FieldNames:     []string{"message"},
		Classification: "Database",
		Message:        e.Error(),
	}
	if err != nil {
		return binding.Errors{newError}
	} else {
		return append(err, newError)
	}
}
