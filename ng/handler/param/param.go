package param

import "reflect"

type Parser func(args [][]byte) (v interface{}, consumed int, err error)

func New(name string, t reflect.Type, p Parser) *Param {
	return &Param{Name: name, Type: t, Parser: p}
}

type Param struct {
	Name   string
	Type   reflect.Type
	Parser Parser
}

func (p *Param) Parse(args [][]byte) (interface{}, int, error) {
	return p.Parser(args)
}

func (p *Param) Specialize(name string, next func(interface{}) (interface{}, error)) *Param {
	return &Param{
		Name: name,
		Type: p.Type,
		Parser: func(args [][]byte) (interface{}, int, error) {
			r, c, err := p.Parse(args)
			if err != nil {
				return nil, 0, err
			}
			r, err = next(r)
			if err != nil {
				return nil, 0, err
			}
			return r, c, nil
		},
	}
}
