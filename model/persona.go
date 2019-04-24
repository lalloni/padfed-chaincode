package model

func NewPersona() *Persona {
	return &Persona{
		Actividades:    map[string]*PersonaActividad{},
		Impuestos:      map[string]*PersonaImpuesto{},
		Domicilios:     map[string]*PersonaDomicilio{},
		Telefonos:      map[string]*PersonaTelefono{},
		Jurisdicciones: map[string]*PersonaJurisdiccion{},
		Emails:         map[string]*PersonaEmail{},
		Archivos:       map[string]*PersonaArchivo{},
		Categorias:     map[string]*PersonaCategoria{},
		Etiquetas:      map[string]*PersonaEtiqueta{},
		Contribuciones: map[string]*PersonaContribucion{},
		Relaciones:     map[string]*PersonaRelacion{},
	}
}

type Persona struct {

	// siempre requerido
	ID uint64 `json:"id"`

	// atributos de persona
	Persona *PersonaBasica `json:"persona,omitempty"`

	// colecciones
	Actividades    map[string]*PersonaActividad    `json:"actividades,omitempty"`
	Impuestos      map[string]*PersonaImpuesto     `json:"impuestos,omitempty"`
	Domicilios     map[string]*PersonaDomicilio    `json:"domicilios,omitempty"`
	Telefonos      map[string]*PersonaTelefono     `json:"telefonos,omitempty"`
	Jurisdicciones map[string]*PersonaJurisdiccion `json:"jurisdicciones,omitempty"`
	Emails         map[string]*PersonaEmail        `json:"emails,omitempty"`
	Archivos       map[string]*PersonaArchivo      `json:"archivos,omitempty"`
	Categorias     map[string]*PersonaCategoria    `json:"categorias,omitempty"`
	Etiquetas      map[string]*PersonaEtiqueta     `json:"etiquetas,omitempty"`
	Contribuciones map[string]*PersonaContribucion `json:"contribmunis,omitempty"`
	Relaciones     map[string]*PersonaRelacion     `json:"relaciones,omitempty"`
}

type PersonaBasica struct {

	// atributos comunes
	Tipo        string       `json:"tipo,omitempty"`
	ID          uint64       `json:"id,omitempty"`
	TipoID      string       `json:"tipoid,omitempty"`
	ActivoID    uint         `json:"activoid,omitempty"`
	Estado      string       `json:"estado,omitempty"`
	Pais        uint         `json:"pais,omitempty"`
	Inscripcion *Inscripcion `json:"inscripcion,omitempty"`
	CH          []string     `json:"ch,omitempty"`
	DS          string       `json:"ds,omitempty"`

	// atributos de persona fisica
	Nombre        string     `json:"nombre,omitempty"`
	Apellido      string     `json:"apellido,omitempty"`
	Materno       string     `json:"materno,omitempty"`
	Sexo          string     `json:"sexo,omitempty"`
	Documento     *Documento `json:"documento,omitempty"`
	Nacimiento    *Fecha     `json:"nacimiento,omitempty"`
	Fallecimiento *Fecha     `json:"fallecimiento,omitempty"`

	// atributos de persona juridica
	RazonSocial    string `json:"razonsocial,omitempty"`
	FormaJuridica  uint   `json:"formajuridica,omitempty"`
	MesCierre      uint   `json:"mescierre,omitempty"`
	ContratoSocial *Fecha `json:"contratosocial,omitempty"`
}

type Inscripcion struct {
	Registro uint `json:"registro,omitempty"`
	Numero   uint `json:"numero,omitempty"`
}

type Documento struct {
	Tipo   uint   `json:"tipo,omitempty"`
	Numero string `json:"numero,omitempty"`
}

type PersonaActividad struct {
	Actividad string `json:"actividad,omitempty"`
	Orden     uint   `json:"prden,omitempty"`
	Periodo   uint   `json:"periodo,omitempty"`
	Estado    string `json:"estado,omitempty"`
	DS        *Fecha `json:"ds,omitempty"`
}

