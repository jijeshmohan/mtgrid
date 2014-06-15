package models

import (
	"errors"
	"net/http"
	"strings"

	"github.com/coopernurse/gorp"
	"github.com/martini-contrib/binding"
)

type Device struct {
	Id        int64
	Name      string `form:"name" binding:"required" json:"name"`
	Type      string `form:"type" binding:"required" json:"type"`
	OS        string `form:"os" binding:"required" json:"os"`
	OsVersion string `form:"os_version" json:"os_version" binding:"required"`
	Status    string
}

func (d Device) Validate(errs binding.Errors, req *http.Request) binding.Errors {
	if len(strings.Trim(d.Name, " ")) < 5 {
		errs = append(errs, binding.Error{
			FieldNames:     []string{"name"},
			Classification: "LengthError",
			Message:        "Name is too short. Minimum 5 characters",
		})
	}

	return errs
}

func (d *Device) PreInsert(s gorp.SqlExecutor) error {
	d.Status = "Disconnected"
	d.Name = strings.ToLower(strings.Trim(d.Name, " "))
	return nil
}

func addDeviceTable() {
	tbl := orp.AddTableWithName(Device{}, "devices").SetKeys(true, "Id")
	tbl.ColMap("Name").SetUnique(true).SetMaxSize(25)

	tbl.ColMap("Type").SetMaxSize(20)
	tbl.ColMap("OS").SetMaxSize(20)
	tbl.ColMap("OsVersion").SetMaxSize(20)
	tbl.ColMap("Status").SetMaxSize(20)
}

func AddDevice(device *Device) error {
	return orp.Insert(device)
}

func GetAllDevices() (devices []Device, err error) {
	_, err = orp.Select(&devices, "select * from devices order by id")
	return
}

func GetDevice(id int64) (device *Device, err error) {
	var obj interface{}
	obj, err = orp.Get(Device{}, id)
	if err != nil || obj == nil {
		return nil, errors.New("Unable to find device")
	}
	device = obj.(*Device)
	return
}

func GetDeviceWithName(name string) (*Device, error) {
	var d Device
	err := orp.SelectOne(&d, "select * from devices where [name]=:name", map[string]interface{}{"name": name})
	return &d, err
}

func UpdateDeviceStatus(device *Device, status string) (err error) {
	device.Status = status
	_, err = orp.Update(device)
	return err
}
