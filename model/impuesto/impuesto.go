package impuesto

type Impuesto struct {
	ID          uint64 `json:"id,omitempty"`
	Org         uint64 `json:"org,omitempty"`
	Nombre      string `json:"nombre,omitempty"`
	Abreviatura string `json:"abreviatura,omitempty"`
}
