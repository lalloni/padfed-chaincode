# API

## Generalidades

Todas las funciones (no deprecadas) del chaincode responden utilizando la estructura de respuesta de Fabric de la siguiente manera:

- Status: se utiliza un status code con [valores y semántica tomados de HTTP](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)
- Message: se utiliza para informar mensajes de error (si lo hubiera) mostrables a usuarios finales
- Payload: se utiliza para contener una estructura de respuesta estandarizada para todas las funciones (ver [Payload](#payload) a continuación)

### Payload

El Payload siempre tenrá la siguiente forma general:

```json
{
    "client": {},
    "chaincode": {},
    "transaction": {},
    "content": {},
    "content-encoding": "",
    "fault": {}
}
```

En el cual:

- `client`: es un objeto que contiene información relativa al sistema cliente que realizó la invocación (ver [Client](#client))
- `chaincode`: es un objeto que contiene información relativa al chaincode que recibió la invicación (ver [Chaincode](#chaincode))
- `transaction`: es un objeto que contiene información relativa a la transacción en la que se ejecutó la invocación (ver [Transaction](#transaction))
- `content`: es un objeto que contiene la respuesta a la operación
- `content-encoding`: es una cadena de caracteres que contiene el identificador del algoritmo de codificación que se utilizó para representar en `content` una secuencia arbitraria de bytes, si hubiera sido necesario (esto ocurre cuando los bytes no pueden representarse como una cadena de UTF-8 válida)
- `fault`: es un objeto o cadena de caracteres que contiene información acerca de la causa que impidió el éxito de la invocación, si este hubiera fallado

Condiciones:

- `fault` y `content` son mutuamente excluyentes
- `content-encoding` es requerido si se codificó `content` (actualmente solo puede tener el valor "base64" si el contenido de `content` no era expresable como una string UTF-8 y fue necesaria la codificación, caso contrario no estára presente indicando que no hubo codificación y que `content` contiene una cadena de caracteres UTF-8 literal)
- `client`, `chaincode` y `transaction` solo estarán presentes si se realizó la invocación incluyendo la opción `debug`, caso contrario no estarán presentes

#### Client

```json
{
    "mspid": "",
    "subject": "",
    "issuer": ""
}
```

Donde:

- `mspid`: contiene una cadena de caracteres con el identificador del MSP que autenticó al sistema cliente
- `subject`: contiene una [cadena de caracteres que representa](https://tools.ietf.org/html/rfc1779) el [X.500 DN](https://ldapwiki.com/wiki/Distinguished%20Names) del subject del certificado con el que se autenticó el sistema cliente
- `issuer`: contiene una [cadena de caracteres que representa](https://tools.ietf.org/html/rfc1779) el [X.500 DN](https://ldapwiki.com/wiki/Distinguished%20Names) del issuer del certificado con el que se autenticó el sistema cliente

#### Chaincode

```json
{
    "version": ""
}
```

Donde:

- `version`: es la string que identifica la versión del chaincode

#### Transaction

```json
{
    "id": "",
    "function": ""
}
```

Donde:

- `id`: contiene una cadena de caracteres con una representación del UUID que identifica la transacción
- `function`: contiene una cadena de caracteres con el nombre de la función invocada

## Funciones

| Nombre                              | ACL      | Argumentos                                                     | Payload                                                             |
| ----------------------------------- | -------- | -------------------------------------------------------------- | ------------------------------------------------------------------- |
| GetVersion                          | -        | -                                                              | version [String](#string)                                           |
| GetFunctions                        | -        | -                                                              | functions [JSON](#json) string array                                |
| GetPersona                          | -        | cuit [DecimalString](#decimalstring)                           | [Persona](#persona) [JSON](#json)                                   |
| GetPersonaBasica                    | -        | cuit [DecimalString](#decimalstring)                           | PersonaBasica [JSON](#json) object                                  |
| GetPersona[{Collection}](#cols)     | -        | cuit [DecimalString](#decimalstring)                           | [JSON](#json) array                                                 |
| GetPersona[{Collection}](#cols)Item | -        | cuit [DecimalString](#decimalstring), itemid [String](#string) | [JSON](#json)                                                       |
| DelPersona                          | OnlyAFIP | cuit [DecimalString](#decimalstring)                           | -                                                                   |
| PutPersona                          | OnlyAFIP | [persona](#persona) [JSON](#json)                              | -                                                                   |
| PutPersonaList                      | OnlyAFIP | [personalist](#personalist) [JSON](#json)                      | cantidad [DecimalString](#decimalstring)                            |
| GetPersonaRange                     | OnlyAFIP | cuit1, cuit2 [DecimalString](#decimalstring)                   | [Persona](#persona) [JSON](#json) array                             |
| DelPersonaRange                     | OnlyAFIP | cuit1, cuit2 [DecimalString](#decimalstring)                   | ids [JSON](#json) number array                                      |
| GetPersonaAll                       | OnlyAFIP | -                                                              | [Persona](#persona) [JSON](#json) array                             |
| GetImpuesto                         | -        | cuit [DecimalString](#decimalstring)                           | [Impuesto](#impuesto) [JSON](#json)                                 |
| DelImpuesto                         | OnlyAFIP | codigoimpuesto [DecimalString](#decimalstring)                 | -                                                                   |
| PutImpuesto                         | OnlyAFIP | [impuesto](#impuesto) [JSON](#json)                            | -                                                                   |
| PutImpuestoList                     | OnlyAFIP | [impuestolist](#impuestolist) [JSON](#json)                    | cantidad [DecimalString](#decimalstring)                            |
| GetImpuestoRange                    | -        | codigo1, codigo2 [DecimalString](#decimalstring)               | [Impuesto](#impuesto) [JSON](#json) array                           |
| DelImpuestoRange                    | OnlyAFIP | codigo1, codigo2 [DecimalString](#decimalstring)               | ids [JSON](#json) number array                                      |
| GetImpuestoAll                      | -        | -                                                              | [Impuesto](#impuesto) [JSON](#json) array                           |
| GetOrganizacion                     | -        | id [DecimalString](#decimalstring)                             | [Organización](#organizacion) [JSON](#json)                         |
| GetOrganizacionAll                  | -        | -                                                              | [Organización](#organizacion) [JSON](#json) array                   |
| PutStates                           | OnlyAFIP | key [String](#string), value [Bytes](#bytes), ...              | cantidad [DecimalString](#decimalstring)                            |
| GetStates                           | OnlyAFIP | key [String](#string) ǁ [Ranges](#ranges) [JSON](#json)        | [GetStatesResponse](#getstatesresponse) [JSON](#json)               |
| DelStates                           | OnlyAFIP | key [String](#string) ǁ [Ranges](#ranges) [JSON](#json)        | keys [JSON](#json) string array                                     |
| GetStatesHistory                    | OnlyAFIP | key [String](#string) ǁ [Ranges](#ranges) [JSON](#json)        | [GetStatesHistoryResponse](#getstateshistoryresponse) [JSON](#json) |

## Tipos básicos

### Bytes

Es un secuencia de bytes arbitrarios.

### String

Es una cadena de caracteres [UNICODE](https://en.wikipedia.org/wiki/Unicode) codificada en [UTF-8](https://en.wikipedia.org/wiki/UTF-8).

### DecimalString

Es un número entero decimal representado como una cadena de caracteres decimales [ASCII](https://en.wikipedia.org/wiki/ASCII)/[UTF-8](https://en.wikipedia.org/wiki/UTF-8).

### JSON

Es un [documento JSON válido](https://json.org) representado como una cadena de caracteres [UNICODE](https://en.wikipedia.org/wiki/Unicode) codificada en [UTF-8](https://en.wikipedia.org/wiki/UTF-8).

## Estructuras

### Organismo

[JSON](#json) que cumple el [schema de organización](https://github.com/padfed/padfed-chaincode/blob/master/resources/schemas/organismo.yaml).

### Impuesto

[JSON](#json) que cumple el [schema de impuesto](https://github.com/padfed/padfed-chaincode/blob/master/resources/schemas/impuesto.yaml).

### Persona

[JSON](#json) que cumple el [schema de persona](https://github.com/padfed/padfed-chaincode/blob/master/resources/schemas/persona.yaml).

#### Colecciones de una persona {id=cols}

Son las colecciones contenidas en [Persona](#persona) [JSON](#json).

Actualmente éstas son:

- Actividades
- Impuestos
- Domicilios
- DomiciliosRoles
- Telefonos
- Jurisdicciones
- Emails
- Archivos
- Categorias
- Etiquetas
- Contribuciones
- Relaciones
- CMSedes

Cada una cumplirá con el sub-schema correspondiente del [schema de persona](https://github.com/padfed/padfed-chaincode/blob/master/resources/schemas/persona.yaml).

### PersonaList

[JSON](#json) que cumple el [schema de lista de personas](https://github.com/padfed/padfed-chaincode/blob/master/resources/schemas/persona-list.yaml).

### PersonaBasica

[JSON](#json) que cumple el sub-schema `#/definitions/persona` del [schema de persona](https://github.com/padfed/padfed-chaincode/blob/master/resources/schemas/persona.yaml).

### GetStatesResponse

Dependiendo del tipo de argumento usado en la invocación a [GetStates](#getstates) podrá distintas estas estructuras según esta tabla:

| Argumento                       | Estructura de Respuesta                           |
| ------------------------------- | ------------------------------------------------- |
| [String](#string)               | raw bytes                                         |
| [Ranges](#ranges) [JSON](#json) | [State](#state) [RangesResponse](#rangesresponse) |

### GetStatesHistoryResponse

Dependiendo del tipo de argumento usado en la invocación a [GetStatesHistory](#getstateshistory) podrá distintas estas estructuras según esta tabla:

| Argumento                       | Estructura de Respuesta                                         |
| ------------------------------- | --------------------------------------------------------------- |
| [String](#string)               | [StateHistory](#statehistory)                                   |
| [Ranges](#ranges) [JSON](#json) | [StateHistory](#statehistory) [RangesResponse](#rangesresponse) |

### State

Es un objeto [JSON](https://json.org) con la siguiente estructura:

| Atributo   | Tipo JSON | Semántica                                                                                                                                                                                                                            |
| ---------- | --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `key`      | string    | Key del state.                                                                                                                                                                                                                       |
| `content`  | string    | Si el contenido binario del state es una [String](#string) [UTF-8](https://en.wikipedia.org/wiki/UTF-8) válida tendrá su valor literal, caso contrario tendrá el valor codificado en [Base64](https://en.wikipedia.org/wiki/Base64). |
| `encoding` | string    | Si fue necesario codificar el valor tendrá el nombre del algoritmo utilizado (`base64`), caso contrario no estará presente.                                                                                                          |

### StateHistory

Es un objeto [JSON](https://json.org) con la siguiente estructura:

| Atributo  | Tipo JSON                         | Semántica                   |
| --------- | --------------------------------- | --------------------------- |
| `key`     | string                            | Key del state.              |
| `history` | [StateChange](#statechange) array | Lista de cambios del state. |

### StateChange

Es un objeto [JSON](https://json.org) con la siguiente estructura:

| Atributo   | Tipo JSON | Semántica                                                                                                                                                                                                                                                        |
| ---------- | --------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `txid`     | string    | Identificador de la transacción en la que se produjo el cambio.                                                                                                                                                                                                  |
| `time`     | string    | Timestamp de la modificación en formato [RFC3339](https://www.rfc-editor.org/rfc/rfc3339.html).                                                                                                                                                                  |
| `delete`   | bool      | Indica si la modificación fue la eliminación del state.                                                                                                                                                                                                          |
| `content`  | string    | Contenido del state luego de la modificación. Si el contenido es una [String](#string) [UTF-8](https://en.wikipedia.org/wiki/UTF-8) válida tendrá su valor literal, caso contrario tendrá el valor codificado en [Base64](https://en.wikipedia.org/wiki/Base64). |
| `encoding` | string    | Si fue necesario codificar el contenido tendrá el nombre del algoritmo utilizado (`base64`), caso contrario no estará presente.                                                                                                                                  |

### Ranges

Es un array [JSON](https://json.org) de longitud arbitraria cuyos **elementos** pueden ser cualquiera de las siguientes posibilidades:

| Tipo Elemento | Tipo JSON       | Ejemplo         | Semántica                                                                                                                                                                                                                                                             |
| ------------- | --------------- | --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `key`         | string          | `"foo"`         | Especifica una key.                                                                                                                                                                                                                                                   |
| `range`       | string array[2] | `["foo","bar"]` | Especifica un rango de keys. Coincidirá toda key que sea **igual o mayor** que la primera y **menor** que la segunda. Tanto la primera como la segunda pueden estar vacías (`""`) indicando que se desea obtener desde el principio o hasta el final respectivamente. |
| `prefix`      | string array[1] | `["foo"]`       | Especifica un rango de keys **por prefijo**. Coincidirá toda key que comience con la string suministrada. Esta puede estar vacía (`""`) indicando un prefijo vacío (coincide con todas las keys).                                                                     |

#### Ejemplos

Un elemento que coincide con la key `foo`:

```json
["foo"]
```

Un elemento que coincide con las keys desde `bar` (inclusive) hasta `foo` (no inclusive):

```json
[["bar","foo"]]
```

Dos elementos donde el primero coincide con las keys desde `bar` (inclusive) hasta `foo` (no inclusive) y el segundo coincide únicamente con la key `qux`:

```json
[["bar","foo"],"qux"]
```

Todas las keys desde `bar` (inclusive) hasta `foo` (no inclusive) y todas las keys con prefijo `qux`:

```json
[["bar","foo"],["qux"]]
```

Todas las keys desde `bar` (inclusive) hasta `foo` (no inclusive), la key específica `qux` y todas las keys con prefijo `wal`:

```json
[["bar","foo"],"qux",["wal"]]
```

Todas las keys:

```json
[[""]]
```

o bien

```json
[["",""]]
```

Todas las keys desde `qux` (inclusive) en adelante:

```json
[["qux",""]]
```

Todas las keys desde la primera que existe hasta `qux` (no inclusive):

```json
[["","qux"]]
```

### RangesResponse

Es un arreglo [JSON](#json) de longitud variable cuyos elementos tendrán una estructura dada según el elemento correspondiente de la estructure [Ranges](#ranges) a la que responde.

| Tipo elemento en array solicitud [Ranges](#ranges) | Tipo elemento en array resultado |
| -------------------------------------------------- | -------------------------------- |
| `key`                                              | [State](#state)                  |
| `range`                                            | [State](#state) array            |
| `prefix`                                           | [State](#state) array            |

Esto quiere decir que si el elemento en la posición `N` de la solicitud es un...

- `key`: entonces el elemento en la posición `N` de la respuesta será un [State](#state)
- `range`: entonces el elemento en la posición `N` de la respuesta será un **arreglo de [State](#state)**
- `prefix`: entonces el elemento en la posición `N` de la respuesta será un **arreglo de [State](#state)**

#### Ejemplos

Solicitud:

```json
[ "foo" ]
```

Respuesta posible:

```json
[
    {"key":"foo","content":"x"}
]
```

Solicitud:

```json
[ ["bar","foo"] ]
```

Respuesta posible:

```json
[
    [ {"key":"bar","content":"x"}, {"key":"bar1","content":"x"}, {"key":"car","content":"x"} ]
]
```

Solicitud:

```json
[ ["bar","foo"], "qux" ]
```

Respuesta posible:

```json
[
    [ {"key":"bar","content":"x"}, {"key":"bar1","content":"x"}, {"key":"car","content":"x"} ],
      {"key":"qux","content":"x"}
]
```

Solicitud:

```json
[ ["bar","foo"], ["qux"] ]
```

Respuesta posible:

```json
[
    [ {"key":"bar","content":"x"}, {"key":"car1","content":"x"}, {"key":"foa32","content":"x"} ],
    [ {"key":"qux","content":"x"}, {"key":"quxb1","content":"x"} ]
]
```

Solicitud:

```json
[ ["bar","foo"], "wal", ["qux"] ]
```

Respuesta posible:

```json
[
    [ {"key":"bar","content":"x"}, {"key":"car1","content":"x"}, {"key":"foa32","content":"x"} ],
      {"key":"wal","content":"x"},
    [ {"key":"qux","content":"x"}, {"key":"quxb1","content":"x"} ]
]
```
