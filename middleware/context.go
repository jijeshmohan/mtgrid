package middleware

import (
	"github.com/go-martini/martini"
)

type Data map[string]interface{}

func InitContext() martini.Handler {
	return func(c martini.Context) {
		data := make(Data)
		c.Map(data)
		c.Next()
	}
}
