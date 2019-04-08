package archiver

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

func NewTarXZ(target string) (Archiver, error) {
	return NewTar(target, xzCompressor)
}

func NewTarGZ(target string) (Archiver, error) {
	return NewTar(target, gzCompressor)
}

func NewTar(target string, filter func(io.WriteCloser) (io.WriteCloser, error)) (Archiver, error) {
	w, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	ww, err := filter(w)
	if err != nil {
		w.Close()
		return nil, err
	}
	tw := tar.NewWriter(ww)
	return &tarArchiver{
		w:  w,
		tw: tw,
	}, nil
}

type tarArchiver struct {
	w  io.WriteCloser
	tw *tar.Writer
}

func (a *tarArchiver) AddAll(files []string, namer func(string) string) error {
	return archiveFiles(a, files, namer)
}

func (a *tarArchiver) Add(file string, name string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	h, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	h.Name = name
	h.Uid = 0
	h.Gid = 0
	err = a.tw.WriteHeader(h)
	if err != nil {
		return err
	}
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(a.tw, in)
	return err
}

func (a *tarArchiver) Close() error {
	err := a.tw.Close()
	if err != nil {
		a.w.Close()
		return err
	}
	return a.w.Close()
}

func xzCompressor(w io.WriteCloser) (io.WriteCloser, error) {
	return xz.NewWriter(w)
}

func gzCompressor(w io.WriteCloser) (io.WriteCloser, error) {
	return gzip.NewWriter(w), nil
}
