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
	Start     int     `form:"start"`
	Limit     int     `form:"limit"`
	K         float32 `form:"k"`
	D         float32 `form:"d"`
	D1        float32 `form:"d1"`
	D2        float32 `form:"d2"`
}

type Params map[string]interface{}

func (this URLOptionsScheme) SetDefaults() URLOptionsScheme {
	if this.Quality == 0 {
		this.Quality = 80
	}

	return this
}

func (this URLOptionsScheme) IsImgOptions() bool {
	if this.Scale != "" || this.Quality != 0 {
		return true
	}

	return false
}

// @todo rewrite it
func (this URLOptionsScheme) ToHash() string {
	var s string

	if this.Height != 0 {
		s += ".he" + strconv.FormatFloat(this.Height, 'f', 1, 64) + "-"
	}

	if this.Chest != 0 {
		if s == "" {
			s = "."
		}

		s += "ce" + strconv.FormatFloat(this.Chest, 'f', 1, 64) + "-"
	}

	if this.Underbust != 0 {
		if s == "" {
			s = "."
		}

		s += "un" + strconv.FormatFloat(this.Underbust, 'f', 1, 64) + "-"
	}

	if this.Waist != 0 {
		if s == "" {
			s = "."
		}

		s += "wa" + strconv.FormatFloat(this.Waist, 'f', 1, 64) + "-"
	}

	if this.Hips != 0 {
		if s == "" {
			s = "."
		}

		s += "hi" + strconv.FormatFloat(this.Hips, 'f', 1, 64) + "-"
	}

	return s
}

// @todo rewrite it
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

	if this.Underbust != 0 {
		result["underbust"] = this.Underbust
	}

	return result
}
