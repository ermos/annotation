package annotation

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

type Result []Annotation

type Annotation struct {
	Method  string
	Key 	string
	Data 	string
}

// Parse controller dir for extract all annotation and make a json file of their.
func Fetch (dir string, v interface{}, resultFunc func(reflect.Value, Result) error) error {
	var result Result
	// Get all comments
	fset := token.NewFileSet()
	d, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, dir := range d {
		for _, file := range dir.Files {
			for _, r := range file.Decls {
				fn, ok := r.(*ast.FuncDecl)
				if !ok || fn.Doc.Text() == "" {
					continue
				}
				method := fn.Name.String()
				annotations := strings.Split(fn.Doc.Text(), "\n")
				for _, annotation := range annotations {
					key, data := extractAnnotation(annotation)
					if key != "" {
						result = append(result, Annotation{
							Method: method,
							Key: key,
							Data: data,
						})
					}
				}
			}
		}
	}
	// Read interface
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("pointer is not set or interface is nil")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Slice {
		return errors.New("slice needed")
	}
	// Run custom func
	return resultFunc(rv, result)
}

func cleanAnnotations(value string) string {
	res := strings.TrimSpace(value)
	rgx := regexp.MustCompile(`\s+`)
	res = rgx.ReplaceAllString(res, " ")
	return res
}

func extractAnnotation(value string) (name string, data string) {
	value = cleanAnnotations(value)
	rgx := regexp.MustCompile(`^@(.*)\((.*)\)$`)
	if !rgx.Match([]byte(value)) {
		return "", ""
	}
	res := rgx.FindStringSubmatch(value)
	return res[1], res[2]
}