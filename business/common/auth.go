package common

import (
	"regexp"
	"strconv"

	"github.com/lalloni/fabrikit/chaincode/authorization"
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/pkg/errors"
)

var (
	AFIP = authorization.MSPID("AFIP")
	Free = authorization.Allowed
)

var cuitSerialNumberRegexp = regexp.MustCompile(`^CUIT (\d{11})$`)

func ExtractCUITFromSerialNumber(ctx *context.Context) (uint64, error) {
	cert, err := ctx.ClientCertificate()
	if err != nil {
		return 0, errors.Wrap(err, "getting client certificate")
	}
	match := cuitSerialNumberRegexp.FindStringSubmatch(cert.Subject.SerialNumber)
	if match == nil {
		return 0, errors.New("invalid cuit syntax in certificate subject serialnumber")
	}
	cuit, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "invalid cuit in certificate subject serialnumber")
	}
	return cuit, nil
}
