# Visual tests

This applies to any test in this package that describes trees via the
compact string round trip used to draw them for tests. Look at
rbt_visual_test.go for the current examples of this pattern before adding
to it.

## Naming

- All case names within one test must follow a single, consistent template
  built from the same axes — check the existing cases in the test you're
  extending and match their template. Don't mix a "position" phrasing in
  one case with a "shape" phrasing in another case of the same test.
- Fields holding a drawing should be named after the operation and node
  that produce it, not generic names like `before`/`after`.

## What visual tests cover

- Visual tests demonstrate what an operation does to the shape and colors of a
  tree. They are not for testing panics or other error paths: every case must
  start from a valid tree and exercise the operation on valid inputs. Guard
  clauses (nil checks, sentinel-destination checks, and similar) belong in other
  tests.
- The tree an operation produces does not have to be a valid BST or red-black
  tree. Low-level helpers such as rotations and transplant may leave the tree
  invalid, which is fine since constraint validation is tested elsewhere.

## Drawings

- Every drawing must make the operation and the node(s) involved
  unambiguous. Identify the relevant nodes through the case's own fields,
  and make sure each before/after pair actually demonstrates the operation
  rather than restating the same shape.
- Each case includes its own input tree as a field, even when several cases
  use the same input. Don't move it out to one shared value: the case should
  be readable on its own, without having to look somewhere else to see what
  the operation started from.
- Every drawing starts on a new line (a newline right after the opening
  backtick) so the first line lines up with the rest. Remove that leading
  newline before using the string.
- Never hand-type a drawing. Generate it from the real formatter: build the
  nodes in a throwaway test in this package, run it, capture the exact
  output, then delete the throwaway code before finishing. Hand-picked
  cases only prove the shapes you thought of, so when changing the
  formatting or parsing logic itself, also check it against many random
  trees — either write a fuzz test or reuse the random-permutation approach
  already present elsewhere in this package (look for `rand/v2`'s `PCG`
  used with a fixed seed) — and explicitly against edge cases (the empty
  tree, a single node, a node with no siblings, the smallest/largest value
  in the tree), since those are easy to miss with random generation alone.

## Scope of changes

- Only implement the test cases actually requested. Proposing additional
  cases for coverage is welcome, but don't add them unprompted.
- Before proposing a new case, check whether it exercises a genuinely
  uncovered code path or combination. If the coverage already exists, say
  so instead of adding a case "just in case".
- If a request would require changing shared scaffolding (the case struct,
  the runner, or the shape/values of an existing case) rather than just
  appending a new case, describe the change you think is needed and ask
  before making it — don't restructure existing cases or scaffolding on
  your own judgment.
