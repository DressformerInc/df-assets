package utils

import (
	"testing"
)

type Case struct {
	filename    string
	format      string
	contenttype string
	flag        int
}

var cases = []Case{
	{"xxx", "", "application/octet-stream", General},
	{"xxx.png", "png", "image/png", Png},
	{"xxx.jpg", "jpg", "image/jpg", Jpg},
}

func TestExt(t *testing.T) {
	for _, tcase := range cases {
		typ := NewTypeFromExt(tcase.filename)
		if typ.ContentType() != tcase.contenttype {
			t.Fatal("Wrong content type:", typ.ContentType(), "Expected:", tcase.contenttype)
		}

		if typ.Format() != tcase.format {
			t.Fatal("Wrong format:", typ.Format(), "Expected:", tcase.format)
		}
	}
}

func TestFlag(t *testing.T) {
	for _, tcase := range cases {
		typ := NewTypeFromFlag(tcase.flag)
		if typ.ContentType() != tcase.contenttype {
			t.Fatal("Wrong content type:", typ.ContentType(), "Expected:", tcase.contenttype)
		}

		if typ.Format() != tcase.format {
			t.Fatal("Wrong format:", typ.Format(), "Expected:", tcase.format)
		}
	}
}
