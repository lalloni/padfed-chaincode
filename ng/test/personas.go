package test

import (
	"math"
	"sort"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-chaincode.git/model"
)

func SummaryPersonasID(pers []model.Persona) (min uint64, max uint64, index map[uint64]model.Persona, ids []uint64) {
	min = uint64(math.MaxUint64)
	max = uint64(0)
	index = map[uint64]model.Persona{}
	ids = []uint64{}
	for _, per := range pers {
		if per.ID < min {
			min = per.ID
		}
		if per.ID > max {
			max = per.ID
		}
		index[per.ID] = per
	}
	for id := range index {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return
}
