# padfed validator

## ¿Qué es esto?

Una librería Go para realizar validaciones de los datos que pasan a través de las API de padfed.

## Objetivos principales

- Ser usable dentro y fuera de padfed
- Ofrecer mensajes claros mostrables al usuario final y que permitan resolver los problemas de validación

## Principios de diseño

- Datos a validar (input) codificados en JSON
- Validaciones definidas en JSON-Schema
- Validación exhaustiva y devolución de reporte con todos los problemas

## Diseño de implementación

- Se utiliza YAML para representar internamente los JSON-Schema para aprovechar características ausentes en JSON:
  - Comentarios
  - Referencias internas
  - Humanidad!
- La librería es un módulo Go

## Uso básico

Desde un módulo Go ejecutar...

```sh
go get gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git
```

Y en el código...

Importar:

```go
import validator "gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator"
```

Instanciar:

```go
v := validator.New()
```

Ejecutar validación:

```go
res, err := v.ValidatePersonaJSON(persona)

if err != nil {
    // hubo un error en el proceso de validación (r no debe usarse)
    return errors.Wrap(err, "validando persona")
}
```

Usar resultados:

```go
if res.Valid() {
    // no se encontraron problemas de validación
    return nil
}
```

Informar problemas:

```go
for _, e := range res.Errors {
    // cada e es un problema de validación que contiene
    // el elemento problemático y la descripción del problema
    log.Printf("Error en %s: %s", e.Field, e.Description)
}
```

## Para desarrolladores de este proyecto

### Requisitos previos

- [Go](https://golang.org/dl/) v1.11.0 o posterior

#### Mage

El proyecto adopta [Mage](https://magefile.org/) como herramienta de ejecución de tareas pero no lo requiere preinstalado gracias a que al ser código Go se puede invocar directamente con `go run`.

Para explotar esta posibilidad se incorpora el archivo [mage.go](mage.go) en el raíz del proyecto que incluye el poquísimo código necesario para hacerlo.

De aquí en adelante se documentan todas las tareas utilizando esta característica.

Todos los comandos tendrán la forma:

```sh
go run mage.go ...
```

Por otro lado, si se [instaló Mage](https://magefile.org/#installation) en el entorno del usuario, se podrán ejecutar todas las tareas invocando directamente al mismo, para hacerlo de esta manera se debe reemplazar `go run mage.go` por `mage` en todos los comandos documentados.

En este caso, los comandos tendrán la forma:

```sh
mage ...
```

Siendo esta forma de ejecución marginalmente más eficiente.

### Tareas disponibles

Ver salida de:

```sh
go run mage.go
```

### Correr tests

Ejecuta los tests del proyecto informando fallos.

Ejecutar con:

```sh
go run mage.go test
```

### Lanzar GoConvey

GoConvey permite lanzar un proceso que:

- Sirve por HTTP una consola web de testing
- Monitorea cambios en archivos del proyecto
- Abre en el web browser la consola web

El proceso de monitoreo ejecuta automáticamente los tests ante los eventos de modificación de archivos del proyecto y notifica a la consola web para que ésta informe al usuario en tiempo real de todos los sucesos.

La consola web se integra con las notificaciones del desktop del usuario para mostrar avisos de resultados de las ejecuciones de los tests.

Y todo eso se obtiene ejecutando simplemente:

```sh
go run mage.go convey
```

### Release

Dado que el proyecto es una librería Go su proceso de release es:

1. Generación de recursos
2. Compilación
3. Ejecución de tests
4. Ejecución de linter
5. Creación del tag en Git
6. Publicación del nuevo tag a GitLab

```sh
env version=1.2.3 go run mage.go release
```
