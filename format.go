package tree

import (
	"cmp"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
)

const (
	_PREFIX_LEFT  = "├──"
	_PREFIX_RIGHT = "└──"

	_EXTRA_LEFT  = "│   "
	_EXTRA_RIGHT = "    "
)

const (
	/*
	   FormatLinuxTree formats the tree nicely as a string.

	   	    4
	   	   / \
	   	  2   12
	   	 /   / \
	   	1   8   13
	   	   / \
	   	  6   9
	   	 /     \
	   	5       11

	   would get converted to

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
	   		└── 13
	*/
	FormatLinuxTree = "FormatLinuxTree"
	/*
	   FormatHorizontal formats the tree horizontally. This is the most common
	   format used to represent trees in text files.

	      4
	     / \
	    1   8
	       / \
	      6   9
	     /     \
	    5       11
	             \
	              13
	*/
	FormatHorizontal = "FormatHorizontal"
	/*
	   FormatHorizontalSquared formats the tree horizontally but using squared
	   branches. This is useful if you want the output to be wider rather than
	   longer, for example if using nodes with long labels.

	                         5sfhskfuceskjvsdnkvjkdsn
	               ┌────────────────────┬────────────────────┐
	    1dbfalkfbdslkjfbadslkfbl                 8dsbflkjsdbfjzhklsdbfljkds
	               └┐                          ┌─────────────┬─────────────┐
	     3dbfalkfbdslkjfbadslkfbl  7dsbflkjsdbfjzhklsdbfljkds  9dsbflkjsdbfjzhklsdbfljkds
	*/
	FormatHorizontalSquared = "FormatHorizontalSquared"
)

var availableFormats = []string{FormatLinuxTree, FormatHorizontal, FormatHorizontalSquared}

// FormatTree will return a string representation of the tree, based on the
// format options provided.
func FormatTree[T cmp.Ordered](t Tree[T], formatType string) string {
	if t == nil {
		return "nil tree"
	}

	if isNilOrSentinel(t.Root()) {
		return "empty tree"
	}

	if !slices.Contains(availableFormats, formatType) {
		formatType = FormatHorizontal
	}
	switch formatType {
	case FormatLinuxTree:
		return formatLinuxTree(t)
	case FormatHorizontal:
		b := strings.Builder{}
		hf := newhf[T](&b, 2, false)
		hf.formatTree(t.Root())
		return b.String()
	case FormatHorizontalSquared:
		b := strings.Builder{}
		hf := newhf[T](&b, 2, true)
		hf.formatTree(t.Root())
		return b.String()
	}
	return ""
}

func formatLinuxTree[T cmp.Ordered](t coloredTree[T]) string {
	if isNilOrSentinel(t.Root().Left()) && isNilOrSentinel(t.Root().Right()) {
		return fmt.Sprint(t.Root().Value())
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("%v\n", t.Root().Value()))
	prefix := []string{}

	type stkobj struct {
		n   Node[T]
		cnt int
	}

	stack := []*stkobj{}
	stack = append(stack, &stkobj{
		n:   t.Root(),
		cnt: 0,
	})

	for len(stack) != 0 {
		cobj := stack[len(stack)-1]
		n := cobj.n

		// processed both left and right
		if cobj.cnt >= 2 {
			stack = stack[:len(stack)-1]
			if len(prefix) != 0 {
				prefix = prefix[:len(prefix)-1]
			}
		}

		if isNilOrSentinel(n) || (isNilOrSentinel(n.Left()) && isNilOrSentinel(n.Right())) {
			stack = stack[:len(stack)-1]
			if len(prefix) != 0 {
				prefix = prefix[:len(prefix)-1]
			}
			continue
		}

		if cobj.cnt == 0 {
			l := n.Left()
			var toprint string
			if isNilOrSentinel(l) {
				toprint = "*"
			} else {
				toprint = getTtyColoredValue(l)
			}
			out.WriteString(fmt.Sprintf("%s%s %v\n", strings.Join(prefix, ""), _PREFIX_LEFT, toprint))
			prefix = append(prefix, _EXTRA_LEFT)
			stack = append(stack, &stkobj{
				n:   l,
				cnt: 0,
			})
			cobj.cnt += 1
			continue
		} else if cobj.cnt == 1 {
			r := n.Right()
			var toprint string
			if isNilOrSentinel(r) {
				toprint = "*"
			} else {
				toprint = getTtyColoredValue(r)
			}
			out.WriteString(fmt.Sprintf("%s%s %v\n", strings.Join(prefix, ""), _PREFIX_RIGHT, toprint))
			prefix = append(prefix, _EXTRA_RIGHT)
			stack = append(stack, &stkobj{
				n:   n.Right(),
				cnt: 0,
			})
			cobj.cnt += 1
			continue
		}
	}

	return strings.TrimSuffix(out.String(), "\n")
}

