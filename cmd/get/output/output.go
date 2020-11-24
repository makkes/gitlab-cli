package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var ErrUnknownFormatRequested error = fmt.Errorf("unknown output format requested")

func Print(b []byte, format string, out io.Writer, tableFunc func() error, nameFunc func() error, items []interface{}) error {
	switch {
	case format == "json":
		var buf bytes.Buffer
		json.Indent(&buf, b, "", "    ")
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
		for _, p := range items {
			err = tmpl.Execute(out, p)
			if err != nil {
				return fmt.Errorf("template parsing error: %s", err)
			}
			fmt.Fprintln(out)
		}
		return nil
	default:
		return ErrUnknownFormatRequested
	}
}

func AddFlag(cmd *cobra.Command) *string {
	return cmd.PersistentFlags().StringP("output", "o", "", "Output format. One of: json|name|table|go-template")
}
