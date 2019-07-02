package validator

type locale struct{}

func (l locale) Required() string {
	return `Se requiere {{.property}}`
}

func (l locale) InvalidType() string {
	return `Tipo inválido: debe ser {{.expected}} pero se recibió {{.given}}`
}

func (l locale) NumberAnyOf() string {
	return `Debe {{ if .subtitles }}ser uno de {{ range $i, $e := .subtitles }}{{ if $i }}, {{ end}}{{ . }}{{ end }}{{else}}cumplir al menos un sub esquema{{ end }}`
}

func (l locale) NumberOneOf() string {
	return `Debe {{ if .subtitles }}ser uno y sólo uno de {{ range $i, $e := .subtitles }}{{ if $i }}, {{ end}}{{ . }}{{ end }}{{ else }}cumplir un y sólo un sub esquema{{ end }}`
}

func (l locale) NumberAllOf() string {
	return `Debe {{ if .subtitles }}ser {{ range $i, $e := .subtitles }}{{ if $i }}, {{ end}}{{ . }}{{ end }}{{ else }}cumplir todos los sub esquemas{{ end }}`
}

func (l locale) NumberNot() string {
	return `No debe {{ if .subtitle }}ser {{ .subtitle }}{{ else }}cumplir el sub esquema{{ end }}`
}

func (l locale) MissingDependency() string {
	return `Depende de {{.dependency}}`
}

func (l locale) Internal() string {
	return `Error interno {{.error}}`
}

func (l locale) Const() string {
	return `{{.field}} no coincide con: {{.allowed}}`
}

func (l locale) Enum() string {
	return `{{.field}} debe ser uno de los siguientes: {{.allowed}}`
}

func (l locale) ArrayNoAdditionalItems() string {
	return `No se permiten items adicionales en el arreglo`
}

func (l locale) ArrayNotEnoughItems() string {
	return `No hay suficientes items en el arreglo para coincidir con la lista de esquemas`
}

func (l locale) ArrayMinItems() string {
	return `El arreglo debe tener al menos {{.min}} items`
}

func (l locale) ArrayMaxItems() string {
	return `El arreglo debe tener como máximo {{.max}} items`
}

func (l locale) Unique() string {
	return `items de {{.type}} deben ser únicos ({{.i}} y {{.j}} son iguales)`
}

func (l locale) ArrayContains() string {
	return `Al menos uno de los items debe coincidir`
}

func (l locale) ArrayMinProperties() string {
	return `Debe tener al menos {{.min}} propiedades`
}

func (l locale) ArrayMaxProperties() string {
	return `Debe tener como máximo {{.max}} propiedades`
}

func (l locale) AdditionalPropertyNotAllowed() string {
	return `La propiedad {{.property}} no está permitida`
}

func (l locale) InvalidPropertyPattern() string {
	return `La propiedad {{.property}} no coincide con el patrón '{{.pattern}}'`
}

func (l locale) InvalidPropertyName() string {
	return `El nombre de la propiedad "{{.property}}" no coincide`
}

func (l locale) StringGTE() string {
	return `La longitud debe ser mayor o igual que {{.min}}`
}

func (l locale) StringLTE() string {
	return `La longitud debe ser menor o igual que {{.max}}`
}

func (l locale) DoesNotMatchPattern() string {
	return `No coincide con el patrón "{{.pattern}}"`
}

func (l locale) DoesNotMatchFormat() string {
	return `No coincide con el formato {{.format}}`
}

func (l locale) MultipleOf() string {
	return `Debe ser un múltiplo de {{ .multiple.RatString }}`
}

func (l locale) NumberGTE() string {
	return `Debe ser mayor o igual que {{ .min.RatString }}`
}

func (l locale) NumberGT() string {
	return `Debe ser mayor que {{ .min.RatString }}`
}

func (l locale) NumberLTE() string {
	return `Debe ser menor o igual que {{ .max.RatString }}`
}

func (l locale) NumberLT() string {
	return `Debe ser menor que {{ .max.RatString }}`
}

// Schema validators
func (l locale) RegexPattern() string {
	return `Expresón regular inválida '{{.pattern}}'`
}

func (l locale) GreaterThanZero() string {
	return `{{.number}} debe ser mayor que 0`
}

func (l locale) MustBeOfA() string {
	return `{{.x}} debe ser de un {{.y}}`
}

func (l locale) MustBeOfAn() string {
	return `{{.x}} debe ser de un {{.y}}`
}

func (l locale) CannotBeUsedWithout() string {
	return `{{.x}} no puede ser usado sin {{.y}}`
}

func (l locale) CannotBeGT() string {
	return `{{.x}} no puede ser mayor que {{.y}}`
}

func (l locale) MustBeOfType() string {
	return `{{.key}} debe ser de tipo {{.type}}`
}

func (l locale) MustBeValidRegex() string {
	return `{{.key}} debe ser una expresión regular válida`
}

func (l locale) MustBeValidFormat() string {
	return `{{.key}} debe ser de un formato válido {{.given}}`
}

func (l locale) MustBeGTEZero() string {
	return `{{.key}} debe ser mayor o igual que 0`
}

func (l locale) KeyCannotBeGreaterThan() string {
	return `{{.key}} no puede ser mayor que {{.y}}`
}

func (l locale) KeyItemsMustBeOfType() string {
	return `items de {{.key}} deben ser {{.type}}`
}

func (l locale) KeyItemsMustBeUnique() string {
	return `items de {{.key}} deben ser únicos`
}

func (l locale) ReferenceMustBeCanonical() string {
	return `La referencia {{.reference}} debe ser canónica`
}

func (l locale) NotAValidType() string {
	return `Tiene tipo inválido. Se espera uno de: {{.expected}}, se recibió: {{.given}}`
}

func (l locale) Duplicated() string {
	return `El tipo {{.type}} está duplicado`
}

func (l locale) HttpBadStatus() string { // nolint: golint
	return `No se pudo leer esquema desde HTTP, el estado de la respuesta es {{.status}}`
}

func (l locale) ErrorFormat() string {
	return `{{.field}}: {{.description}}`
}

func (l locale) ParseError() string {
	return `Se espera: {{.expected}}, se recibió: Invalid JSON`
}

func (l locale) ConditionThen() string {
	return `Debe cumplir "{{ if .thentitle }}{{ .thentitle }}{{ else }}then{{ end }}" dado que se cumple "{{ if .iftitle }}{{ .iftitle }}{{ else }}if{{ end }}"`
}

func (l locale) ConditionElse() string {
	return `Debe cumplir "{{ if .elsetitle }}{{ .elsetitle }}{{ else }}else{{ end }}" dado que no se cumple "{{ if .iftitle }}{{ .iftitle }}{{ else }}if{{ end }}"`
}
