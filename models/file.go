package models

type FileScheme struct {
	Id   string `json:"id"`
	Name string `json:"orig_name"`
	Err  string `json:"err_msg,omitempty"`
	Blob []byte `json:"-"`
	Size int    `json:"-"`
}

type File struct{}

func (*File) Construct(args ...interface{}) interface{} {
	return &File{}
}
