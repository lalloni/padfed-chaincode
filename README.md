# padfed-chaincode

- **Roadmap** https://redmine-blockchain.afip.gob.ar/projects/padfed-cc/roadmap
- **Redmine** https://redmine-blockchain.afip.gob.ar/projects/padfed-cc
- **GitLab** https://gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode
- **Nexus** https://nexus.cloudint.afip.gob.ar/nexus/service/rest/repository/browse/padfed-bc-maven/afip/padfed/padfedcc/

## Objetivo

Este proyecto mantiene el código fuente del chaincode (SmartContract) **padfed-chaincode** escrito en Go.

## Desarrollo

### Entorno

#### Prerrequisitos

- Go 1.11 o posterior (ver sección siguiente para instalar o actualizar)
- Git

#### Go

Si no se cuenta con una versión adecuada de Go se puede instalar o actualizar ejecutando:

```sh
sudo sh -c 'curl -sSL# https://dl.google.com/go/go1.12.linux-amd64.tar.gz -o /tmp/g.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xaf /tmp/g.tar.gz'
```

#### Task

Instalar o actualizar ejecutando:

```sh
sh -c 'curl -sL https://taskfile.dev/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0'
```

### Tareas

El Taskfile incluye definiciones de tareas para:

#### Ejecutar tests

```sh
task test
```

#### Informar tests code coverage

Genera un reporte de cobertura de tests y abre un web browser para visualizarlo.

```sh
task cover
```

#### Construir el binario final

```sh
task compile
```

#### Empaquetar fuentes

```sh
task package
```

El paquete generado se versionará según último tag en Git.

#### Publicar versión snapshot

Publicará una versión de trabajo del paquete de fuentes del CC.

Realizará:

- Ejecución de tests y análisis estático de código
- Empaquetamiento
- Upload a Nexus

El paquete publicado se versionará según último tag en Git.

```sh
task publish
```

#### Generar y publicar versión release de paquete de fuentes del CC

Realizará:

- Ejecución de tests y análisis estático de código
- Tag en Git
- Empaquetamiento
- Upload a Nexus

```sh
task release version=1.2.3
```