// horizontalFomrmatter renders a horizontal ASCII representation of a binary tree.
type horizontalFomrmatter[T cmp.Ordered] struct {
	out io.Writer
	// squareBranches prints branches using Unicode box‑drawing characters
	// instead of classic / and \.
	squareBranches bool
	// hspace is the minimum number of spaces between adjacent node labels in a
	// single tree. Must be positive. Default is 2.
	hspace int
}

func newhf[T cmp.Ordered](out io.Writer, hspace int, squarebranches bool) *horizontalFomrmatter[T] {
	p := &horizontalFomrmatter[T]{
		out:            out,
		hspace:         hspace,
		squareBranches: squarebranches,
	}
	return p
}

// formatTree renders a single tree rooted at root.
func (p *horizontalFomrmatter[T]) formatTree(root Node[T]) {
	p.writeTreeLines(p.buildTreeLines(root))
}

func (p *horizontalFomrmatter[T]) buildTreeLines(root Node[T]) []treeLine {

	if root == nil {
		return nil
	}

	if s, ok := (root).(nodeWithSentinel[T]); ok {
		if s.isSentinel() {
			return nil
		}
	}

	rootLabel := fmt.Sprint(root.Value())
	var color string
	if rootc, ok := root.(coloredNode[T]); ok {
		color = rootc.ttycolor()
	}
	if color == _COLOR_RED {
		rootLabel = ttyRed + fmt.Sprint(root.Value()) + ttyColorReset
	}
	leftLines := p.buildTreeLines(root.Left())
	rightLines := p.buildTreeLines(root.Right())

	// fmt.Println("treeliunbes", "left", leftLines, "right", rightLines)
	leftCnt := len(leftLines)
	rightCnt := len(rightLines)
	minCnt := min(leftCnt, rightCnt)
	maxCnt := max(leftCnt, rightCnt)

	// The left and right subtree print representations have jagged edges, and we essentially we have to
	// figure out how close together we can bring the left and right roots so that the edges just meet on
	// some line.  Then we add hspace, and round up to next odd number.
	maxRootSpacing := 0
	// we don't really care about lines that are only on one side because there's
	// no line on the other side to intersect with.
	for i := range minCnt {
		overlap := leftLines[i].rightOffset - rightLines[i].leftOffset
		if overlap > maxRootSpacing {
			maxRootSpacing = overlap
		}
	}
	minSpacingBetweenChildren := maxRootSpacing + p.hspace
	if minSpacingBetweenChildren%2 == 0 {
		minSpacingBetweenChildren++
	}
	rendered := stripANSI(rootLabel)
	renderedLen := len(rendered)

	// build lines again, including ones that were previously generated, with
	// new offsets.
	var allTreeLines []treeLine
	// on this current line, assume node label center is the center is the zero
	// offset.
	allTreeLines = append(allTreeLines, treeLine{
		line:        rootLabel,
		leftOffset:  -(renderedLen - 1) / 2,
		rightOffset: renderedLen / 2,
	})

	var leftAdj, rightAdj int
	switch {
	// No children. Done.
	case len(leftLines) == 0 && len(rightLines) == 0:

	// if there are lines only on one side, we don't have to be careful
	case len(leftLines) == 0:
		if p.squareBranches {
			allTreeLines = append(allTreeLines, treeLine{line: "└┐", leftOffset: 0, rightOffset: 1})
			rightAdj = 1
		} else {
			allTreeLines = append(allTreeLines, treeLine{line: "\\", leftOffset: 1, rightOffset: 1})
			rightAdj = 2
		}
	case len(rightLines) == 0:
		if p.squareBranches {
			allTreeLines = append(allTreeLines, treeLine{line: "┌┘", leftOffset: -1, rightOffset: 0})
			leftAdj = -1
		} else {
			allTreeLines = append(allTreeLines, treeLine{line: "/", leftOffset: -1, rightOffset: -1})
			leftAdj = -2
		}

	// if both children are present, you need both "branches". The distance between
	// branches will vary based on the node label length and the distance required
	// to fix all child line overlaps.
	//
	// code below just adds the line required to draw the branches. child adjustments
	// are done later
	default:
		if p.squareBranches {
			adjust := minSpacingBetweenChildren/2 + 1
			horizontal := strings.Repeat("─", minSpacingBetweenChildren/2)
			branch := "┌" + horizontal + "┬" + horizontal + "┐"
			allTreeLines = append(allTreeLines, treeLine{line: branch, leftOffset: -adjust, rightOffset: adjust})
			leftAdj = -adjust
			rightAdj = adjust
		} else {
			if minSpacingBetweenChildren == 1 {
				allTreeLines = append(allTreeLines, treeLine{line: "/ \\", leftOffset: -1, rightOffset: 1})
				leftAdj = -2
				rightAdj = 2
			} else {
				for i := 1; i < minSpacingBetweenChildren; i += 2 {
					allTreeLines = append(allTreeLines, treeLine{
						line:        "/" + spaces(i) + "\\",
						leftOffset:  -((i + 1) / 2),
						rightOffset: (i + 1) / 2,
					})
				}
				leftAdj = -(minSpacingBetweenChildren/2 + 1)
				rightAdj = minSpacingBetweenChildren/2 + 1
			}
		}
	}

	// adjust offsets for all lines that were already generated. also, merge
	// the two lines on the same row together
	for i := range maxCnt {

		// merge not required, just align to parent

		if i >= len(leftLines) {
			rl := rightLines[i]
			rl.leftOffset += rightAdj
			rl.rightOffset += rightAdj
			allTreeLines = append(allTreeLines, rl)
			continue
		}
		if i >= len(rightLines) {
			ll := leftLines[i]
			ll.leftOffset += leftAdj
			ll.rightOffset += leftAdj
			allTreeLines = append(allTreeLines, ll)
			continue
		}

		// merge required, then align to parent

		ll := leftLines[i]
		rl := rightLines[i]
		adjSpacing := minSpacingBetweenChildren
		if minSpacingBetweenChildren == 1 {
			if p.squareBranches {
				adjSpacing = 1
			} else {
				adjSpacing = 3
			}
		}
		line := ll.line + spaces(adjSpacing-ll.rightOffset+rl.leftOffset) + rl.line
		allTreeLines = append(allTreeLines, treeLine{
			line:        line,
			leftOffset:  ll.leftOffset + leftAdj,
			rightOffset: rl.rightOffset + rightAdj,
		})
	}
	return allTreeLines
}

