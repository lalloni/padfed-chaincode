package persona

import (
	"reflect"
	"testing"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/context"
	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/ng/response"
)

func TestGetPersonaHandler(t *testing.T) {
	type args struct {
		ctx *context.Context
	}
	tests := []struct {
		name string
		args args
		want *response.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPersonaHandler(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPersonaHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
