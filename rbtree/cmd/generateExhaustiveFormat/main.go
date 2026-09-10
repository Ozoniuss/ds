// This file is AI-generated.
//
// go run ./rbtree/cmd/generateExhaustiveFormat
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Ozoniuss/ds/rbtree"
)

func permute(values []string) [][]string {
	if len(values) == 0 {
		return [][]string{{}}
	}
	var result [][]string
	for i := range values {
		rest := make([]string, 0, len(values)-1)
		rest = append(rest, values[:i]...)
		rest = append(rest, values[i+1:]...)
		for _, p := range permute(rest) {
			perm := append([]string{values[i]}, p...)
			result = append(result, perm)
		}
	}
	return result
}

func main() {

	opts := []struct {
		out   string
		shape int
	}{
		{
			out:   "rbtree/testdata/exhaustiveFormat1.txt",
			shape: rbtree.SHAPE_TREE,
		},
		{
			out:   "rbtree/testdata/exhaustiveFormat2.txt",
			shape: rbtree.SHAPE_SQUARE,
		},
	}

	values := []string{"552345623", "6324232", "2", "4354", "1999059330313454643", "32"}

	perms := permute(values)

	for _, opt := range opts {

		out := opt.out
		shape := opt.shape

		var sb strings.Builder
		for _, perm := range perms {
			strs := make([]string, len(perm))
			for i, v := range perm {
				strs[i] = fmt.Sprintf("%s", v)
			}
			fmt.Fprintf(&sb, "input: [%s]\n", strings.Join(strs, ","))

			tr := rbtree.NewRBT[string]()
			for _, v := range perm {
				tr.Insert(v)
			}

			lines := rbtree.BuildLines(tr.Root(), shape, 2, true, true)
			for _, l := range lines {
				sb.WriteString(l)
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}

		if err := os.WriteFile(out, []byte(sb.String()), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "failed to write output file:", err)
			os.Exit(1)
		}
	}

}
