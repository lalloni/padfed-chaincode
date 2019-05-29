package model

type Motivo struct {
	ID    uint   `json:"id,omitempty"`
	Desde *Fecha `json:"desde,omitempty"`
	Hasta *Fecha `json:"hasta,omitempty"`
}
