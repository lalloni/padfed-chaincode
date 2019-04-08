package archiver

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var UnknownFormat = errors.New("unknown format")

type format string

var (
	TarXZFormat format = "tar.xz"
	TarGZFormat format = "tar.gz"
	ZipFormat   format = "zip"
)

func (f *format) AddExt(name string) string {
	return name + "." + string(*f)
}

func (f *format) ChangeExt(name string) string {
	return f.AddExt(strings.TrimSuffix(name, filepath.Ext(name)))
}

func (f *format) NewArchiver(target string) (Archiver, error) {
	switch *f {
	case TarXZFormat:
		return NewTarXZ(target)
	case TarGZFormat:
		return NewTarGZ(target)
	case ZipFormat:
		return NewZip(target)
	default:
		return nil, UnknownFormat
	}
}

func Format(s string) (format, error) {
	switch s {
	case string(TarGZFormat):
		return TarGZFormat, nil
	case string(TarXZFormat):
		return TarXZFormat, nil
	case string(ZipFormat):
		return ZipFormat, nil
	default:
		return "", UnknownFormat
	}
}

func FormatDefault(s string, def format) (format, error) {
	if s == "" {
		return def, nil
	}
	return Format(s)
}
