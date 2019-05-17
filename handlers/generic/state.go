package generic

import (
	"time"
	"unicode/utf8"

	"github.com/hyperledger/fabric/protos/ledger/queryresult"
)

type state struct {
	Key      string      `json:"key,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
}

func newstate(key string, content interface{}) *state {
	s := &state{
		Key:     key,
		Content: content,
	}
	if bs, ok := content.([]byte); ok {
		if utf8.Valid(bs) {
			s.Content = string(bs)
		} else {
			s.Encoding = "base64"
		}
	}
	return s
}

type statehistory struct {
	Key     string      `json:"key,omitempty"`
	History []*statemod `json:"history,omitempty"`
}

type statemod struct {
	TxID     string      `json:"txid,omitempty"`
	Time     time.Time   `json:"time,omitempty"`
	Delete   bool        `json:"delete,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
}

func newstatemod(s *queryresult.KeyModification) *statemod {
	var con interface{}
	enc := ""
	if s.Value != nil {
		if utf8.Valid(s.Value) {
			con = string(s.Value)
		} else {
			enc = "base64"
			con = s.Value
		}
	}
	return &statemod{
		TxID:     s.TxId,
		Time:     time.Unix(s.Timestamp.Seconds, int64(s.Timestamp.Nanos)).In(time.Local),
		Delete:   s.IsDelete,
		Content:  con,
		Encoding: enc,
	}
}
