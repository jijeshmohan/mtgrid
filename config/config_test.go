package config

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	Convey("configuration should have default values", t, func() {
		c := New()
		So(c, ShouldNotBeNil)
		So(c.Addr, ShouldEqual, ":3000")
		So(c.DbType, ShouldEqual, "sqlite3")
		So(c.DbPath, ShouldContainSubstring, "development.db")
	})

	Convey("Load config file", t, func() {
		Convey("should override default values", func() {
			content := `db_path="testing.db"`
			withTempFile(content, func(tempFile string) {
				c := New()
				So(c, ShouldNotBeNil)
				So(c.DbPath, ShouldContainSubstring, "development.db")
				So(c.LoadFile(tempFile), ShouldBeNil)
				So(c.DbPath, ShouldContainSubstring, "testing.db")
			})
		})

		Convey("should ignore invalid values present in the file", func() {
			content := `db_path1="testing.db"`
			withTempFile(content, func(tempFile string) {
				c := New()
				So(c, ShouldNotBeNil)
				So(c.DbPath, ShouldContainSubstring, "development.db")
				So(c.LoadFile(tempFile), ShouldBeNil)
				So(c.DbPath, ShouldContainSubstring, "development.db")
			})
		})
	})

}

func withTempFile(content string, fn func(string)) {
	f, _ := ioutil.TempFile("", "")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	fn(f.Name())
}
