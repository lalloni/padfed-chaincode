package model

type Organismo uint

func OrganismoCódigo(c uint) *Organismo {
	o := Organismo(c)
	return &o
}
