package model

type Impuesto struct {
	Impuesto    uint   `json:"impuesto"`
	IDOrg       uint   `json:"idOrg"`
	TipoRegimen string `json:"tipoRegimen"`
	Nombre      string `json:"nombre"`
	FechaDesde  string `json:"fechaDesde"`
	FechaHasta  string `json:"fechaHasta"`
}
