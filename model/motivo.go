package model

type Motivo struct {
	ID    uint   `json:"id"`
	Desde *Fecha `json:"desde"`
	Hasta *Fecha `json:"hasta"`
}
