package models

import (
	. "df/assets/utils"
	"df/gomorph"
	"errors"
	"fmt"
	r "github.com/dancannon/gorethink"
	enc "github.com/dancannon/gorethink/encoding"
	"log"
	"os"
)

type Source struct {
	Id     string  `gorethink:"id"                  json:"id"`
	Weight float64 `gorethink:"weight,omitempty"    json:"weight,omitempty"`
	Name   string  `gorethink:"orig_name,omitempty" json:"orig_name,omitempty"`
}

type MorphTarget struct {
	Section string   `gorethink:"section,omitempty" json:"section,omitempty"`
	Sources []Source `gorethink:"sources,omitempty" json:"sources,omitempty"`
}

type GeometryScheme struct {
	Id           string        `gorethink:"id,omitempty"            json:"id" binding:"-"`
	Base         Source        `gorethink:"base,omitempty"          json:"base,omitempty"`
	Name         string        `gorethink:"name,omitempty"          json:"name,omitempty"`
	IsBody       bool          `gorethink:"is_body,omitempty"       json:"-"`
	MorphTargets []MorphTarget `gorethink:"morph_targets,omitempty" json:"morph_targets,omitempty"`
}

type Geometry struct {
	r.Term
}

func (*Geometry) Construct(args ...interface{}) interface{} {
	return &Geometry{
		r.Db("dressformer").Table("geometry"),
	}
}

func (this *Geometry) FindAll(opt URLOptionsScheme) []GeometryScheme {
	rows, err := this.Skip(opt.Start).Limit(opt.Limit).Run(session())

	if err != nil {
		log.Println("Unable to fetch cursor for all. Error:", err)
		return nil
	}

	var result []GeometryScheme

	if err = rows.All(&result); err != nil {
		log.Println("Unable to get data, err:", err)
	}

	return result
}

func (this *Geometry) Find(id interface{}) *GeometryScheme {
	rows, err := this.Get(id).Run(session())
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
	result, err := this.Insert(payload, r.InsertOpts{ReturnChanges: true}).Run(session())
	if err != nil {
		log.Println("Error inserting data:", err)
		return nil, errors.New("Internal server error")
	}

	response := &r.WriteResponse{}

	if err = result.One(response); err != nil {
		log.Println("Unable to iterate cursor:", err)
		return nil, errors.New("Internal server error")
	}

	if len(response.Changes) != 1 {
		log.Println("Unexpected length of Changes:", len(response.Changes))
		return nil, errors.New("Internal server error")
	}

	newval := &GeometryScheme{}

	if err = enc.Decode(newval, response.Changes[0].NewValue); err != nil {
		log.Println("Decode error:", err)
		return nil, errors.New("Internal server error")
	}

	return newval, nil
}

func (this *Geometry) Put(id string, payload GeometryScheme) (*GeometryScheme, error) {
	result, err := this.Get(id).Update(payload, r.UpdateOpts{ReturnChanges: true}).Run(session())
	if err != nil {
		log.Println("Error updating:", id, "with data:", payload, "error:", err)
		return nil, errors.New("Wrong data")
	}

	response := &r.WriteResponse{}

	if err = result.One(response); err != nil {
		log.Println("Unable to iterate cursor:", err)
		return nil, errors.New("Internal server error")
	}

	if len(response.Changes) != 1 {
		log.Println("Unexpected length of Changes:", len(response.Changes))
		return nil, errors.New("Internal server error")
	}

	newval := &GeometryScheme{}

	if err = enc.Decode(newval, response.Changes[0].NewValue); err != nil {
		log.Println("Decode error:", err)
		return nil, errors.New("Internal server error")
	}

	return newval, nil
}

func (this *Geometry) Remove(id string) error {
	_, err := this.Get(id).Delete().Run(session())
	if err != nil {
		log.Println("Error deleting:", id, "error:", err)
		return errors.New("Internal server error")
	}

	return nil
}

// @bug @todo
// While we're using GUID as a geomtry object id, we can't distribute
// assets. Because we can't store and get node id inside it.
// There is two ways:
// 1. change geometry id to oid
// 2. move geometry api to main api server and use common http assets interface
//    to get files
func (this *GeometryScheme) Morph(dst string, pmap Params, options URLOptionsScheme) ([]byte, error) {
	basefp := AppConfig.StorageFilePath(this.Base.Id)
	targets := []*gomorph.MorphTarget{}

	for name, val := range pmap {
		log.Println("Name:", name, "Val:", val)

		sources := findSection(name, this.MorphTargets)
		if sources == nil {
			return nil, errors.New(fmt.Sprintf("Section for parameter: %s not found", name))
		}

		if len(sources) == 0 {
			return nil, errors.New(fmt.Sprintf("Empty sources for section: %s", name))
		}

		mt := &gomorph.MorphTarget{
			DstWeight: float32(val.(float64)),
			Sources:   [2]*gomorph.Source{},
		}

		if len(sources) != 2 {
			return nil, errors.New("Sources should contain 2 morphtargets.")
		}

		for i := 0; i < 2; i++ {
			source := sources[i]

			fp := AppConfig.StorageFilePath(source.Id)
			if _, err := os.Stat(fp); err != nil {
				return nil, errors.New(fmt.Sprintf("One of morphtargets sources not found: %s\n", fp))
			}

			mt.Sources[i] = &gomorph.Source{
				File:      fp,
				SrcWeight: float32(source.Weight),
			}

		}
		targets = append(targets, mt)
	}

	if len(targets) > 0 {
		p := gomorph.Params{
			K:  options.K,
			D:  options.D,
			D1: options.D1,
			D2: options.D2,
		}

		if obj := gomorph.NewObjectFromSources(basefp, targets, p); obj == nil {
			return nil, errors.New("Server error")
		} else {
			b := obj.Blob(dst)
			log.Println("OK:", len(b))
			return b, nil
		}
	}

	return nil, nil
}

func findSection(name string, targets []MorphTarget) []Source {
	for _, target := range targets {
		if target.Section == name {
			return target.Sources
		}
	}

	return nil
}
