package tree

import (
	"testing"
)

// mustPanic is used to verify that operations on uninitialized trees or
// nodes actually panic.
func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	fn()
}

func TestRBTUninitialized(t *testing.T) {
	t.Parallel()

	type method struct {
		name string
		fn   func(*RBT[int])
	}

	methods := []method{
		{name: "Root", fn: func(tr *RBT[int]) { tr.Root() }},
		{name: "Size", fn: func(tr *RBT[int]) { tr.Size() }},
		{name: "Count", fn: func(tr *RBT[int]) { tr.Count(0) }},
		{name: "Insert", fn: func(tr *RBT[int]) { _ = tr.Insert(1) }},
		{name: "Delete", fn: func(tr *RBT[int]) { _ = tr.Delete(1) }},
		{name: "String", fn: func(tr *RBT[int]) { _ = tr.String() }},
	}

	trees := []struct {
		name string
		tree *RBT[int]
	}{
		{name: "nilReceiver", tree: nil},
		{name: "zeroValue", tree: &RBT[int]{}},
		{name: "missingSentinel", tree: &RBT[int]{size: 0, root: nil, tnil: nil}},
	}

	for _, tr := range trees {
		for _, m := range methods {
			t.Run(tr.name+"/"+m.name, func(t *testing.T) {
				t.Parallel()
				mustPanic(t, func() { m.fn(tr.tree) })
			})
		}
	}
}

func TestRBTEmpty(t *testing.T) {
	t.Parallel()

	tr := NewRBT[int]()

	if tr.Size() != 0 {
		t.Errorf("empty tree size: got %d, want 0", tr.Size())
	}
	if tr.Root() != nil {
		t.Errorf("empty tree root: got %#v, want nil", tr.Root())
	}
	if got := tr.Count(0); got != 0 {
		t.Errorf("Count(0) on empty tree: got %d, want 0", got)
	}
	if got := tr.Count(1); got != 0 {
		t.Errorf("Count(1) on empty tree: got %d, want 0", got)
	}

	err := tr.Delete(1)
	if err == nil {
		t.Fatal("Delete on empty tree: got nil error, want value not found")
	}
	if err.Error() != "value not found" {
		t.Errorf("Delete on empty tree: got %q, want %q", err.Error(), "value not found")
	}
	if tr.Size() != 0 {
		t.Errorf("Delete on empty tree changed size: got %d, want 0", tr.Size())
	}
	if tr.Root() != nil {
		t.Errorf("Delete on empty tree set a root: got %#v", tr.Root())
	}

	if got := tr.String(); got != "empty tree" {
		t.Errorf("String on empty tree: got %q, want %q", got, "empty tree")
	}
}

func TestRBTInsertIntoEmpty(t *testing.T) {
	t.Parallel()

	tr := NewRBT[int]()
	if err := tr.Insert(0); err != nil {
		t.Fatalf("Insert(0) into empty tree: %v", err)
	}

	if tr.Size() != 1 {
		t.Errorf("size after first insert: got %d, want 1", tr.Size())
	}
	root := tr.Root()
	if root == nil {
		t.Fatal("root after first insert is nil")
	}
	if root.Value() != 0 {
		t.Errorf("root value: got %d, want 0", root.Value())
	}
	if root.Count() != 1 {
		t.Errorf("root count: got %d, want 1", root.Count())
	}
	if root.Parent() != nil {
		t.Errorf("root parent: got %#v, want nil", root.Parent())
	}
	if root.Left() != nil {
		t.Errorf("root left: got %#v, want nil", root.Left())
	}
	if root.Right() != nil {
		t.Errorf("root right: got %#v, want nil", root.Right())
	}
	if got := tr.Count(0); got != 1 {
		t.Errorf("Count(0): got %d, want 1", got)
	}
	if got := tr.Count(1); got != 0 {
		t.Errorf("Count(1) missing value: got %d, want 0", got)
	}
}

func TestRBTDeleteOnlyNode(t *testing.T) {
	t.Parallel()

	tr := NewRBT[int]()
	if err := tr.Insert(7); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	if err := tr.Delete(7); err != nil {
		t.Fatalf("Delete only node: %v", err)
	}
	if tr.Size() != 0 {
		t.Errorf("size after deleting only node: got %d, want 0", tr.Size())
	}
	if tr.Root() != nil {
		t.Errorf("root after deleting only node: got %#v, want nil", tr.Root())
	}
	if got := tr.Count(7); got != 0 {
		t.Errorf("Count after deleting only node: got %d, want 0", got)
	}

	err := tr.Delete(7)
	if err == nil {
		t.Fatal("second Delete of only node: got nil error")
	}
	if err.Error() != "value not found" {
		t.Errorf("second Delete: got %q, want %q", err.Error(), "value not found")
	}

	if err := tr.Insert(7); err != nil {
		t.Fatalf("Insert into tree emptied by Delete: %v", err)
	}
	if tr.Size() != 1 || tr.Root() == nil || tr.Root().Value() != 7 {
		t.Fatalf("tree not usable after emptying: size=%d root=%#v", tr.Size(), tr.Root())
	}
}

