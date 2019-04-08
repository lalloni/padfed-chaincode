package archiver

import (
	"path/filepath"

	"github.com/pkg/errors"
)

func Package(targetFormat format, target string, files []string) error {
	archiver, err := targetFormat.NewArchiver(target)
	if err != nil {
		return errors.Wrapf(err, "creating %s archiver", targetFormat)
	}
	err = archiver.AddAll(files, filepath.Base)
	if err != nil {
		defer archiver.Close()
		return errors.Wrapf(err, "adding files to %s archiver", targetFormat)
	}
	return archiver.Close()
}

type Archiver interface {
	Add(file string, name string) error
	AddAll(files []string, namer func(string) string) error
	Close() error
}

func archiveFiles(archiver Archiver, files []string, namer func(string) string) error {
	for _, file := range files {
		err := archiver.Add(file, namer(file))
		if err != nil {
			return errors.Wrapf(err, "adding file '%s'", file)
		}
	}
	return nil
}
