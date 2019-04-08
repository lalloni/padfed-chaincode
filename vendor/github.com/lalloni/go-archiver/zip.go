package archiver

import (
	"archive/zip"
	"io"
	"os"
)

func NewZip(target string) (Archiver, error) {
	w, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	zw := zip.NewWriter(w)
	return &zipArchiver{
		w:  w,
		zw: zw,
	}, nil
}

type zipArchiver struct {
	w  io.WriteCloser
	zw *zip.Writer
}

func (a *zipArchiver) AddAll(files []string, namer func(string) string) error {
	return archiveFiles(a, files, namer)
}

func (a *zipArchiver) Add(file string, name string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	h, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	h.Name = name
	h.Method = zip.Deflate
	w, err := a.zw.CreateHeader(h)
	if err != nil {
		return err
	}
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(w, in)
	return err
}

func (a *zipArchiver) Close() error {
	err := a.zw.Close()
	if err != nil {
		a.w.Close()
		return err
	}
	return a.w.Close()
}
