Para desplegar el chaincode con las dependecias adecuadas, se recomienda mover el fuente al GOPATH definido localmente
Ejemplo:
/$HOME/go/src/afip/tribfed/chaincodes/padfed

*ejecutar luego

   govendor init


* Agregar la dependencia externa al directorio actual

    govendor fetch github.com/lalloni/afip/cuit
    govendor fetch github.com/hyperledger/fabric/core/chaincode/shim
    govendor fetch github.com/hyperledger/fabric/core/peer
    govendor fetch github.com/spacemonkeygo/errors
    govendor fetch github.com/xeipuuv/gojsonschema