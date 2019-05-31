# CHANGES

## 0.6.x

- [x] Reingeniería de chaincode
- [x] Rediseño de API de chaincode
- [x] Migrar a store QueryPersona/All, DelPersonaByRange
- [x] Representación opcional LZ4 en State
- [x] Emprolijar API genérica
- [x] ACL
- [x] Rediseño de API externa del chaincode (status, uniformar uso de campos payload y message, uniformar estructura de respuesta)
- [ ] Migrar código genérico a librerías fabric-cc-kit
  - [ ] Paquete `ng`
  - [ ] Paquere `store`
- [ ] Migrar código de negocio de padfed-validator a padfed-chaincode
  - [ ] Esquemas
- [ ] Emprolijar "modo desarrollo" (funciones de cc para testing)
- [ ] Impuesto en Store como composites

## 0.5.x

- [x] Completar modelo de datos
  - [x] Roles de domicilios
  - [x] Actividadae Jurisdiccionales

**TODO** completar elementos faltantes

## 0.4.x

- [x] Introducción de Store
- [x] Validar opcionalidad de singletons y collections
- [x] Migrar GetPersona a Store
- [x] Migrar DelPersona a Store
- [x] Migrar PutPersona a Store
- [x] Migrar PutPersonas a Store
- [x] Eliminar consultas y transacciones que operan sobre PersonaImpuesto
- [x] Actualizar formato de keys en queryPersona/All, delPersonaByRange
