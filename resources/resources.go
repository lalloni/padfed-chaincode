package resources

import (
	"github.com/gobuffalo/packr/v2"
)

var Schemas = packr.New("schemas", "./schemas")

var Data = packr.New("data", "./data")
