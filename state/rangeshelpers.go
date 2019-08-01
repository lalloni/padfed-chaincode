package state

import (
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/pkg/errors"
)

func QueryKeyRanges(ctx *context.Context, query *Ranges) (interface{}, error) {
	if query.IsSingle() {
		s, err := queryKeyItem(ctx, query.Single())
		if err != nil {
			return nil, errors.Wrap(err, "querying single item")
		}
		if query.Single().IsPoint() {
			// caso especial: si se pide getstate(x) se devuelve solo <valor> de
			// contenido en lugar de {key:x,content:<valor>}
			st := s.(*State)
			if st.Nil {
				return nil, nil
			}
			return st.Content, nil
		}
		return s, nil
	}
	result := []interface{}{}
	for _, item := range query.List() {
		r, err := queryKeyItem(ctx, item)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}

func queryKeyItem(ctx *context.Context, item *Item) (interface{}, error) {
	if item.IsPoint() {
		s, err := keyState(ctx, item.Point())
		if err != nil {
			return nil, errors.Wrap(err, "executing single point query")
		}
		return s, nil
	}
	begin, until := item.Range()
	ss, err := keyRangeStates(ctx, begin, until)
	if err != nil {
		return nil, errors.Wrap(err, "executing ranges range query")
	}
	return ss, nil
}

func processKeyRanges(ctx *context.Context, query *Ranges, process func(key string) (interface{}, error)) (interface{}, error) {
	if query.IsSingle() {
		r, err := processKeyItem(ctx, query.Single(), process)
		if err != nil {
			return nil, errors.Wrap(err, "processing single item")
		}
		return r, nil
	}
	rr := []interface{}{}
	for _, item := range query.List() {
		r, err := processKeyItem(ctx, item, process)
		if err != nil {
			return nil, errors.Wrap(err, "processing list item")
		}
		rr = append(rr, r)
	}
	return rr, nil
}

func processKeyItem(ctx *context.Context, item *Item, process func(key string) (interface{}, error)) (interface{}, error) {
	if item.IsPoint() {
		r, err := process(item.Point())
		if err != nil {
			return nil, errors.Wrap(err, "processing point key")
		}
		return r, nil
	}
	begin, until := item.Range()
	ss, err := keyRangeStates(ctx, begin, until)
	if err != nil {
		return nil, errors.Wrap(err, "executing range query")
	}
	rr := []interface{}{}
	for _, s := range ss {
		r, err := process(s.Key)
		if err != nil {
			return nil, errors.Wrap(err, "processing range key")
		}
		rr = append(rr, r)
	}
	return rr, nil
}
