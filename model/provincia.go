package model

type Provincia uint

func ProvinciaCódigo(c uint) *Provincia {
	p := Provincia(c)
	return &p
}