func (p *horizontalFomrmatter[T]) writeTreeLines(lines []treeLine) {
	// fmt.Println("lines", lines)
	if len(lines) == 0 {
		return
	}
	minLeft := minLeftOffset(lines)
	maxRight := maxRightOffset(lines)
	for _, tl := range lines {
		left := -minLeft + tl.leftOffset
		right := maxRight - tl.rightOffset
		io.WriteString(p.out, spaces(left)+tl.line+spaces(right)+"\n")
	}
}

// treeLine represents a printable row from the tree. offsets are the relative
// positions of the starting and ending characters relative to the center of
// the root label (that is computed separately). A list of treeLines is required
// to print the whole tree, and the list basically contains consecutive lines
// of the output.
//
// The smallest offset represents the start of a line in the output file (file
// is a generic name here, it includes stdout for example). The largest offset
// will basically represent the longest line.
//
// Examples:
//
// (line:"┌┘   └┐" l:-3 r:3) (branches)
// (line:"/     \" l:-3 r:3)
// (line:"5" l:-1 r:-1) (single node)
type treeLine struct {
	line        string
	leftOffset  int
	rightOffset int
}

func (t treeLine) String() string {
	return fmt.Sprintf("(line:\"%s\" l:%d r:%d)", t.line, t.leftOffset, t.rightOffset)
}

// TODO: this doesn't match original sequence, verify it is correct.
func stripANSI(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[^0-9;]`)
	return re.ReplaceAllString(s, "")
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}

func minLeftOffset(lines []treeLine) int {
	min := 0
	for i, tl := range lines {
		if i == 0 || tl.leftOffset < min {
			min = tl.leftOffset
		}
	}
	return min
}

func maxRightOffset(lines []treeLine) int {
	maxVal := 0
	for i, tl := range lines {
		if i == 0 || tl.rightOffset > maxVal {
			maxVal = tl.rightOffset
		}
	}
	return maxVal
}

// coloredTree is an internal interface extending the tree interface
// to allow printing colored nodes, e.g. in Red Black trees.
type coloredTree[T cmp.Ordered] interface {
	Tree[T]
}

// coloredNode is an interface extending the node interface to allow printing
// colored nodes, e.g. in Red Black trees.
type coloredNode[T cmp.Ordered] interface {
	ttycolor() string
}

func getTtyColoredValue[T cmp.Ordered](n Node[T]) string {
	label := fmt.Sprint(n.Value())
	var color string
	if c, ok := n.(coloredNode[T]); ok {
		color = c.ttycolor()
	}
	if color == _COLOR_RED {
		label = ttyRed + fmt.Sprint(n.Value()) + ttyColorReset
	}

	return label
}

type nodeWithSentinel[T cmp.Ordered] interface {
	isSentinel() bool
}

func isNodeSentinel[T cmp.Ordered](n Node[T]) bool {
	if n == nil {
		return false
	}
	if s, ok := (n).(nodeWithSentinel[T]); ok {
		if s.isSentinel() {
			return true
		}
	}
	return false
}

func isNilOrSentinel[T cmp.Ordered](n Node[T]) bool {
	if n == nil {
		return true
	}
	if isNodeSentinel(n) {
		return true
	}
	return false
}

var ttyColorReset = "\033[0m"
var ttyRed = "\033[31m"
var ttyGreen = "\033[32m"
var ttyYellow = "\033[33m"
var ttyBlue = "\033[34m"
var ttyPurple = "\033[35m"
var ttyCyan = "\033[36m"
var ttyGray = "\033[37m"
var ttyWhite = "\033[97m"
