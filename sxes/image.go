package sxes

type Image struct {
	Uuid string
}

func MakeImage(uuid string) *Image {
	return &Image{
		Uuid: uuid,
	}
}
