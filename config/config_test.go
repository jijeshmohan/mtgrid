package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("Load a configuration file that does exist", t, func() {
		c, err := NewConfig("invalid")
		So(err, ShouldNotBeNil)
		So(c, ShouldBeNil)
	})

	Convey("Default to development config if the value empty", t, func() {
		c, err := NewConfig("")
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)
		So(c.Dbpath, ShouldEqual, "db_development.db")
	})
}
