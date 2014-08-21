package models

import (
	"github.com/3d0c/martini-contrib/linker"
	rethink "github.com/dancannon/gorethink"
)

type Model interface {
	Construct(arg ...interface{}) interface{}
}

func session() *rethink.Session {
	return linker.Get().Session().(*rethink.Session)
}
