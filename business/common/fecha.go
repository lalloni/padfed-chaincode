package common

import (
	"reflect"
	"strings"
	"time"
)

const layout = "2006-01-02"

type Fecha time.Time

func (f *Fecha) UnmarshalJSON(bs []byte) error {
	p, err := time.ParseInLocation(layout, strings.Trim(string(bs), `"`), time.Local)
	if err != nil {
		return err
	}
	*f = Fecha(p)
	return nil
}

func (f *Fecha) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*f).Format(layout) + `"`), nil
}

func FechaHoy() *Fecha {
	d := Fecha(time.Now().Truncate(24 * time.Hour))
	return &d
}

func FechaEn(año, mes, día int) *Fecha {
	d := Fecha(time.Date(año, time.Month(mes), día, 0, 0, 0, 0, time.Local))
	return &d
}

func FechaDecodeHook(f, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}
	if t != reflect.TypeOf(Fecha{}) {
		return data, nil
	}
	d, err := time.Parse(layout, data.(string))
	if err != nil {
		return nil, err
	}
	return Fecha(d), nil
}
