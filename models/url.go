package models

import (
	"strconv"
)

type URLOptionsScheme struct {
	Scale     string  `form:"scale"`
	Format    string  `form:"format"`
	Quality   int     `form:"q"`
	Height    float64 `form:"height"`
	Chest     float64 `form:"chest"`
	Underbust float64 `form:"underbust"`
	Waist     float64 `form:"waist"`
	Hips      float64 `form:"hips"`
}

type Params map[string]interface{}

func (this URLOptionsScheme) SetDefaults() URLOptionsScheme {
	if this.Format != "jpeg" || this.Format != "png" {
		this.Format = "jpeg"
	}

	if this.Quality == 0 {
		this.Quality = 80
	}

	return this
}

// returns string ".he{height}-ch{chest}-un{underbust}-hi{hips}"
func (this URLOptionsScheme) ToHash(p Params) string {
	var s string

	if len(p) != 0 {
		s = "."
	}

	for name, val := range p {
		s += name[0:2] + strconv.FormatFloat(val.(float64), 'f', 1, 64) + "-"
	}

	return s
}

func (this URLOptionsScheme) ToMap() Params {
	result := Params{}

	if this.Height != 0 {
		result["height"] = this.Height
	}

	if this.Chest != 0 {
		result["chest"] = this.Chest
	}

	if this.Waist != 0 {
		result["waist"] = this.Waist
	}

	if this.Hips != 0 {
		result["hips"] = this.Hips
	}

	return result
}
