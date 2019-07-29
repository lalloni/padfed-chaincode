package organizaciones

func GetAll() []*Org {
	return orgArray
}

func GetByID(id uint64) *Org {
	return orgByIDIndex[id]
}

func GetByCUIT(cuit uint64) *Org {
	return orgByCUITIndex[cuit]
}
