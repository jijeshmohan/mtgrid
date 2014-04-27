package routes

import (
	"log"

	"github.com/jijeshmohan/mtgrid/middleware"
	"github.com/jijeshmohan/mtgrid/models"
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
