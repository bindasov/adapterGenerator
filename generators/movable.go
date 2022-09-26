package generators

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/fatih/camelcase"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type MethodName struct {
	Action   string
	Property string
}

type Method struct {
	Name       *MethodName
	InputParam string
	Output     string
}

func extractMethodParams(s string) string {
	i := strings.Index(s, "(")
	if i >= 0 {
		j := strings.Index(s, ")")
		if j >= 0 {
			return s[i+1 : j]
		}
	}
	return ""
}

func extractOutput(s string) string {
	i := strings.Index(s, ")")
	if i >= 0 {
		j := len(s)
		if j >= i {
			return s[i+1 : j]
		}
	}
	return ""
}

func processTemplate(movableTemplate string, data []*Method) string {
	tmpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Parse(movableTemplate))
	var processed bytes.Buffer
	err := tmpl.Execute(&processed, data)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}
	return string(formatted)
}

func Generate(adapterType interface{}) string {
	var methods []*Method

	t := adapterType.(reflect.Type)
	for i := 0; i < t.NumMethod(); i++ {
		splitted := camelcase.Split(t.Method(i).Name)
		method := &Method{
			Name: &MethodName{
				Action:   splitted[0],
				Property: splitted[1],
			},
			InputParam: extractMethodParams(t.Method(i).Type.String()),
			Output:     extractOutput(t.Method(i).Type.String()),
		}
		methods = append(methods, method)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	movableTemplate, err := os.ReadFile(filepath.Join(path, "movableAdapter.tmpl"))
	if err != nil {
		fmt.Print(err)
	}
	formatted := processTemplate(string(movableTemplate), methods)

	f, _ := os.Create(filepath.Join(path, "adapters/movable.go"))
	w := bufio.NewWriter(f)
	w.WriteString(formatted)
	w.Flush()
	return formatted
}
