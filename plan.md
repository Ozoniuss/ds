# Pool-allocate RBT nodes into a single backing slice

## Context

`rbtree/rbt.go` heap-allocates one `*RBTNode[T]` per `Insert`. At 100k inserts that is 100k
separate allocations, each 64 bytes, and the GC has to walk 3 pointers per node on every cycle.
The benchmarks in `rbt_bench_test.go` already report `ns/insert` and `inserts/sec`, so the win
is directly measurable.

The goal is a pool allocator: one growable `[]node[T]` per tree, nodes addressed by `uint32`
index instead of pointer, freed slots recycled through a free list. Go still owns the slice's
lifetime, so there is nothing to release manually.

Two decisions already made:

- **Index links, handle-based public nodes.** Storage becomes an unexported `node[T]` with
  `uint32` links. `RBTNode[T]` becomes a small `{tree, index}` handle so `Root()`, `Left()`,
  `Right()` and `Parent()` keep their exact signatures and keep returning `nil` — the public
  API compiles unchanged.
- **Shrink the node layout in the same pass**, since the struct is being rewritten anyway.
  `color string` (16 bytes plus a `memequal` on every color test in the fixup loops) becomes a
  `bool`; `sentinel bool` disappears because index 0 *is* the sentinel; `count` is dropped
  entirely (see below).

Net effect for `T = int`: 64 bytes/node → 24, one amortized allocation per doubling instead of
one per insert, and — because `node[T]` holds no pointers when `T` holds none — the GC stops
scanning the backing array entirely.

### On `count`, and why not a color bitset

`count` is written as `1` at creation and never incremented — dead since e78991f removed
duplicates. It is worth deleting rather than narrowing to `uint32`, because of how Go lays the
struct out for `T = int`:

```
value 0..8   left 8..12   right 12..16   parent 16..20   count 20..24   red 24..25
                                                          → struct align 8 → 32 bytes
```

`red` is a single byte forcing a 7-byte tail pad, so moving colors into a separate bitset would
buy 8 bytes and land at 24. But `count` costs exactly those same 8 bytes, so dropping it reaches
24 too, with `red` left as an inline bool and no bitset to maintain. A bitset turns every color
read into a second dependent load plus a shift and mask, adds a region to grow and clear in
`alloc`/`release`, and saves no cache traffic — both fixups read a node's color right alongside
its `left`/`right`/`parent`, so that line is already resident.

Consequences: `RBTNode.Count()` returns `1` and `RBT.Count(value)` returns 0 or 1. If duplicate
support is ever restored, put `count` back and steal the color into its top bit — *not* into
`parent`, which the fixup loops read constantly and where one missed mask silently corrupts a
link.

## Design

All of this is in `rbtree/rbt.go`.

### Storage

```go
const tnil uint32 = 0

type node[T cmp.Ordered] struct {
	value       T
	left, right uint32
	parent      uint32
	red         bool
}

type RBT[T cmp.Ordered] struct {
	nodes []node[T]  // index 0 is the tnil sentinel
	root  uint32
	free  uint32     // head of the free list; tnil means empty
	size  int
}
```

The `tnil *RBTNode[T]` field goes away. `lazyInit` becomes:

```go
func (t *RBT[T]) lazyInit() {
	if len(t.nodes) == 0 {
		t.nodes = append(t.nodes, node[T]{}) // sentinel: black, all links tnil
		t.root = tnil
	}
}
```

The zero-value tree still works: `root == 0 == tnil` means empty. Read-only methods that
currently compare against `t.tnil` (`Root`, `Count`) must call `lazyInit` or guard on
`len(t.nodes) == 0`, since they can now run before any write.

The sentinel's fields are still written during `Delete`/`rbtDeleteFixup` (via `rbtransplant`
setting `nodes[tnil].parent`), exactly as today — index 0 behaves like the old `t.tnil` object.

### Allocator

```go
// alloc may grow t.nodes, invalidating any *node[T] taken from it. Callers must
// hold only uint32 indices across this call.
func (t *RBT[T]) alloc(value T) uint32 {
	if t.free != tnil {
		i := t.free
		t.free = t.nodes[i].parent
		t.nodes[i] = node[T]{value: value}
		return i
	}
	t.nodes = append(t.nodes, node[T]{value: value})
	return uint32(len(t.nodes) - 1)
}

func (t *RBT[T]) release(i uint32) {
	t.nodes[i] = node[T]{parent: t.free} // free list threads through .parent
	t.free = i
}
```

The free list costs no extra memory — it reuses the `parent` field of dead nodes.

### Public handle

```go
type RBTNode[T cmp.Ordered] struct {
	t *RBT[T]
	i uint32
}

func (n *RBTNode[T]) Left() *RBTNode[T] {
	panicIfNilRBTNodeOrSentinel(n)
	l := n.t.nodes[n.i].left
	if l == tnil {
		return nil
	}
	return &RBTNode[T]{n.t, l}
}
```

`Value`, `Parent` and `Right` follow the same shape (`Count` just returns 1, per the section
above). `panicIfNilRBTNodeOrSentinel` checks
`n == nil || n.t == nil || n.i == tnil`, replacing `isSentinel()`. The free `sentinel[T]()`
constructor is only used by tests now — keep it as a test helper returning
`&RBTNode[int]{t: tr, i: tnil}`.

