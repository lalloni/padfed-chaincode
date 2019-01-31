Este proyecto mantiene el codigo fuente del chaincode (SmartContract) **padfedcc** escrito en GOLANG

- **Redmine** https://redmine-blockchain.afip.gob.ar/projects/padfed-cc

# Maven

El pom.xml incluye las siguientes funcionalidades:

- Replicación de código fuente a directorio $GO_HOME o $HOME/GO segun entorno configurado

- Empaquetado de fuentes y vendor deps en un jar para ser publicado en Nexus

Ejecucion | Acción
--- | ---
mvn package | Genera un jar y copia los fuentes .go y vendor en $HOME/go o $GO_HOME. En el directorio de destino generá un directorio versionado según la version indicada en el pom.xml.

