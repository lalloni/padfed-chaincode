package generic

import (
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
		key, value, err := ctx.ArgKV(i)
		if err != nil {
			return response.BadRequest("getting key-value argument at %d: %v", i, err)
		}
		res := kvput(ctx, key, value)
		if res != nil {
			return res
		}
		count++
	}
	return response.OK(count)
}

// GetStatesHandler es un handler que puede recibir como argumento un raw string
// (no JSON) o bien un JSON array.
//
// En caso de recibir un raw string se retorna un array de bytes con el contenido
// recuperado de la key. En caso de no existir la key retorna un status code not
// found.
//
// En caso de recibir un JSON array, sus elementos pueden ser strings que
// representan keys cuyos values se desea obtener, o bien arrays que contienen:
//
//     - Un par [string1,string2] que será interpretado como un rango desde-hasta
//       de keys.
//       Si la primer string es "" (string vacía) significa leer desde el principio.
//       Si la segunda string es "" (string vacía) significa leer hasta el final.
//     - Una única string [string] que será interpretada como el prefijo del rango
//       de keys.
//       Si la string es "" (string vacía) equivale a la lectura completa.
//
// En caso de recibir un JSON array, retorna un JSON array de objetos que
// contienen los atributos "key", "content" y opcionalmente "encoding".
//
// Si en el arreglo se especificó una key puntual inexistente entonces el objeto
// correspondiente no incluirá el atributo "content" en la respuesta.
//
// Si los bytes de "content" no pudieran ser representados como una string
// codificada en UTF-8, será codificado en Base64 y se asignará el valor "base64"
// al atributo "encoding", caso contrario no se incluye.
func GetStatesHandler(ctx *context.Context) *response.Response {
	arg, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("getting argument: %v", err)
	}
	return processRanges(ctx, arg, func(single bool, key string) (interface{}, error) {
		s, res := kvget(ctx, key)
		if res != nil {
			return res, nil
		}
		if single && s.Nil {
			return response.NotFound(), nil
		}
		if single {
			return s.Content, nil
		}
		return s, nil
	})
}

// DelStatesHandler es un handler que puede recibir como argumento un raw string
// (no JSON) o bien un JSON array de strings y las marca para eliminación.
func DelStatesHandler(ctx *context.Context) *response.Response {
	arg, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("getting argument: %v", err)
	}
	return processRanges(ctx, arg, func(single bool, key string) (interface{}, error) {
		err := ctx.Stub.DelState(key)
		if err != nil {
			return response.Error("deleting state: %v", err), nil
		}
		return key, nil
	})
}

// GetStatesHistoryHandler es un handler que puede recibir como argumento un raw
// string (no JSON) o bien un JSON array de strings.
//
// En caso de recibir un raw string se retorna un array de JSON objects que
// describen cada modificación de la key recibida.
//
// En caso de recibir un JSON array de strings retorna un JSON array de objetos que
// contienen los atributos "key" y "history" donde este último es un array de
// JSON objects que describen cada modificación de la key.
//
// En ambos casos los objetos que describen las modificaciones contienen los
// atributos "txid", "time", "delete", "content" Y "encoding".
//
// "txid" es una string con el identificador de la transacción que incluyó la
// modificación.
//
// "time" es una string con fecha y hora local correspondiente a la transacción.
//
// "delete" es un boolean que indica si en la transacción se eliminó la key.
//
// "content" es una string con el contenido que adoptó el state en la modificación.
//
// "encoding" es una string que tendrá el valor "base64" si "content" no se
// puede representar como una string UTF-8 y se codificó en Base64. En caso
// contrario estará ausente.
func GetStatesHistoryHandler(ctx *context.Context) *response.Response {
	arg, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest("getting argument: %v", err)
	}
	return processRanges(ctx, arg, func(single bool, key string) (interface{}, error) {
		mods, res := khget(ctx, key)
		if res != nil {
			return res, nil
		}
		return &statehistory{Key: key, History: mods}, nil
	})
}

// GetAllHandler devuelve todos los states
func GetAllHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented()
}
