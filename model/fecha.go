package model

import (
	"strings"
	"time"
)

const layout = "2006-01-02"

type Fecha time.Time

func (f *Fecha) UnmarshalJSON(bs []byte) error {
	p, err := time.Parse(layout, strings.Trim(string(bs), `"`))
	if err != nil {
		return err
	}
	*f = Fecha(p)
	return nil
}

func (f *Fecha) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*f).Format(layout) + `"`), nil
}
