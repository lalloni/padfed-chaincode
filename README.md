En este proyecto se mantiene el codigo fuente del chaincode **padfed_cc** escrito en GOLANG

# Redmine

https://redmine-blockchain.afip.gob.ar/projects/padfed-cc

# Codigo fuente del chaincode

En esta seccion exlicar temas referidos al codigo funete

# Gestion del proyecto con maven

En esta primera versión el pom.xml incluye las siguientes funcionalidades:

    - Replicación de código fuente a directorio $GO_HOME o $HOME/GO segun entorno configurado

    - Empaquetado de fuentes y vendor deps para publicar en Nexus

Ejemplo:

    mvn package --> Genera empaquetado (aun no utilizado) y copia fuentes .go y vendor deps en $HOME/go o $GO_HOME

    En el directorio de destino generá un directorio versionado según la version indicada en el pom.xml asignada para mantener distintas versiones del chaincode

