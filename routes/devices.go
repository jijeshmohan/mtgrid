package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jijeshmohan/mtgrid/middleware"
	"github.com/jijeshmohan/mtgrid/models"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

func ListDevices(r render.Render, data middleware.Data, req *http.Request) {
	devices, err := models.GetAllDevices()
	if err != nil {
		log.Println("error ", err)
		handleError(500, r)
		return
	}
	data["devices"] = devices
	data["host"] = req.Host
	r.HTML(200, "devices/index", data)
}

func NewDevice(r render.Render, data middleware.Data) {
	r.HTML(200, "devices/new", data)
}

func CreateDevice(r render.Render, err binding.Errors, device models.Device, data middleware.Data, req *http.Request) {
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
	websocketClient(req.Host, device)
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

func websocketClient(host string, device models.Device) {
	dlr := &websocket.Dialer{}
	con, _, err := dlr.Dial("ws://"+host+"/sock", http.Header{})
	if err != nil {
		log.Println(err)
	}
	d := struct {
		Status string
		Device models.Device
	}{
		"CREATED",
		device,
	}
	err = con.WriteJSON(d)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	con.Close()
}
