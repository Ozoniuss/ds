package tree

import (
	"strings"
	"testing"
)

func TestFormatLinuxTree(t *testing.T) {

	type testCase struct {
		t    Tree[int]
		repr string
	}

	testcases := []testCase{
		{
			t:    nil,
			repr: `nil tree`,
		},
		{
			t:    NewBST[int](),
			repr: `empty tree`,
		},
		{
			t: func() Tree[int] {
				tr := NewBST[int]()
				tr.Insert(5)
				tr.Insert(2)
				tr.Insert(10)
				tr.Insert(8)
				tr.Insert(12)
				tr.Insert(6)
				tr.Insert(9)
				return tr
			}(),
			repr: `
5
├── 2
└── 10
    ├── 8
    │   ├── 6
    │   └── 9
    └── 12`,
		},
		{
			t: func() Tree[int] {
				tr := NewBST[int]()
				tr.Insert(4)
				tr.Insert(12)
				tr.Insert(2)
				tr.Insert(1)
				tr.Insert(8)
				tr.Insert(13)
				tr.Insert(6)
				tr.Insert(9)
				tr.Insert(5)
				tr.Insert(11)
				return tr
			}(),
			repr: `
4
├── 2
│   ├── 1
│   └── *
└── 12
    ├── 8
    │   ├── 6
    │   │   ├── 5
    │   │   └── *
    │   └── 9
    │       ├── *
    │       └── 11
    └── 13`,
		},
	}

	for _, tc := range testcases {
		out := FormatTree(tc.t, string(FormatLinuxTree))
		if out != strings.TrimPrefix(tc.repr, "\n") {
			t.Errorf("expected:\n%s (len: %d)\ngot:\n%s (len: %d)", tc.repr, len(tc.repr), out, len(out))
		}
	}
}
