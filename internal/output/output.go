package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	fmt.Fprintln(Stdout, v)
}

func PrintJSON(v any) {
	enc := json.NewEncoder(Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
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
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, h)
	}
	fmt.Fprintln(w)

	for _, row := range rows {
		for i, col := range row {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, col)
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func Error(format string, args ...any) {
	fmt.Fprintf(Stderr, "error: "+format+"\n", args...)
}

func Success(format string, args ...any) {
	fmt.Fprintf(Stdout, format+"\n", args...)
}
