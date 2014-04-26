package models

import "errors"

type Device struct {
	Id        int64
	Name      string `binding:"required"`
	Type      string ``
	OS        string
	OsVersion string
	Status    string
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

func GetAllDevices() (devices *[]Device, err error) {
	_, err = orp.Select(devices, "select * from devices order by id")
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
