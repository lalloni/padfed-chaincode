package support

import "github.com/hyperledger/fabric/protos/peer"

type Handler func(*Context) peer.Response

type HandlerMap map[string]Handler
