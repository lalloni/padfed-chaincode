# API

## Generalidades

Todas las funciones (no deprecadas) del chaincode responden utilizando la estructura de respuesta de Fabric de la siguiente manera:

- Status: se utiliza un status code con [valores y semántica tomados de HTTP](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)
- Message: se utiliza para informar mensajes de error (si lo hubiera) mostrables a usuarios finales
- Payload: se utiliza para contener una estructura de respuesta estandarizada para todas las funciones (ver [Payload](#payload))

### Payload

El Payload siempre tenrá la siguiente forma general:

```json
{
    "client": {},
    "transaction": {},
    "content": {},
    "content-encoding": "",
    "fault": {}
}
```

En el cual:

- `client`: es un objeto que contiene información relativa al sistema cliente que realizó la invocación (ver [Client](#client))
- `transaction`: es un objeto que contiene información relativa a la transacción en la que se ejecutó la invocación (ver [Transaction](#transaction))
- `content`: es un objeto que contiene la respuesta a la operación
- `content-encoding`: es una cadena de caracteres que contiene el identificador del algoritmo de codificación que se utilizó para representar en `content` una secuencia arbitraria de bytes, si hubiera sido necesario (esto ocurre cuando los bytes no pueden representarse como una cadena de UTF-8 válida)
- `fault`: es un objeto o cadena de caracteres que contiene información acerca de la causa que impidió el éxito de la invocación, si este hubiera fallado

Condiciones:

- `fault` y `content` son mutuamente excluyentes
- `content-encoding` es requerido si se codificó `content` (actualmente solo puede tener el valor "base64" si el contenido de `content` no era expresable como una string UTF-8 y fue necesaria la codificación, caso contrario no estára presente indicando que no hubo codificación y que `content` contiene una cadena de caracteres UTF-8 literal)
- `client` y `transaction` solo estarán presentes si se realizó la invocación incluyendo la opción `debug`, caso contrario no estarán presentes

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

| Nombre           | ACL      | Argumentos                                               | Payload                           | Comportamiento |
| ---------------- | -------- | -------------------------------------------------------- | --------------------------------- | -------------- |
| GetVersion       | -        | -                                                        |                                   |                |
| GetFunctions     | -        | -                                                        |                                   |                |
| GetPersona       | -        | cuit [DecimalString](#decimalstring)                     | [Persona](#persona) [JSON](#json) |                |
| DelPersona       | OnlyAFIP | cuit [DecimalString](#decimalstring)                     |                                   |                |
| PutPersona       | OnlyAFIP | [persona](#persona) [JSON](#json)                        |                                   |                |
| PutPersonaList   | OnlyAFIP | [personalist](#personalist) [JSON](#json)                |                                   |                |
| GetPersonaRange  | OnlyAFIP | cuit1, cuit2 [DecimalString](#decimalstring)             |                                   |                |
| DelPersonaRange  | OnlyAFIP | cuit1, cuit2 [DecimalString](#decimalstring)             |                                   |                |
| GetPersonaAll    | OnlyAFIP | -                                                        |                                   |                |
| PutStates        | OnlyAFIP | key [String](#string), value [Bytes](#bytes) ...         |                                   |                |
| GetStates        | OnlyAFIP | key [String](#string) \| [ranges](#ranges) [JSON](#json) |                                   |                |
| DelStates        | OnlyAFIP | key [String](#string) \| [ranges](#ranges) [JSON](#json) |                                   |                |
| GetStatesHistory | OnlyAFIP | key [String](#string) \| [ranges](#ranges) [JSON](#json) |                                   |                |

## Tipos

### String

Es una cadena de caracteres codificada en [UTF-8](https://en.wikipedia.org/wiki/UTF-8).

### DecimalString

Es una número entero decimal representado como una cadena de caracteres decimales.

### JSON

Es un [documento JSON válido](https://json.org) representado como una cadena de caracteres codificada en [UTF-8](https://en.wikipedia.org/wiki/UTF-8).

## Estructuras

### Persona

[JSON](#json) que cumple el [schema de persona](https://github.com/padfed/padfed-validator/blob/master/doc/schemas/persona.json).

### PersonaList

[JSON](#json) que cumple el [schema de lista de personas](https://github.com/padfed/padfed-validator/blob/master/doc/schemas/persona-list.json).

### Ranges

Es un array [JSON](https://json.org) de longitud arbitraria cuyos elementos pueden ser cualquiera de las siguientes posibilidades:

| Tipo   | JSON   | Ejemplo         | Semántica                                                                                                                                                                                                                                                  |
| ------ | ------ | --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| key    | string | `"foo"`         | Especifica una key.                                                                                                                                                                                                                                        |
| range  | array  | `["foo","bar"]` | Especifica un rango de keys. Coincidirá toda key que sea **igual o mayor** que la primera y **menor** que la segunda. Tanto la primera como la segunda pueden ser `""` indicando que se desea obtener desde el principio o hasta el final respectivamente. |
| prefix | array  | `["foo"]`       | Especifica un rango de keys **por prefijo**. Coincidirá toda key que comience con la string. La string puede ser `""` indicando un prefijo vacío (todas las keys).                                                                                         |

#### Ejemplos

Sólo la key `foo`:

```json
["foo"]
```

Desde `foo` hasta `bar` (no inclusive):

```json
[["foo","bar"]]
```

Desde `foo` hasta `bar` (no inclusive) y la key específica `baz`:

```json
[["foo","bar"],"baz"]
```

Desde `foo` hasta `bar` (no inclusive) y las keys con prefijo `baz`:

```json
[["foo","bar"],["baz"]]
```

Desde `foo` hasta `bar` (no inclusive), la key específica `qux` y las keys con prefijo `baz`:

```json
[["foo","bar"],"qux",["baz"]]
```

Todas las keys:

```json
[[""]]
```

o bien

```json
[["",""]]
```

Todas las keys desde `qux` en adelante:

```json
[["qux",""]]
```

Todas las keys desde la primera que existe hasta `qux`:

```json
[["","qux"]]
```
