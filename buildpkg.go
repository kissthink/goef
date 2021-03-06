// Package goef helps you embed file/data/assets/resources/binary directly in Go
// code. There are many tools can help you embed too, but this package tries to
// be with minimal features (files are read-only) and easy to use.
package goef

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const gofile = `package {{.PkgName}}

import (
	"encoding/base64"
	"os"
)

var virtualFilesystem = map[string]string{
{{ range .Files }}"{{ .Name }}": "{{ .Base64Content }}",
{{ end }}}

func ReadFile(filename string) ([]byte, error) {
	content, ok := virtualFilesystem[filename]
	if ok {
		return base64.StdEncoding.DecodeString(content)
	}
	return nil, os.ErrNotExist
}
`

type pkgData struct {
	PkgName string
	Files   []pkgFile
}

type pkgFile struct {
	Name          string
	Base64Content string
}

func getFilenameContent(dirpath, path string, info os.FileInfo) (name, content string, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	content = base64.StdEncoding.EncodeToString(b)
	name, err = filepath.Rel(dirpath, path)
	return
}

// GenerateGoPackage will generate a single Go file which contains the files in
// *dirpath* directory, and the name of the package is *pkgname*.
//
// You can put the generated Go file in your source code, and read the embedded
// files with the following method:
//
//   ReadFile(filename string) ([]byte, error)
//
// The usage of the above method is the same as ioutil.ReadFile in Go standard
// library.
//
// You can also put the generated Go file in a separate package, import and read
// embedded files in the same way.
func GenerateGoPackage(pkgname, dirpath, outputpath string) (err error) {
	fo, err := os.Create(outputpath)
	if err != nil {
		return
	}
	defer fo.Close()

	pd := pkgData{PkgName: pkgname}
	err = filepath.Walk(dirpath, func(filepath string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.Mode().IsRegular() {
			name, content, errf := getFilenameContent(dirpath, filepath, info)
			if errf != nil {
				return errf
			}

			pd.Files = append(pd.Files,
				pkgFile{
					Name:          name,
					Base64Content: content,
				})
		}

		return nil
	})
	if err != nil {
		return
	}

	tmpl, err := template.New("goembed").Parse(gofile)
	if err != nil {
		return
	}

	err = tmpl.Execute(fo, pd)
	return
}
