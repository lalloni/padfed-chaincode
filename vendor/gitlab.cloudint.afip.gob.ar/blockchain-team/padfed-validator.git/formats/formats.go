package formats

type FormatCheckerFunc func(input interface{}) bool

func (f FormatCheckerFunc) IsFormat(input interface{}) bool {
	return f(input)
}
