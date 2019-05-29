// +build mage

package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/lalloni/go-archiver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	prefix "github.com/x-cray/logrus-prefixed-formatter"

	"gitlab.cloudint.afip.gob.ar/blockchain-team/padfed-validator.git/build"
)

func init() {
	log.SetFormatter(&prefix.TextFormatter{
		FullTimestamp: true,
	})
	log.Info("magefile initialized")
}

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
	mg.Deps(Genall)
	return build.RunLinter("run")
}

// Ejecuta análisis estático de código fuente y tests
func Verify() {
	mg.Deps(Check, Test)
}

// Ejecuta compilación de librería de validación
func Compilelibrary() error {
	mg.Deps(Genall)
	return sh.Run("go", "build", "./...")
}

// Ejecuta compilación de herramienta de validación
func Compilevalidatortool() error {
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
	mg.Deps(Compilelibrary, Compilevalidatortool)
	return nil
}

// Empaqueta los binarios del proyecto
func Package() error {
	mg.Deps(Clean, Compilevalidatortool)
	p := "target/pkg/validator.zip"
	log.Infof("packaging binaries into %s", p)
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
		log.Infof("adding %s as %s", n, nn)
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
	log.Info("checking parameters")
	version := os.Getenv("ver")
	if version == "" {
		return errors.New(`Version is required for release.
You must set the version to be released using the environment variable 'ver'.
On unix-like shells you could do something like:
    env ver=1.2.3 mage release`)
	}
	if _, err := semver.NewVersion(version); err != nil {
		return errors.Wrapf(err, "checking syntax of version %q", version)
	}

	tag := "v" + version
	log.Infof("releasing version %s with tag %s", version, tag)

	log.Info("checking release tag does not exist")
	out, err := sh.Output("git", "tag")
	if err != nil {
		return errors.Wrap(err, "getting git tags")
	}
	s := bufio.NewScanner(strings.NewReader(out))
	for s.Scan() {
		if tag == s.Text() {
			return errors.Errorf("release tag %q already exists", tag)
		}
	}

	log.Info("updating generated resources")
	mg.SerialDeps(Genall)

	log.Info("checking working tree is not dirty")
	out, err = sh.Output("git", "status", "-s")
	if err != nil {
		return errors.Wrap(err, "getting git status")
	}
	if len(out) > 0 {
		return errors.Errorf("working directory is dirty")
	}

	log.Info("running linter, compiler & tests")
	mg.Deps(Compile, Check, Test)

	log.Infof("creating tag %s", tag)
	if err := sh.RunV("git", "tag", "-s", "-m", "Release "+version, tag); err != nil {
		return errors.Wrap(err, "creating git tag")
	}

	log.Infof("pushing tag %s to 'origin' remote", tag)
	if err := sh.RunV("git", "push", "origin", tag); err != nil {
		return errors.Wrap(err, "pushing tag to origin remote")
	}

	log.Infof("pushing current branch to 'origin' remote", tag)
	if err := sh.RunV("git", "push", "origin"); err != nil {
		return errors.Wrap(err, "pushing current branch to origin remote")
	}

	log.Info("release successfuly completed")

	return nil
}

// Construye un binario estático de este build
func Buildbuild() error {
	return sh.RunV("mage", "-compile", "magestatic")
}

// Ejecuta los tests ante cambios en el proyecto
func Testwatch() error {
	c := make(chan build.Event, 1000)
	err := build.Monitor(".", c, "-.*/**", "-target", "-target/**")
	if err != nil {
		return err
	}
	log.Info("Running tests for the first time...")
	if Test() == nil {
		log.Info("SUCCESS")
	} else {
		log.Error("FAILED")
	}
	for e := range c {
		log.Infof("Running tests after receiving %s...", e.String())
		if Test() == nil {
			log.Info("SUCCESS")
		} else {
			log.Error("FAILED")
		}
	}
	return nil
}
