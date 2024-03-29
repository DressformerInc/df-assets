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
	IsBody       bool          `gorethink:"is_body,omitempty"       json:"is_body,omitempty"`
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

func (this *Geometry) FindAll(ids []string, opts URLOptionsScheme) []GeometryScheme {
	var query r.Term

	if opts.Limit == 0 || opts.Limit > 100 {
		opts.Limit = 25
	}

	if len(ids) > 0 {
		query = this.GetAll(r.Args(ids))
	} else {
		query = this.Skip(opts.Start).Limit(opts.Limit)
	}

	rows, err := query.Run(session())

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
//
// @todo
// This is just a prototype. Rewrite it!
//
func Morph(dstNames []string, geoms []GeometryScheme, pmap Params, options URLOptionsScheme) ([]byte, error) {
	bodySources := []*gomorph.Source{}
	params := map[string]float32{}

	dummyModel := (*Dummy).Construct(nil).(*Dummy)
	defaultDummy := dummyModel.Find("")
	defDummyGeometry := (*Geometry).Construct(nil).(*Geometry).Find(defaultDummy.Assets.Geometry.Id)

	bodySources = append(bodySources, &gomorph.Source{
		Input:     AppConfig.StorageFilePath(defDummyGeometry.Base.Id),
		Output:    dstNames[0],
		ParamName: "base",
		SrcWeight: float32(defaultDummy.Body.Height),
	})

	for name, val := range pmap {
		log.Println("Name:", name, "Val:", val)

		// make sources list for dummy

		sources := findSection(name, defDummyGeometry.MorphTargets)
		if sources == nil {
			return nil, errors.New(fmt.Sprintf("Section for parameter: %s not found", name))
		}

		if len(sources) == 0 {
			return nil, errors.New(fmt.Sprintf("Empty sources for section: %s", name))
		}

		if len(sources) != 2 {
			return nil, errors.New("Sources should contain 2 morphtargets.")
		}

		if name == "underbust" {
			name = "underchest"
		}

		params[name] = float32(val.(float64))

		for i := 0; i < 2; i++ {
			source := sources[i]

			fp := AppConfig.StorageFilePath(source.Id)
			if _, err := os.Stat(fp); err != nil {
				return nil, errors.New(fmt.Sprintf("One of morphtargets sources not found: %s\n", fp))
			}

			bodySources = append(bodySources, &gomorph.Source{
				Input:     fp,
				ParamName: name,
				SrcWeight: float32(source.Weight),
			})
		}
	}

	dummy := gomorph.NewDummy("default")
	dummy.AddMorphTargets(bodySources)

	if geoms[0].IsBody {
		dummy.Morph(dstNames[0], params)
	} else {
		garmentSources := []*gomorph.Source{}

		for i, name := range dstNames {
			garmentSources = append(garmentSources, &gomorph.Source{
				Input:  AppConfig.StorageFilePath(geoms[i].Base.Id),
				Output: name})
		}

		dummy.Morph("", params)
		dummy.PutOn(params, garmentSources)
	}

	dummy.Release()

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
