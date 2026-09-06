// This file is AI-generated.
//
// go run ./rbtree/cmd/generateExhaustiveFormat -n=6 -out=rbtree/testdata/exhaustiveFormat.txt
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Ozoniuss/ds/rbtree"
)

func permute(values []int) [][]int {
	if len(values) == 0 {
		return [][]int{{}}
	}
	var result [][]int
	for i := range values {
		rest := make([]int, 0, len(values)-1)
		rest = append(rest, values[:i]...)
		rest = append(rest, values[i+1:]...)
		for _, p := range permute(rest) {
			perm := append([]int{values[i]}, p...)
			result = append(result, perm)
		}
	}
	return result
}

func main() {
	n := flag.Int("n", 6, "length of the list [1..n] to permute")
	out := flag.String("out", "rbtree/testdata/exhaustiveFormat.txt", "output file path")
	flag.Parse()

	values := make([]int, *n)
	for i := range values {
		values[i] = i + 1
	}

	perms := permute(values)

	var sb strings.Builder
	for _, perm := range perms {
		strs := make([]string, len(perm))
		for i, v := range perm {
			strs[i] = fmt.Sprintf("%d", v)
		}
		fmt.Fprintf(&sb, "input: [%s]\n", strings.Join(strs, ","))

		tr := rbtree.NewRBT[int]()
		for _, v := range perm {
			tr.Insert(v)
		}

		lines := rbtree.BuildLines(tr.Root(), 2, true, true)
		for _, l := range lines {
			sb.WriteString(l)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	if err := os.WriteFile(*out, []byte(sb.String()), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write output file:", err)
		os.Exit(1)
	}
}
