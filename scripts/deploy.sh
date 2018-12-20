echo $*
echo Copiando fuentes del Chaincode a directorio de destino
echo Se utilizara 'GO_HOME' o bien si la misma no se encuentra definida el directorio default HOME/go

if [ "$4" = "" ]
then
    echo Variable GO no recibida se asumira $HOME/go
    GO_DEST=$1/go
  else
    echo Variable GO detectada en el ambiente $GO_HOME
    GO_DEST=$GO_HOME
fi

mkdir -p $GO_DEST/src/afip/tribfed/chaincodes/padfed/$2

cp -u -r $3/src/main/afip/tribfed/chaincodes/padfed $GO_DEST/src/afip/tribfed/chaincodes/padfed/$2