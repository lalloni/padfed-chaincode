// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lalloni/go-archiver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/build"
)

// Limpia directorio de proyecto de artefactos temporales generados
func Clean() {
	sh.Rm("target")
}

// Genera JSON Schema
func Genschema() error {
	return build.Convert("schemas/resources", "doc/schemas")
}

// Genera recursos embebidos en código fuente
func Genpack() error {
	return build.RunPackr()
}

// Genera todos los "generables"
func Genall() error {
	mg.SerialDeps(Genschema, Genpack)
	return nil
}

// Ejecuta tests
func Test() error {
	mg.Deps(Genall)
	return sh.RunV("go", "test", "./...")
}

// Ejecuta análisis estático de código fuente
func Check() error {
	return build.RunLinter("run")
}

// Ejecuta compilación de librería de validación
func CompileLibrary() error {
	mg.Deps(Genall)
	return sh.Run("go", "build", "./...")
}

// Ejecuta compilación de herramienta de validación
func CompileValidatorTool() error {
	mg.Deps(Genall)
	base := "target/bin/"
	for _, goos := range []string{"windows", "linux"} {
		for _, goarch := range []string{"amd64"} {
			out := filepath.Join(base, goos+"-"+goarch, "validator")
			if goos == "windows" {
				out = out + ".exe"
			}
			env := map[string]string{
				"GOOS":   goos,
				"GOARCH": goarch,
			}
			err := sh.RunWith(env, "go", "build", "-o", out, "./cmd/validator")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Ejecuta todas las tareas de compilcación
func Compile() error {
	mg.Deps(CompileLibrary, CompileValidatorTool)
	return nil
}

// Empaqueta los binarios del proyecto
func Package() error {
	mg.Deps(Clean, CompileValidatorTool)
	p := "target/pkg/validator.zip"
	log.Printf("packaging binaries into %s", p)
	os.MkdirAll(filepath.Dir(p), 0777)
	a, err := archiver.NewZip(p)
	if err != nil {
		return err
	}
	defer a.Close()
	fs, err := filepath.Glob("target/bin/*/*")
	if err != nil {
		return err
	}
	err = a.AddAll(fs, func(n string) string {
		d, f := filepath.Split(n)
		d = filepath.Base(d)
		nn := filepath.Join(d, f)
		log.Printf("adding %s as %s", n, nn)
		return nn
	})
	if err != nil {
		return err
	}
	return a.Close()
}

// Lanza GoConvey (http://goconvey.co/)
func Convey() error {
	err := build.RunPackr("clean")
	if err != nil {
		return err
	}
	return build.RunGoConvey("-port=9999", "-watchedSuffixes=.go,.yaml", "-packages=1")
}

// Ejecuta el proceso de release
func Release() error {
	version := os.Getenv("version")
	if version == "" {
		return errors.New(`Version is required for release.
You must set the version to be released using the environment variable 'version'.
On unix-like shells you could do something like:
    env version=1.2.3 mage release`)
	}
	fmt.Printf("Releasing version: %s\n", version)
	mg.SerialDeps(Genall, Check, Compile, Test)
	return errors.New("still not implemented")
}

// Construye un binario estático de este build
func Buildbuild() error {
	return sh.RunV("mage", "-compile", "magestatic")
}
