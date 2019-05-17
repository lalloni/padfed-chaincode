package generic

import (
	"encoding/json"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

// PutStatesHandler es un handler variádico que recibe como argumentos una
// secuencia de pares de key (raw string) y contenido (raw bytes) a guardar
// en states con el key literal recibido.
//
// Retorna la cantidad de states almacenados.
func PutStatesHandler(ctx *context.Context) *response.Response {
	args := ctx.Stub.GetArgs()
	// NOTE: args[0] is the function so len(args) should be odd
	if len(args)%2 == 0 {
		return response.BadRequest("arguments must be sequence of one or more key (string) and value (bytes) pairs")
	}
	count := 0
	for i := 1; i < len(args); i += 2 {
		key, value, res := kvarg(ctx, i)
		if res != nil {
			return res
		}
		res = kvput(ctx, key, value)
		if res != nil {
			return res
		}
		count++
	}
	return response.OK(count)
}

// GetStatesHandler es un handler que puede recibir como argumento un raw string
// (no JSON) o bien un JSON array de strings.
//
// En caso de recibir un raw string se retorna un array de bytes con el contenido
// recuperado de la key. En caso de no existir la key retorna un status code not
// found.
//
// En caso de recibir un JSON array de strings retorna un JSON array de objetos que
// contienen los atributos "key", "content" y opcionalmente "encoding". Si la key
// no existe no incluye el atributo "content" en la respuesta.
//
// Si "content" no pudiera ser representado como una string codificada en UTF-8,
// será codificado en Base64 y se asignará el valor "base64" al atributo
// "encoding", caso contrario no se incluye.
func GetStatesHandler(ctx *context.Context) *response.Response {
	arg, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("getting argument: %v", err)
	}
	keys := []string{}
	err = json.Unmarshal(arg, &keys)
	if err != nil {
		// interpretar arg como una raw string
		key := string(arg)
		bs, res := kvget(ctx, key)
		if res != nil {
			return res
		}
		if bs == nil {
			return response.NotFound()
		}
		return response.OK(bs)
	}
	result := []*state{}
	for _, key := range keys {
		bs, res := kvget(ctx, key)
		if res != nil {
			return res
		}
		result = append(result, newstate(key, bs))
	}
	return response.OK(result)
}

// DelStatesHandler es un handler que puede recibir como argumento un raw string
// (no JSON) o bien un JSON array de strings y las marca para eliminación.
func DelStatesHandler(ctx *context.Context) *response.Response {
	arg, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("getting argument: %v", err)
	}
	keys := []string{}
	err = json.Unmarshal(arg, &keys)
	if err != nil {
		// interpretar arg como una raw string
		keys = []string{string(arg)}
	}
	for _, key := range keys {
		err := ctx.Stub.DelState(key)
		if err != nil {
			return response.Error("deleting key: %v", err)
		}
	}
	return response.OK(nil)
}

func GetStatesHistoryHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented()
}

func GetStatesRangeHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented()
}

func DelStatesRangeHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented()
}

func kvarg(ctx *context.Context, pos int) (string, []byte, *response.Response) {
	key, err := ctx.ArgString(pos)
	if err != nil {
		return "", nil, response.BadRequest("getting key argument: %v", err)
	}
	value, err := ctx.ArgBytes(pos + 1)
	if err != nil {
		return "", nil, response.BadRequest("getting value argument: %v", err)
	}
	return key, value, nil
}

func kvput(ctx *context.Context, key string, value []byte) *response.Response {
	err := ctx.Stub.PutState(key, value)
	if err != nil {
		return response.Error("putting state: %v", err)
	}
	return nil
}

func kvget(ctx *context.Context, key string) ([]byte, *response.Response) {
	bs, err := ctx.Stub.GetState(key)
	if err != nil {
		return nil, response.Error("getting key state: %v", err)
	}
	return bs, nil
}
