// MIT License
//
// Copyright (c) 2019 Pablo Ignacio Lalloni
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cuit

import "strconv"

const (
	// PersonaFísica es el tipo de las personas físicas ("humanas")
	PersonaFísica TipoPersona = iota

	// PersonaJurídica es el tipo de las personas jurídicas
	PersonaJurídica
)

// tipopersona es un tipo privado para impedir creación de TipoPersona fuera de este paquete
type tipopersona uint8

// TipoPersona es un enumerado de los tipos de persona
type TipoPersona tipopersona

func (t TipoPersona) String() string {
	switch t {
	case PersonaFísica:
		return "Persona Física"
	case PersonaJurídica:
		return "Persona Jurídica"
	default:
		return "TipoPersona(" + strconv.Itoa(int(t)) + ")"
	}
}

const minjurídica = 3e10

// TipoPersonaCUIT retorna el TipoPersona que corresponde al cuit suministrado según su rango
func TipoPersonaCUIT(cuit uint64) TipoPersona {
	if cuit < minjurídica {
		return PersonaFísica
	}
	return PersonaJurídica
}
