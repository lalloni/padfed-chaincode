package main

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestValidPersonaProto(t *testing.T) {

	persona1 := &Persona{
		CUIT:     30679638943,
		Nombre:   "Pepe",
		Apellido: "Sanchez",
	}
	persona2 := &Persona{
		CUIT:     30679638943,
		Nombre:   "Sancho",
		Apellido: "Panza",
	}

	personas := Personas{}
	personas.Personas = []*Persona{persona1, persona2}

	proto.RegisterType((*Impuesto)(nil), "main.Impuesto")

	data, err := proto.Marshal(&personas)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	var personasNew Personas
	err = proto.Unmarshal(data, &personasNew)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	fmt.Println("Persona nombre " + strconv.Itoa(len(personasNew.Personas)))
}
