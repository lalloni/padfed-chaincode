package model

// Persona asset
type Persona struct {
	CUIT          uint64      `json:"cuit"`
	Nombre        string      `json:"nombre,omitempty"`
	Apellido      string      `json:"apellido,omitempty"`
	RazonSocial   string      `json:"razonSocial,omitempty"`
	Tipo          string      `json:"tipo,omitempty"`
	Estado        string      `json:"estado,omitempty"`
	FormaJuridica int32       `json:"formaJuridica,omitempty"`
	TipoDoc       int32       `json:"tipoDoc,omitempty"`
	Doc           string      `json:"doc,omitempty"`
	Sexo          string      `json:"sexo,omitempty"`
	MesCierre     int32       `json:"mesCierre,omitempty"`
	Nacimiento    string      `json:"nacimiento,omitempty"`
	Fallecimiento string      `json:"fallecimiento,omitempty"`
	Inscripcion   string      `json:"inscripcion,omitempty"`
	FechaCierre   string      `json:"fechaCierre,omitempty"`
	NuevaCUIT     uint64      `json:"nuevaCuit,omitempty"`
	Materno       string      `json:"materno,omitempty"`
	Pais          string      `json:"pais,omitempty"`
	CH            []string    `json:"ch,omitempty"`
	DS            string      `json:"ds,omitempty"`
	Impuestos     []*Impuesto `json:"impuestos,omitempty"`
}

// Impuesto asset
type Impuesto struct {
	Impuesto    int32  `json:"impuesto"`
	IDOrganismo int32  `json:"idOrg,omitempty"`
	Periodo     int32  `json:"periodo"`
	Dia         int32  `json:"dia,omitempty"`
	IDTxc       uint64 `json:"idTxc,omitempty"`
	Inscripcion string `json:"inscripcion,omitempty"`
	Estado      string `json:"estado"`
	DS          string `json:"ds,omitempty"`
	Motivo      string `json:"motivo,omitempty"`
}

type Personas struct {
	Personas []*Persona `json:"personas"`
}

// Impuestos es una estructura solo usada para poder extraer un array de impuestos.
type Impuestos struct {
	Impuestos []*Impuesto `json:"impuestos"`
}

// ParamImpuesto asset
type ParamImpuesto struct {
	Impuesto    int32  `json:"impuesto"`
	IDOrg       int32  `json:"idOrg"`
	TipoRegimen string `json:"tipoRegimen"`
	Nombre      string `json:"nombre"`
	FechaDesde  string `json:"fechaDesde"`
	FechaHasta  string `json:"fechaHasta"`
}
