package personas

import (
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
)

type Persona struct {

	// siempre requerido
	ID uint64 `json:"id"`

	// atributos de persona
	Persona *Basica `json:"persona,omitempty"`

	// colecciones
	Actividades     map[string]*Actividad    `json:"actividades,omitempty"`
	Impuestos       map[string]*Impuesto     `json:"impuestos,omitempty"`
	Domicilios      map[string]*Domicilio    `json:"domicilios,omitempty"`
	DomiciliosRoles map[string]*DomicilioRol `json:"domisroles,omitempty"`
	Telefonos       map[string]*Telefono     `json:"telefonos,omitempty"`
	Jurisdicciones  map[string]*Jurisdiccion `json:"jurisdicciones,omitempty"`
	Emails          map[string]*Email        `json:"emails,omitempty"`
	Archivos        map[string]*Archivo      `json:"archivos,omitempty"`
	Categorias      map[string]*Categoria    `json:"categorias,omitempty"`
	Etiquetas       map[string]*Etiqueta     `json:"etiquetas,omitempty"`
	Contribuciones  map[string]*Contribucion `json:"contribmunis,omitempty"`
	Relaciones      map[string]*Relacion     `json:"relaciones,omitempty"`
	CMSedes         map[string]*CMSede       `json:"cmsedes,omitempty"`

	// errores de lectura de singletons y collections si el store está en modo lenient
	Errors interface{} `json:"errors,omitempty"`
}

type Basica struct {

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
	Nombre        string        `json:"nombre,omitempty"`
	Apellido      string        `json:"apellido,omitempty"`
	Materno       string        `json:"materno,omitempty"`
	Sexo          string        `json:"sexo,omitempty"`
	Documento     *Documento    `json:"documento,omitempty"`
	Nacimiento    *common.Fecha `json:"nacimiento,omitempty"`
	Fallecimiento *common.Fecha `json:"fallecimiento,omitempty"`

	// atributos de persona juridica
	RazonSocial    string        `json:"razonsocial,omitempty"`
	FormaJuridica  uint          `json:"formajuridica,omitempty"`
	MesCierre      uint          `json:"mescierre,omitempty"`
	ContratoSocial *common.Fecha `json:"contratosocial,omitempty"`
}

type Inscripcion struct {
	Registro uint `json:"registro,omitempty"`
	Numero   uint `json:"numero,omitempty"`
}

type Documento struct {
	Tipo   uint   `json:"tipo,omitempty"`
	Numero string `json:"numero,omitempty"`
}

type Actividad struct {
	Org       uint          `json:"org,omitempty"`
	Actividad string        `json:"actividad,omitempty"`
	Orden     uint          `json:"orden,omitempty"`
	Articulo  uint          `json:"articulo,omitempty"`
	Desde     *common.Fecha `json:"desde,omitempty"`
	Hasta     *common.Fecha `json:"hasta,omitempty"`
	DS        *common.Fecha `json:"ds,omitempty"`
}

type Domicilio struct {
	Nombre      string            `json:"nombre,omitempty"`
	Orden       uint              `json:"orden,omitempty"`
	Org         uint              `json:"org,omitempty"`
	Tipo        uint              `json:"tipo,omitempty"`
	Estado      uint              `json:"estado,omitempty"`
	Calle       string            `json:"calle,omitempty"`
	Numero      uint              `json:"numero,omitempty"`
	Piso        string            `json:"piso,omitempty"`
	Sector      string            `json:"sector,omitempty"`
	Manzana     string            `json:"manzana,omitempty"`
	Torre       string            `json:"torre,omitempty"`
	Unidad      string            `json:"unidad,omitempty"`
	Provincia   *common.Provincia `json:"provincia,omitempty"` // 0 es un valor de negocio
	Localidad   string            `json:"localidad,omitempty"`
	CP          string            `json:"cp,omitempty"`
	Nomenclador string            `json:"nomenclador,omitempty"`
	Adicional   *Adicional        `json:"adicional,omitempty"`
	Baja        *common.Fecha     `json:"baja,omitempty"`
	DS          *common.Fecha     `json:"ds,omitempty"`
}