Two behaviour changes worth a doc comment on `RBTNode`:

- Handles are no longer comparable by pointer identity — `a.Left() == b` will not work. Add an
  `Equal` method or document comparing `Value()`.
- After `Delete`, a handle to the removed node is invalid and its slot may be reused by a later
  `Insert`, so reading it can silently return another node's value. (Today a deleted node's
  pointer stays readable but detached.) A generation counter in the node, bumped on `release`
  and checked by the handle, would make this detectable; not worth it until something needs it —
  and note it would cost the 8 bytes just reclaimed.

### Porting the algorithms

`insertFixup`, `rbtDeleteFixup`, `leftRotate`, `rightRotate`, `rbtransplant` and
`treeMinimumRbt` all keep their structure; `*RBTNode[T]` parameters become `uint32`, and
`x.color == _COLOR_BLACK` becomes `!t.nodes[x].red`.

The one genuine risk is readability and staleness of chained access. `z.parent.parent.right`
becomes unreadable if inlined, so bind locals:

```go
p := t.nodes[z].parent
g := t.nodes[p].parent
y := t.nodes[g].right
```

**Indices are stable, but the field values read into those locals are not.** Rotations rewrite
parent links, so any cached `p`/`g` must be re-read after a rotation — the current code gets
this for free by re-evaluating `z.parent` each time. Every rotation site in `insertFixup` and
`rbtDeleteFixup` needs checking against this.

A `func (t *RBT[T]) n(i uint32) *node[T] { return &t.nodes[i] }` helper is safe inside the
fixups and rotations because none of them allocate — `Insert` calls `alloc` once, before
`insertFixup`. Do not use it in `Insert` around the allocation site.

### Insert

Restructure so the search runs first and `alloc` happens only on a real insertion — the current
code allocates `z` before the loop and abandons it on a duplicate, which the arena would turn
into a leaked slot. Indices survive the `append` inside `alloc`, which is exactly why this is
safer than holding raw pointers into a growable slice:

```go
y, x := tnil, t.root
for x != tnil {
	y = x
	switch {
	case value < t.nodes[x].value: x = t.nodes[x].left
	case value > t.nodes[x].value: x = t.nodes[x].right
	default: return
	}
}
z := t.alloc(value) // may move t.nodes; y, x and t.root are indices, so unaffected
```

The existing special case for an empty tree (`y == tnil` → z becomes a black root) stays.

### Delete

Logic is unchanged. `z` is the node detached from the tree in all three branches (in the
two-child case `y` takes `z`'s place and keeps its own identity), so add a single
`t.release(z)` after `rbtDeleteFixup` and before/after `t.size--`. Releasing after the fixup is
required: the fixup reads `nodes[tnil].parent`, which `rbtransplant` set from `z.parent`.

## Files to change

- `rbtree/rbt.go` — everything above.
- `rbtree/rbt_test.go` — mechanical but the bulk of the churn (~490 lines):
  - `assertRBTProperties`' `walk` closure currently recurses over `*RBTNode[T]` reading
    `.left`/`.right`/`.color`/`.isSentinel()`. Convert it to recurse over `uint32`, starting at
    `tree.root`, reading `tree.nodes[i]`. Drop the "invalid color" check — vacuous for a bool.
  - The `leftRotate`/`rightRotate`/`rbtransplant` tests build nodes with literals like
    `&RBTNode[int]{parent: tr.tnil, left: tr.tnil, ...}`. Replace with `tr.alloc(v)` plus
    explicit link assignment; add a small test helper if the repetition is annoying.
  - `TestRBTNilAndSentinelNodes` needs `sentinel[int]()` to produce a tnil handle bound to a tree.
- `rbtree/rbt_bench_test.go` — add a delete benchmark; the insert benchmarks already carry
  `-benchmem`.
- `rbtree/README.md` — only if it describes node layout; the grep showed it discusses the
  interface design, not internals, so likely untouched.

## Verification

1. `go test ./rbtree/ -run . -count 1` — `assertRBTProperties` already validates all five
   red-black invariants plus the BST ordering and node count after every operation, so a
   correct port is well covered by the existing suite.
2. `go vet ./rbtree/` and `go test ./rbtree/ -race`.
3. Add a targeted test for slot reuse: insert N, delete all, insert N again, assert
   `len(tr.nodes)` stayed at `N+1` and the tree properties still hold. This is the one behaviour
   the current suite cannot already exercise.
4. Baseline before touching anything, then compare:
   `go test ./rbtree/ -run '^$' -bench '^Benchmark(RBT|Slice)' -benchmem -count 5 > before.txt`
   and diff with `benchstat`. Expect `allocs/op` for `BenchmarkRBTInsertRandom/100000` to drop
   from ~100k to the number of slice doublings, and `B/op` to roughly halve.

## Possible follow-ups (measure first)

- `NewRBTWithCapacity(n)` to preallocate the backing slice and skip the doubling copies.
- Reserving capacity is also what makes a `Clear()` cheap (`nodes = nodes[:1]; free = tnil`).
- Only if handle allocation shows up in a traversal profile: a value-handle API
  (`Left() (RBTNode[T], bool)`) removes the per-accessor 16-byte allocation, at the cost of
  breaking the current signatures.
