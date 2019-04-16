# ROADMAP

## 0.4.0

### Hitos

- Introducción de Store

### Cambios específicos

- Validar opcionalidad de singletons y collections
- Migrar GetPersona a Store
- Migrar DelPersona a Store
- Migrar PutPersona a Store
- Migrar PutPersonas a Store
- Eliminar consultas y transacciones que operan sobre PersonaImpuesto
- Actualizar formato de keys en queryPersona/All, delPersonaByRange

### Tareas

- Deploy en testnet
- Enviar novedades

## 0.5.0

### Hitos

- Emprolijar modo desarrollo
- Impuesto en Store como composites

### Cambios específicos

- Nueva transacción PutEstadoPersonaImpuesto (modifica PersonaImpuesto)
- Nuevas consultas GetPersonaSimple y GetPersonaImpuestoList
- Migrar a store QueryPersona/All, DelPersonaByRange
- Nueva transacción PutImpuesto
- Nueva consulta GetImpuesto

## Previo a 1.0.0

### Hitos

- ACL
- Cache de parámetros persistidos en State
- Emprolijar API genérica
- Metadata de representación binaria en State
- Representación LZ4 en State