En este proyecto se mantiene el codigo fuente del chaincode **padfed_cc** escrito en GOLANG

# Codigo fuente del chaincode

En esta seccion exlicar temas referidos al codigo funete

# Gestion del proyecto con maven

El proyecto contiene un pom.xml para ser gestionado mediante maven

En esta primera versión del pom.xml incluye la siguiente funcionalidad:

    - Replicación de código fuente a directorio $GO_HOME o $HOME/GO segun entorno configurado

    - Empaquetado de fuentes y vendor deps para publicar en Nexus

Ejemplo:

    mvn package --> Genera empaquetado (aun no utilizado) y copia fuentes .go y vendor deps en $HOME/go o $GO_HOME

    En el directorio de destino la tool generá un directorio versionado según la version pom.xml definida para mantener distintas versiones de fuentes

