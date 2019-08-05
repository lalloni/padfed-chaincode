package organizaciones

import (
	"io/ioutil"

	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/resources"
)

var (
	orgArray   = []*Org{}
	orgByID    = map[uint64]*Org{}
	orgByCUIT  = map[uint64]*Org{}
	orgByMSPID = map[string]*Org{}
)

func init() {
	orgs := mustLoad()
	for _, org := range orgs {
		org := org // alocar nueva variable y copiarle el valor
		orgArray = append(orgArray, &org)
		orgByID[org.ID] = &org
		orgByCUIT[org.CUIT] = &org
		orgByMSPID[org.MSPID] = &org
	}
	AFIP = orgByID[1]
	MORGS = orgByID[100]
}

func mustLoad() []Org {
	orgs, err := load()
	if err != nil {
		panic(err)
	}
	return orgs
}

func load() ([]Org, error) {
	f, err := resources.Data.Open("organizaciones.csv")
	if err != nil {
		return nil, errors.Wrap(err, "opening orgs")
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "reading orgs")
	}
	orgs := []Org(nil)
	err = csvutil.Unmarshal(bs, &orgs)
	if err != nil {
		return nil, errors.Wrap(err, "parsing orgs")
	}
	return orgs, nil
}
