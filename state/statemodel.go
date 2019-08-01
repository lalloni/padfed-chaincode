package state

import (
	"time"
	"unicode/utf8"

	"github.com/hyperledger/fabric/protos/ledger/queryresult"
)

type State struct {
	Key      string      `json:"key,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
	Nil      bool        `json:"-"`
}

func New(key string, content []byte) *State {
	s := &State{Key: key}
	if content == nil {
		s.Nil = true
		return s
	}
	if utf8.Valid(content) {
		s.Content = string(content)
	} else {
		s.Encoding = "base64"
		s.Content = content
	}
	return s
}

type statehistory struct {
	Key     string      `json:"key,omitempty"`
	History []*statemod `json:"history,omitempty"`
}

type statemod struct {
	Version  uint64      `json:"version,omitempty"`
	TxID     string      `json:"txid,omitempty"`
	Block    uint64      `json:"block,omitempty"`
	Time     string      `json:"time,omitempty"`
	Delete   bool        `json:"delete,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Encoding string      `json:"encoding,omitempty"`
}

func newstatemod(s *queryresult.KeyModification, block uint64, version uint64) *statemod {
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
		Version:  version,
		TxID:     s.TxId,
		Block:    block,
		Time:     time.Unix(s.Timestamp.Seconds, int64(s.Timestamp.Nanos)).In(time.Local).Format(time.RFC3339Nano),
		Delete:   s.IsDelete,
		Content:  con,
		Encoding: enc,
	}
}
