package models

/*

#cgo CFLAGS: -Wall -I/usr/local/libmorph/include -std=c99 -Wunused-variable
#cgo LDFLAGS: -L/usr/local/libmorph/lib -lmorph

#include "proc.h"

*/
import "C"

import (
	. "df/assets/utils"
	"errors"
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
	"os"
	"unsafe"
)

type Source struct {
	Id     string  `gorethink:"id"     json:"id"`
	Weight float64 `gorethink:"weight" json:"weight"`
}

type MorphTarget struct {
	Section string   `gorethink:"section" json:"section"`
	Sources []Source `gorethink:"sources" json:"sources"`
}

type GeometryScheme struct {
	Id           string        `gorethink:"id,omitempty"   json:"id"   binding:"-"`
	Base         string        `gorethink:"base"           json:"base"`
	MorphTargets []MorphTarget `gorethink:"morph_targets"  json:"morph_targets,omitempty"`
}

type Geometry struct {
	table r.Term
}

func (*Geometry) Construct(args ...interface{}) interface{} {
	return &Geometry{
		table: r.Table("geometry"),
	}
}

func (this *Geometry) Find(id interface{}) *GeometryScheme {
	rows, err := this.table.Get(id).Run(session())
	if err != nil {
		log.Println("Unable to fetch cursor for id:", id, "Error:", err)
		return nil
	}

	var result *GeometryScheme

	if err = rows.One(&result); err != nil {
		log.Println("Unable to get data, err:", err)
		return nil
	}

	return result
}

func (this *Geometry) Create(payload GeometryScheme) (*GeometryScheme, error) {
	result, err := this.table.Insert(payload, r.InsertOpts{ReturnVals: true}).Run(session())
	if err != nil {
		log.Println("Error inserting data:", err)
		return nil, errors.New("Internal server error")
	}

	response := &r.WriteResponse{NewValue: &GeometryScheme{}}

	if err = result.One(response); err != nil {
		log.Println("Unable to iterate cursor:", err)
		return nil, errors.New("Internal server error")
	}

	return response.NewValue.(*GeometryScheme), nil
}

// func (this *GeometryScheme) Morph(dst string, pmap Params) error {
// 	return nil
// }

// @bug @todo
// While we're using GUID as a geomtry object id, we can't distribute
// assets. Because we can't store and get node id inside it.
// There is two ways:
// 1. change geometry id to oid
// 2. move geometry api to main api server and use common http assets interface
//    to get files

func (this *GeometryScheme) Morph(dst string, pmap Params) error {
	// if _, err := os.Stat(dst); err == nil || len(pmap) == 0 {
	// 	return nil
	// }

	c_src := C.CString(AppConfig.StorageFilePath(this.Base))
	c_dst := C.CString(dst)
	defer C.free(unsafe.Pointer(c_src))
	defer C.free(unsafe.Pointer(c_dst))

	mptr := C.initMorpher()
	mobj := C.initMobj(unsafe.Pointer(mptr))

	C.processAddBaseObject(c_src, unsafe.Pointer(mptr))

	i := 0

	for name, val := range pmap {
		sources := findSection(name, this.MorphTargets)
		if sources == nil {
			return errors.New(fmt.Sprintf("Section for parameter: %s not found", name))
		}

		if len(sources) == 0 {
			return errors.New(fmt.Sprintf("Empty sources for section: %s", name))
		}

		for _, source := range sources {
			fp := AppConfig.StorageFilePath(source.Id)
			if _, err := os.Stat(fp); err != nil {
				return errors.New(fmt.Sprintf("One of morphtargets sources not found: %s\n", fp))
			}

			log.Println("Adding source:", fp, source.Weight)

			c_fp := C.CString(fp)

			C.processAddMorphTargetObject(c_fp, C.size_t(i), C.double(source.Weight), unsafe.Pointer(mptr))

			C.free(unsafe.Pointer(c_fp))

		}

		C.procAddUid(unsafe.Pointer(mobj), C.int(i))
		C.procAddWeight(unsafe.Pointer(mobj), C.double(val.(float64)))

		i++
	}

	C.build(unsafe.Pointer(mobj))
	C.saveObject(c_dst, unsafe.Pointer(mobj))

	C.releaseMorpher(unsafe.Pointer(mptr))
	C.releaseMobj(unsafe.Pointer(mobj))

	return nil
}

func findSection(name string, targets []MorphTarget) []Source {
	for _, target := range targets {
		if target.Section == name {
			return target.Sources
		}
	}

	return nil
}
