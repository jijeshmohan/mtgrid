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
			withTempFile(content, func(c *Config, tempFile string) {
				So(c.DbPath, ShouldContainSubstring, "development.db")
				So(c.LoadFile(tempFile), ShouldBeNil)
				So(c.DbPath, ShouldContainSubstring, "testing.db")
			})
		})

		Convey("should throw error when the file is missing", func() {
			c := New()
			So(c.LoadFile("invalidfile"), ShouldNotBeNil)
			So(c.DbPath, ShouldContainSubstring, "development.db")
		})

		Convey("should throw error when the file is in different format", func() {
			content := `{db_path="testing.db"}`
			withTempFile(content, func(c *Config, tempFile string) {
				So(c.LoadFile(tempFile), ShouldNotBeNil)
				So(c.DbPath, ShouldContainSubstring, "development.db")
			})
		})

		Convey("should ignore invalid values present in the file", func() {
			content := `db_path1="testing.db"`
			withTempFile(content, func(c *Config, tempFile string) {
				So(c.DbPath, ShouldContainSubstring, "development.db")
				So(c.LoadFile(tempFile), ShouldBeNil)
				So(c.DbPath, ShouldContainSubstring, "development.db")
			})
		})
	})

	Convey("Load env", t, func() {

		Convey("should override default values", func() {
			os.Setenv("MTG_BIND_ADDR", ":4000")
			os.Setenv("MTG_DB_TYPE", "mysql")
			os.Setenv("MTG_DB_PATH", "another.db")

			c := New()
			So(c.LoadEnv(), ShouldBeNil)

			So(c.Addr, ShouldEqual, ":4000")
			So(c.DbPath, ShouldEqual, "another.db")
			So(c.DbType, ShouldEqual, "mysql")
		})

	})

}

func withTempFile(content string, fn func(c *Config, tempFile string)) {
	f, _ := ioutil.TempFile("", "")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	c := New()
	fn(c, f.Name())
}
