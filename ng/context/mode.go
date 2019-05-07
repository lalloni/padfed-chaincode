package context

type mode uint8

const (
	Development mode = iota
	Testing
	Production
)
