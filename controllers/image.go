package controllers

type Image struct{}

func (*Image) Construct(args ...interface{}) interface{} {
	return &Image{}
}
