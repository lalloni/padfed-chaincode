package model

// Persona asset
type Persona struct {

	// common
	Tipo                      string                           `json:"tipo,omitempty"`
	ID                        uint64                           `json:"id"`
	TipoID                    string                           `json:"tipoid"`
	ActivoID                  uint                             `json:"activoid,omitempty"`
	Estado                    string                           `json:"estado,omitempty"`
	Pais                      uint                             `json:"pais,omitempty"`
	Inscripcion               *Inscripcion                     `json:"inscripcion,omitempty"`
	CH                        []string                         `json:"ch,omitempty"`
	DS                        string                           `json:"ds,omitempty"`
	Actividades               map[string]Actividad             `json:"actividades,omitempty"`
	Impuestos                 map[string]Impuesto              `json:"impuestos,omitempty"`
	Domicilios                map[string]Domicilio             `json:"domicilios,omitempty"`
	Telefonos                 map[string]Telefono              `json:"telefonos,omitempty"`
	Jurisdicciones            map[string]Jurisdiccion          `json:"jurisdicciones,omitempty"`
	Emails                    map[string]Email                 `json:"emails,omitempty"`
	Archivos                  map[string]Archivo               `json:"archivos,omitempty"`
	Categorias                map[string]Categoria             `json:"categorias,omitempty"`
	Etiquetas                 map[string]Etiqueta              `json:"etiquetas,omitempty"`
	ContribucionesMunicipales map[string]ContribucionMunicipal `json:"contribmunis,omitempty"`

	// fisica
	Nombre        string     `json:"nombre,omitempty"`
	Apellido      string     `json:"apellido,omitempty"`
	Materno       string     `json:"materno,omitempty"`
	Sexo          string     `json:"sexo,omitempty"`
	Documento     *Documento `json:"documento,omitempty"`
	Nacimiento    *Fecha     `json:"nacimiento,omitempty"`
	Fallecimiento *Fecha     `json:"fallecimiento,omitempty"`

	// juridica
	RazonSocial    string              `json:"razonsocial,omitempty"`
	FormaJuridica  uint                `json:"formajuridica,omitempty"`
	MesCierre      uint                `json:"mescierre,omitempty"`
	ContratoSocial *Fecha              `json:"contratosocial,omitempty"`
	Relaciones     map[string]Relacion `json:"relaciones,omitempty"`
}

type Actividad struct {
	Actividad string `json:"actividad,omitempty"`
	Orden     uint   `json:"prden,omitempty"`
	Periodo   uint   `json:"periodo,omitempty"`
	Estado    string `json:"estado,omitempty"`
	DS        *Fecha `json:"ds,omitempty"`
}

type Domicilio struct {
	Nombre      string     `json:"nombre,omitempty"`
	Orden       uint       `json:"orden,omitempty"`
	Tipo        uint       `json:"tipo,omitempty"`
	Estado      uint       `json:"estado,omitempty"`
	Calle       string     `json:"calle,omitempty"`
	Numero      uint       `json:"numero,omitempty"`
	Piso        string     `json:"piso,omitempty"`
	Sector      string     `json:"sector,omitempty"`
	Manzana     string     `json:"manzana,omitempty"`
	Torre       string     `json:"torre,omitempty"`
	Unidad      string     `json:"unidad,omitempty"`
	Provincia   *uint      `json:"provincia,omitempty"` // puntero porque 0 es un valor de negocio
	Localidad   string     `json:"localidad,omitempty"`
	CP          string     `json:"cp,omitempty"`
	Nomenclador string     `json:"nomenclador,omitempty"`
	Adicional   *Adicional `json:"adicional,omitempty"`
	Nota        string     `json:"nota,omitempty"`
	Baja        *Fecha     `json:"baja,omitempty"`
	DS          *Fecha     `json:"ds,omitempty"`
}

