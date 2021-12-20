package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var ErrUnknownFormatRequested error = fmt.Errorf("unknown output format requested")

type Printer struct {
	noListWithSingleEntry bool
}

type Opt func(p *Printer)

func NewPrinter(opts ...Opt) Printer {
	p := Printer{}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func NoListWithSingleEntry() Opt {
	return func(p *Printer) {
		p.noListWithSingleEntry = true
	}
}

func (p Printer) Print(format string, out io.Writer, tableFunc func() error, nameFunc func() error, items interface{}) error {
	if p.noListWithSingleEntry && reflect.TypeOf(items).Kind() == reflect.Slice && reflect.ValueOf(items).Len() == 1 {
		items = reflect.ValueOf(items).Index(0).Interface()
	}

	switch {
	case format == "json":
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "    ")
		if err := enc.Encode(items); err != nil {
			return fmt.Errorf("error encoding items: %w", err)
		}
		buf.WriteTo(out)
		out.Write([]byte("\n"))
		return nil
	case format == "name":
		return nameFunc()
	case format == "table", format == "":
		return tableFunc()
	case strings.HasPrefix(format, "go-template="):
		tmplString := strings.TrimPrefix(format, "go-template=")
		tmpl, err := template.New("").Parse(tmplString)
		if err != nil {
			return fmt.Errorf("template parsing error: %s", err)
		}
		err = tmpl.Execute(out, items)
		if err != nil {
			return fmt.Errorf("template parsing error: %s", err)
		}
		return nil
	default:
		return ErrUnknownFormatRequested
	}
}

func AddFlag(cmd *cobra.Command) *string {
	return cmd.PersistentFlags().StringP("output", "o", "", "Output format. One of: json|name|table|go-template")
}
