# padfed-chaincode

- **Roadmap** https://redmine-blockchain.afip.gob.ar/projects/padfed-cc/roadmap
- **Redmine** https://redmine-blockchain.afip.gob.ar/projects/padfed-cc
- **GitLab** https://gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode
- **Nexus** https://nexus.cloudint.afip.gob.ar/nexus/service/rest/repository/browse/padfed-bc-maven/afip/padfed/padfedcc/

## Objetivo

Este proyecto mantiene el codigo fuente del chaincode (SmartContract) **padfedcc** escrito en GOLANG

## Maven

El pom.xml incluye las siguientes funcionalidades:

- Replicación de código fuente a directorio $GO_HOME o $HOME/GO segun entorno configurado

- Empaquetado de fuentes y vendor deps en un jar para ser publicado en Nexus

Ejecución | Acción
--- | ---
mvn package | Genera un jar y copia los fuentes .go y vendor en $HOME/go o $GO_HOME. En el directorio de destino generá un directorio versionado según la version indicada en el pom.xml.