type Telefono struct {
	Orden  uint   `json:"orden,omitempty"`
	Pais   uint   `json:"pais,omitempty"`
	Area   uint   `json:"area,omitempty"`
	Numero uint   `json:"numero,omitempty"`
	Tipo   uint   `json:"tipo,omitempty"`
	Linea  uint   `json:"linea,omitempty"`
	Estado uint   `json:"estado,omitempty"`
	DS     *Fecha `json:"ds,omitempty"`
}

type Adicional struct {
	Tipo uint   `json:"tipo,omitempty"`
	Dato string `json:"dato,omitempty"`
}

type Inscripcion struct {
	Registro uint `json:"registro,omitempty"`
	Numero   uint `json:"numero,omitempty"`
}
type Documento struct {
	Tipo   uint   `json:"tipo,omitempty"`
	Numero string `json:"numero,omitempty"`
}

type Jurisdiccion struct {
	Provincia *uint  `json:"provincia,omitempty"`
	Desde     *Fecha `json:"desde,omitempty"`
	Hasta     *Fecha `json:"hasta,omitempty"`
	Sede      *bool  `json:"sede,omitempty"`
	DS        *Fecha `json:"ds,omitempty"`
}

// Impuesto asset
type Impuesto struct {
	Impuesto    uint   `json:"impuesto,omitempty"`
	Inscripcion Fecha  `json:"inscripcion,omitempty"`
	Estado      string `json:"estado"`
	Dia         uint   `json:"dia,omitempty"`
	Periodo     uint   `json:"periodo"`
	Motivo      uint   `json:"motivo,omitempty"`
	DS          Fecha  `json:"ds,omitempty"`
}

type Email struct {
	Direccion string `json:"direccion,ompitempty"`
	Orden     uint   `json:"orden,ompitempty"`
	Tipo      uint   `json:"tipo,ompitempty"`
	Estado    uint   `json:"estado,ompitempty"`
	DS        Fecha  `json:"ds,ompitempty"`
}

type Archivo struct {
	Descripcion string `json:"descripcion,omitempty"`
	Orden       uint   `json:"orden,omitempty"`
	Tipo        uint   `json:"tipo,omitempty"`
	DS          Fecha  `json:"ds,omitempty"`
}

type Categoria struct {
	Categoria uint   `json:"categoria,omitempty"`
	Motivo    uint   `json:"motivo,omitempty"`
	Estado    string `json:"estado,omitempty"`
	Impuesto  uint   `json:"impuesto,omitempty"`
	Periodo   uint   `json:"periodo,omitempty"`
	DS        Fecha  `json:"ds,omitempty"`
}

type Etiqueta struct {
	Etiqueta uint   `json:"etiqueta,omitempty"`
	Periodo  uint   `json:"periodo,omitempty"`
	Estado   string `json:"estado,omitempty"`
	DS       Fecha  `json:"ds,omitempty"`
}

type ContribucionMunicipal struct {
	Impuesto  uint  `json:"impuesto,omitempty"`
	Municipio uint  `json:"municipio,omitempty"`
	Provincia *uint `json:"provincia,omitempty"`
	Desde     Fecha `json:"desde,omitempty"`
	Hasta     Fecha `json:"hasta,omitempty"`
	DS        Fecha `json:"ds,omitempty"`
}

type Relacion struct {
	Persona uint  `json:"persona,omitempty"`
	Tipo    uint  `json:"tipo,omitempty"`
	Subtipo uint  `json:"subtipo,omitempty"`
	Desde   Fecha `json:"desde,omitempty"`
	DS      Fecha `json:"ds,omitempty"`
}

type Personas struct {
	Personas []Persona `json:"personas"`
}

// ParamImpuesto asset
type ParamImpuesto struct {
	Impuesto    uint   `json:"impuesto"`
	IDOrg       uint   `json:"idOrg"`
	TipoRegimen string `json:"tipoRegimen"`
	Nombre      string `json:"nombre"`
	FechaDesde  string `json:"fechaDesde"`
	FechaHasta  string `json:"fechaHasta"`
}