type PersonaDomicilio struct {
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
	Provincia   *Provincia `json:"provincia,omitempty"` // 0 es un valor de negocio
	Localidad   string     `json:"localidad,omitempty"`
	CP          string     `json:"cp,omitempty"`
	Nomenclador string     `json:"nomenclador,omitempty"`
	Adicional   *Adicional `json:"adicional,omitempty"`
	Nota        string     `json:"nota,omitempty"`
	Baja        *Fecha     `json:"baja,omitempty"`
	DS          *Fecha     `json:"ds,omitempty"`
}

type Adicional struct {
	Tipo uint   `json:"tipo,omitempty"`
	Dato string `json:"dato,omitempty"`
}

type PersonaTelefono struct {
	Orden  uint   `json:"orden,omitempty"`
	Pais   uint   `json:"pais,omitempty"`
	Area   uint   `json:"area,omitempty"`
	Numero uint   `json:"numero,omitempty"`
	Tipo   uint   `json:"tipo,omitempty"`
	Linea  uint   `json:"linea,omitempty"`
	Estado uint   `json:"estado,omitempty"`
	DS     *Fecha `json:"ds,omitempty"`
}

type PersonaJurisdiccion struct {
	Provincia *Provincia `json:"provincia,omitempty"` // 0 tiene valor de negocio
	Desde     *Fecha     `json:"desde,omitempty"`
	Hasta     *Fecha     `json:"hasta,omitempty"`
	Sede      bool       `json:"sede,omitempty"`
	DS        *Fecha     `json:"ds,omitempty"`
}

type PersonaImpuesto struct {
	Impuesto    uint   `json:"impuesto,omitempty"`
	Inscripcion *Fecha `json:"inscripcion,omitempty"`
	Estado      string `json:"estado,omitempty"`
	Dia         uint   `json:"dia,omitempty"`
	Periodo     uint   `json:"periodo,omitempty"`
	Motivo      uint   `json:"motivo,omitempty"`
	DS          *Fecha `json:"ds,omitempty"`
}

type PersonaEmail struct {
	Direccion string `json:"direccion,omitempty"`
	Orden     uint   `json:"orden,omitempty"`
	Tipo      uint   `json:"tipo,omitempty"`
	Estado    uint   `json:"estado,omitempty"`
	DS        *Fecha `json:"ds,omitempty"`
}

type PersonaArchivo struct {
	Descripcion string `json:"descripcion,omitempty"`
	Orden       uint   `json:"orden,omitempty"`
	Tipo        uint   `json:"tipo,omitempty"`
	DS          *Fecha `json:"ds,omitempty"`
}

type PersonaCategoria struct {
	Categoria uint   `json:"categoria,omitempty"`
	Motivo    uint   `json:"motivo,omitempty"`
	Estado    string `json:"estado,omitempty"`
	Impuesto  uint   `json:"impuesto,omitempty"`
	Periodo   uint   `json:"periodo,omitempty"`
	DS        *Fecha `json:"ds,omitempty"`
}

type PersonaEtiqueta struct {
	Etiqueta uint   `json:"etiqueta,omitempty"`
	Periodo  uint   `json:"periodo,omitempty"`
	Estado   string `json:"estado,omitempty"`
	DS       *Fecha `json:"ds,omitempty"`
}

type PersonaContribucion struct {
	Impuesto  uint       `json:"impuesto,omitempty"`
	Municipio uint       `json:"municipio,omitempty"`
	Provincia *Provincia `json:"provincia,omitempty"` // 0 tiene valor de negocio
	Desde     *Fecha     `json:"desde,omitempty"`
	Hasta     *Fecha     `json:"hasta,omitempty"`
	DS        *Fecha     `json:"ds,omitempty"`
}

type PersonaRelacion struct {
	Persona uint   `json:"persona,omitempty"`
	Tipo    uint   `json:"tipo,omitempty"`
	Subtipo uint   `json:"subtipo,omitempty"`
	Desde   *Fecha `json:"desde,omitempty"`
	DS      *Fecha `json:"ds,omitempty"`
}