type DomicilioRol struct {
	Orden uint          `json:"orden,omitempty"`
	Org   uint          `json:"org,omitempty"`
	Tipo  uint          `json:"tipo,omitempty"`
	Rol   uint          `json:"rol,omitempty"`
	DS    *common.Fecha `json:"ds,omitempty"`
}

type Adicional struct {
	Tipo uint   `json:"tipo,omitempty"`
	Dato string `json:"dato,omitempty"`
}

type Telefono struct {
	Orden  uint          `json:"orden,omitempty"`
	Pais   uint          `json:"pais,omitempty"`
	Area   uint          `json:"area,omitempty"`
	Numero uint          `json:"numero,omitempty"`
	Tipo   uint          `json:"tipo,omitempty"`
	Linea  uint          `json:"linea,omitempty"`
	Estado uint          `json:"estado,omitempty"`
	DS     *common.Fecha `json:"ds,omitempty"`
}

type Jurisdiccion struct {
	Provincia *common.Provincia `json:"provincia,omitempty"` // 0 tiene valor de negocio
	Desde     *common.Fecha     `json:"desde,omitempty"`
	Hasta     *common.Fecha     `json:"hasta,omitempty"`
	Org       uint              `json:"org,omitempty"`
	DS        *common.Fecha     `json:"ds,omitempty"`
}

type Impuesto struct {
	Impuesto    uint           `json:"impuesto,omitempty"`
	Inscripcion *common.Fecha  `json:"inscripcion,omitempty"`
	Estado      string         `json:"estado,omitempty"`
	Dia         uint           `json:"dia,omitempty"`
	Periodo     uint           `json:"periodo,omitempty"`
	Motivo      *common.Motivo `json:"motivo,omitempty"`
	DS          *common.Fecha  `json:"ds,omitempty"`
}

type Email struct {
	Direccion string        `json:"direccion,omitempty"`
	Orden     uint          `json:"orden,omitempty"`
	Tipo      uint          `json:"tipo,omitempty"`
	Estado    uint          `json:"estado,omitempty"`
	DS        *common.Fecha `json:"ds,omitempty"`
}

type Archivo struct {
	Descripcion string        `json:"descripcion,omitempty"`
	Orden       uint          `json:"orden,omitempty"`
	Tipo        uint          `json:"tipo,omitempty"`
	DS          *common.Fecha `json:"ds,omitempty"`
}

type Categoria struct {
	Categoria uint          `json:"categoria,omitempty"`
	Motivo    uint          `json:"motivo,omitempty"`
	Estado    string        `json:"estado,omitempty"`
	Impuesto  uint          `json:"impuesto,omitempty"`
	Periodo   uint          `json:"periodo,omitempty"`
	DS        *common.Fecha `json:"ds,omitempty"`
}

type Etiqueta struct {
	Etiqueta uint          `json:"etiqueta,omitempty"`
	Periodo  uint          `json:"periodo,omitempty"`
	Estado   string        `json:"estado,omitempty"`
	DS       *common.Fecha `json:"ds,omitempty"`
}

type Contribucion struct {
	Impuesto  uint              `json:"impuesto,omitempty"`
	Municipio uint              `json:"municipio,omitempty"`
	Provincia *common.Provincia `json:"provincia,omitempty"` // 0 tiene valor de negocio
	Desde     *common.Fecha     `json:"desde,omitempty"`
	Hasta     *common.Fecha     `json:"hasta,omitempty"`
	DS        *common.Fecha     `json:"ds,omitempty"`
}

type Relacion struct {
	Persona uint          `json:"persona,omitempty"`
	Tipo    uint          `json:"tipo,omitempty"`
	Subtipo uint          `json:"subtipo,omitempty"`
	Desde   *common.Fecha `json:"desde,omitempty"`
	DS      *common.Fecha `json:"ds,omitempty"`
}

type CMSede struct {
	Provincia *common.Provincia `json:"provincia,omitempty"` // 0 tiene valor de negocio
	Desde     *common.Fecha     `json:"desde,omitempty"`
	Hasta     *common.Fecha     `json:"hasta,omitempty"`
	DS        *common.Fecha     `json:"ds,omitempty"`
}
