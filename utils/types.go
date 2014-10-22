package utils

import (
	"path/filepath"
)

const (
	General = iota
	Jpg
	Png
)

type Type struct {
	flag        int
	format      string
	contentType string
	extensions  []string
	IsImage     bool
}

var types = []Type{
	{General, "", "application/octet-stream", []string{}, false},
	{Jpg, "jpg", "image/jpg", []string{"jpg", "jpeg"}, true},
	{Png, "png", "image/png", []string{"png"}, true},
}

func NewTypeFromExt(s string) Type {
	ext := filepath.Ext(s)

	if len(ext) > 0 {
		ext = ext[1:]
	}

	for i, _ := range types {
		if inSlice(ext, types[i].extensions) {
			return types[i]
		}
	}

	return types[0]
}

func NewTypeFromName(s string) Type {
	for i, _ := range types {
		if s == types[i].format {
			return types[i]
		}
	}

	return types[0]
}

func NewTypeFromFlag(in int) Type {
	for i, _ := range types {
		if in == types[i].flag {
			return types[i]
		}
	}

	return types[0]
}

func (this Type) Flag() int {
	return this.flag
}

func (this Type) ContentType() string {
	return this.contentType
}

func (this Type) Format() string {
	return this.format
}

func inSlice(s string, list []string) bool {
	for i, _ := range list {
		if s == list[i] {
			return true
		}
	}

	return false
}
