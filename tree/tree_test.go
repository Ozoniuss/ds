package tree

import (
	"math/rand"
	"testing"
)

func TestInsertDeletes(t *testing.T) {
	r := rand.New(rand.NewSource(69))

	type testcase struct {
		name string
		tree Tree[int]
	}

	testcases := []testcase{
		{
			name: "bst",
			tree: NewBST[int](),
		},
		{
			name: "rbt",
			tree: NewRBT[int](),
		},
	}

	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// check that the tree works as a new one once everything was deleted
			for range 2 {
				existing := map[int]struct{}{}
				existing[0] = struct{}{}
				iterations := 1000

				for range iterations {
					v := 0
					var ok bool
				LOOP:
					for {
						if _, ok = existing[v]; ok {
							v = r.Int()
						} else {
							existing[v] = struct{}{}
							tc.tree.Insert(v)
							break LOOP
						}
					}
				}

				if tc.tree.Size() != iterations {
					t.Fatalf("expected %d elements, got %d\n", iterations, tc.tree.Size())
				}

				for k := range existing {
					if k == 0 {
						continue
					}
					err := tc.tree.Delete(k)
					if err != nil {
						t.Errorf("expected %d to be in tree, got error %s\n", k, err.Error())
					}
				}

				if tc.tree.Size() != 0 {
					t.Errorf("expected size 0 after all deletions, got %d", tc.tree.Size())
				}

				if tc.tree.Root() != nil {
					t.Errorf("root should be nil after all deletions, got %#v", tc.tree.Root())
				}
			}

		})

	}
}
