package impuestos

type Impuesto struct {
	Codigo      uint64 `json:"codigo,omitempty"`
	Org         uint64 `json:"org,omitempty"`
	Nombre      string `json:"nombre,omitempty"`
	Abreviatura string `json:"abreviatura,omitempty"`
}
