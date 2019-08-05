package organizaciones

var (
	AFIP  *Org
	MORGS *Org
)

func IsMORGS(o *Org) bool {
	return o == MORGS
}

func IsAFIP(o *Org) bool {
	return o == AFIP
}

func GetAll() []*Org {
	return orgArray
}

func GetByID(id uint64) *Org {
	return orgByID[id]
}

func GetByCUIT(cuit uint64) *Org {
	return orgByCUIT[cuit]
}

func GetByMSPID(msp string) *Org {
	return orgByMSPID[msp]
}