func TestRBTDuplicateInsert(t *testing.T) {
	t.Parallel()

	tr := NewRBT[int]()
	if err := tr.Insert(3); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	err := tr.Insert(3)
	if err == nil {
		t.Fatal("duplicate Insert: got nil error")
	}
	if err.Error() != "value already exists" {
		t.Errorf("duplicate Insert: got %q, want %q", err.Error(), "value already exists")
	}
	if tr.Size() != 1 {
		t.Errorf("size after rejected duplicate: got %d, want 1", tr.Size())
	}
	if got := tr.Count(3); got != 1 {
		t.Errorf("Count after rejected duplicate: got %d, want 1", got)
	}
}

func TestRBTDuplicatesAllowed(t *testing.T) {
	t.Parallel()

	tr := NewRBT[int](WithAllowedDuplicates[int]())
	if tr.Size() != 0 {
		t.Fatalf("empty tree with duplicates option: size %d", tr.Size())
	}

	if err := tr.Insert(4); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if err := tr.Insert(4); err != nil {
		t.Fatalf("duplicate Insert with option: %v", err)
	}

	if tr.Size() != 1 {
		t.Errorf("size after duplicate insert: got %d, want 1", tr.Size())
	}
	if got := tr.Count(4); got != 2 {
		t.Errorf("Count after duplicate insert: got %d, want 2", got)
	}
	if got := tr.Root().Count(); got != 2 {
		t.Errorf("root Count after duplicate insert: got %d, want 2", got)
	}

	if err := tr.Delete(4); err != nil {
		t.Fatalf("Delete one of two duplicates: %v", err)
	}
	if tr.Size() != 1 {
		t.Errorf("size after decrementing count: got %d, want 1", tr.Size())
	}
	if got := tr.Count(4); got != 1 {
		t.Errorf("Count after decrementing: got %d, want 1", got)
	}
	if tr.Root() == nil {
		t.Fatal("root removed while count still 1")
	}

	if err := tr.Delete(4); err != nil {
		t.Fatalf("Delete last occurrence: %v", err)
	}
	if tr.Size() != 0 || tr.Root() != nil {
		t.Errorf("after last duplicate delete: size=%d root=%#v", tr.Size(), tr.Root())
	}
}

func TestRBTDeleteMissingValue(t *testing.T) {
	t.Parallel()

	tr := NewRBT[int]()
	if err := tr.Insert(1); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	err := tr.Delete(2)
	if err == nil {
		t.Fatal("Delete missing value: got nil error")
	}
	if err.Error() != "value not found" {
		t.Errorf("Delete missing value: got %q, want %q", err.Error(), "value not found")
	}
	if tr.Size() != 1 {
		t.Errorf("size after failed delete: got %d, want 1", tr.Size())
	}
	if got := tr.Count(1); got != 1 {
		t.Errorf("Count(1) after failed delete: got %d, want 1", got)
	}
}

func TestRBTNilAndSentinelNodes(t *testing.T) {
	t.Parallel()

	type method struct {
		name string
		fn   func(*RBTNode[int])
	}
	methods := []method{
		{name: "Value", fn: func(n *RBTNode[int]) { n.Value() }},
		{name: "Count", fn: func(n *RBTNode[int]) { n.Count() }},
		{name: "Parent", fn: func(n *RBTNode[int]) { n.Parent() }},
		{name: "Left", fn: func(n *RBTNode[int]) { n.Left() }},
		{name: "Right", fn: func(n *RBTNode[int]) { n.Right() }},
	}

	nodes := []struct {
		name string
		node *RBTNode[int]
	}{
		{name: "nil", node: nil},
		{name: "sentinel", node: sentinel[int]()},
	}

	for _, n := range nodes {
		for _, m := range methods {
			t.Run(n.name+"/"+m.name, func(t *testing.T) {
				t.Parallel()
				mustPanic(t, func() { m.fn(n.node) })
			})
		}
	}
}
