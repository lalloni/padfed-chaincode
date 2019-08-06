package personas

import (
	"strconv"

	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handlerutil/crud"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/common"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/impuestos"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/business/organizaciones"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/state"
)

func AddHandlers(r router.Router) {
	addHandlers(r, false)
}

func addTestingHandlers(r router.Router) {
	addHandlers(r, true)
}

func addHandlers(r router.Router, testing bool) {
	opts := append(
		crud.Defaults, // i.e. get, getrange, has, put, putlist, del, delrange
		crud.WithIDParam(CUITParam),
		crud.WithItemParam(PersonaParam),
		crud.WithListParam(PersonaListParam),
		crud.WithPutValidator(validatePersona),
	)
	if !testing {
		opts = append(opts, crud.WithWriteCheck(common.AFIP))
	}
	crud.AddHandlers(r, Schema, opts...)

	r.SetHandler("QueryPersona", common.Free, QueryPersonaHandler)
	r.SetHandler("QueryPersonaBasica", common.Free, QueryPersonaBasicaHandler)

	if testing {
		r.SetHandler("SetPersonaImpuestoEstado", common.Free, SetPersonaImpuestoEstadoHandler(testingMSPID, testingCUIT))
	} else {
		r.SetHandler("SetPersonaImpuestoEstado", common.Free, SetPersonaImpuestoEstadoHandler(contextMSPID, contextCUIT))
	}
}

func validatePersona(ctx *context.Context, v interface{}) *response.Response {

	per := v.(*Persona)

	if per.ID == 0 {
		return response.BadRequest("id required")
	}

	if per.Persona == nil {
		exist, err := ctx.Store.HasComposite(Schema, per.ID)
		if err != nil {
			return response.Error("checking persona existence: %v", err)
		}
		if !exist {
			return response.BadRequest("persona is required when putting a new instance")
		}
	} else if per.Persona.ID != per.ID {
		return response.BadRequest("root id %d must equal persona id %d", per.ID, per.Persona.ID)
	}

	return nil
}

func QueryPersonaHandler(ctx *context.Context) *response.Response {
	var cuit uint64
	_, err := handler.ExtractArgs(ctx.Args(), CUITParamVar(&cuit))
	if err != nil {
		return response.BadRequest("invalid arguments: %v", err)
	}

	prefix := "per:" + strconv.FormatUint(cuit, 10) + "#"
	query := state.Single(state.Range(state.PrefixRange(prefix)))

	r, err := state.QueryKeyRanges(ctx, query)
	if err != nil {
		return response.Error(err.Error())
	}
	return response.OK(r)
}

func QueryPersonaBasicaHandler(ctx *context.Context) *response.Response {
	var cuit uint64
	_, err := handler.ExtractArgs(ctx.Args(), CUITParamVar(&cuit))
	if err != nil {
		return response.BadRequest("invalid arguments: %v", err)
	}
	key := "per:" + strconv.FormatUint(cuit, 10) + "#per"
	query := state.Single(state.Point(key))
	r, err := state.QueryKeyRanges(ctx, query)
	if err != nil {
		return response.Error(err.Error())
	}
	return response.OK(r)
}

type MSPIDProvider func(ctx *context.Context) (string, error)

type CUITProvider func(ctx *context.Context) (uint64, error)

func SetPersonaImpuestoEstadoHandler(clientMSPID MSPIDProvider, clientCUIT CUITProvider) handler.Handler {
	return func(ctx *context.Context) *response.Response {
		var (
			cuit   uint64
			codimp uint64
			estado string
		)
		_, err := handler.ExtractArgs(ctx.Args(), CUITParamVar(&cuit), impuestos.CodigoImpuestoParamVar(&codimp), EstadoParamVar(&estado))
		if err != nil {
			return response.BadRequest("invalid arguments: %v", err)
		}
		if !allowSettingEstado[estado] {
			return response.BadRequest("setting estado to '%s' is not allowed", estado)
		}
		v, err := ctx.Store.GetComposite(impuestos.Schema, codimp)
		if err != nil {
			return response.Error("checking impuesto existence: %v", err)
		}
		if v == nil {
			return response.BadRequest("impuesto %v not found", codimp)
		}
		itemid := strconv.FormatUint(codimp, 10)
		item, err := ctx.Store.GetCompositeCollectionItem(Impuestos, cuit, itemid)
		if err != nil {
			return response.Error("checking inscripción existence: %v", err)
		}
		if item == nil {
			return response.NotFoundWithMessage("inscripción persona %v impuesto %v not found", cuit, codimp)
		}
		inscripción := item.(*Impuesto)
		if inscripción.Estado == estado {
			return response.BadRequest("setting estado to same value is not allowed")
		}
		msp, err := clientMSPID(ctx)
		if err != nil {
			return response.Error("getting client MSPID: %v", err)
		}
		org := organizaciones.GetByMSPID(msp)
		if organizaciones.IsMORGS(org) {
			callercuit, err := clientCUIT(ctx)
			if err != nil {
				return response.BadRequest("getting client CUIT: %v", err)
			}
			org = organizaciones.GetByCUIT(callercuit)
			if org == nil {
				return response.BadRequest("organización with CUIT %v not found", callercuit)
			}
		}
		impuesto := v.(*impuestos.Impuesto)
		if !organizaciones.IsAFIP(org) && org.ID != impuesto.Org {
			return response.Forbidden("organización %v (%s) can not set estado inscripción impuesto %v (%s)", org.ID, org.Nombre, impuesto.Codigo, organizaciones.GetByID(impuesto.Org).Nombre)
		}
		// OK! Actualizar estado y persistir
		inscripción.Estado = estado
		err = ctx.Store.PutCompositeCollectionItem(Impuestos, cuit, itemid, inscripción)
		if err != nil {
			return response.Error("setting estado inscripción: %v", err)
		}
		return response.OK(nil)
	}
}

var allowSettingEstado map[string]bool

func init() {
	ss := []string{"NA", "BD", "EX"}
	allowSettingEstado = make(map[string]bool, len(ss))
	for _, s := range ss {
		allowSettingEstado[s] = true
	}
}

func contextMSPID(ctx *context.Context) (string, error) {
	return ctx.ClientMSPID()
}

func contextCUIT(ctx *context.Context) (uint64, error) {
	return common.ExtractCUITFromSerialNumber(ctx)
}

var (
	testClientMSPID        = "test"
	testClientCUIT  uint64 = 1
)

func testingMSPID(ctx *context.Context) (string, error) {
	return testClientMSPID, nil
}

func testingCUIT(ctx *context.Context) (uint64, error) {
	return testClientCUIT, nil
}
