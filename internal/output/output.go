package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"text/tabwriter"
)

var (
	JSONOutput bool
	Stdout     io.Writer = os.Stdout
	Stderr     io.Writer = os.Stderr
)

func Print(v any) {
	if JSONOutput {
		enc := json.NewEncoder(Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(v)
		return
	}
	_, _ = fmt.Fprintln(Stdout, v)
}

func PrintJSON(v any) {
	enc := json.NewEncoder(Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(normalizeEmpty(v))
}

// normalizeEmpty converts nil slices and maps to empty ones so they
// marshal as [] and {} instead of null.
func normalizeEmpty(v any) any {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return v
	}
	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			return reflect.MakeSlice(rv.Type(), 0, 0).Interface()
		}
	case reflect.Map:
		if rv.IsNil() {
			return reflect.MakeMap(rv.Type()).Interface()
		}
	}
	return v
}

func PrintTable(headers []string, rows [][]string) {
	if JSONOutput {
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			m := make(map[string]string)
			for j, h := range headers {
				if j < len(row) {
					m[h] = row[j]
				}
			}
			data[i] = m
		}
		PrintJSON(data)
		return
	}

	w := tabwriter.NewWriter(Stdout, 0, 0, 2, ' ', 0)
	for i, h := range headers {
		if i > 0 {
			_, _ = fmt.Fprint(w, "\t")
		}
		_, _ = fmt.Fprint(w, h)
	}
	_, _ = fmt.Fprintln(w)

	for _, row := range rows {
		for i, col := range row {
			if i > 0 {
				_, _ = fmt.Fprint(w, "\t")
			}
			_, _ = fmt.Fprint(w, col)
		}
		_, _ = fmt.Fprintln(w)
	}
	_ = w.Flush()
}

func Error(format string, args ...any) {
	_, _ = fmt.Fprintf(Stderr, "error: "+format+"\n", args...)
}

func Success(format string, args ...any) {
	_, _ = fmt.Fprintf(Stdout, format+"\n", args...)
}
