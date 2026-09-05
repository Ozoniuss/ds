package rbtree

import (
	"cmp"
	"fmt"
	"strings"
)

// BuildLines recursively builds the tree's visual representation as an array of
// strings, each element representing a line from the tree's string representation
// starting with the root.
//
// minDistBetweenSubtrees represents the minimum distance that there will be
// between a left child's non-empty characters and a right child's non-empty
// characters on every line of their representation.
//
// The two most important properties of the printing algorithm are as follows:
// 1. The way a subtree is printed will look identical regardless of what other
// nodes are in a tree.
// 2. A parent will always be written in the middle of its two children (well,
// technically on the bisector, not the actual middle) and its two children
// will always be drawn on the same line.
//
// In case we have a single child, the edge between the parent and the child
// will always have the smallest length.
//
// With these properties we can design an algorithm for drawing a subtree rooted
// at x as follows:
// - If there are no children, just write the node's label as a single line.
// - "Draw" the representation of x's left and right children as if they were
// separate trees. Note that they will look identical even they are drawn as part
// of x's representation; the only difference is going to be the alignment in
// the tree.
// - Note that the two representations are actually at the same level, based on
// property 2. Now, we "merge" the two representations such that the last character
// of the left subtree's representation is at least minDistBetweenSubtrees characters away
// the right subtree's representation, on every line. Basically, this requires
// finding an offset for shifting the right subtree's representation to the right
// such that when we "overwrite" the left subtree's representation on top of the
// shifted one of the right subtree, only empty characters will get overwritten.
// We also want to ensure that the distance between where the left root ends and
// the right root starts is odd, so that we can reliably find a middle for the
// node.
// - Once we found the middle, we can easily write the lines from the node to
// the left and right subtree.
// - In case the node only has a single child, there is no merging step and we
// just need to draw the node such that the edge connecting the node with its
// child has the smallest possible length. If we have a right subtree, this may
// require an offset for its representation (e.g. if we draw a right skewed tree,
// the root would need to be placed on the left of its child's representation
// which would result in a negative x position without adjusting the child
// representation).
func BuildLines[T cmp.Ordered](n *RBTNode[T], minDistBetweenSubtrees int, shouldCenter bool) []string {
	if n == nil || n.isSentinel() {
		panic("BuildLines: n must be an internal node")
	}

	if n.Left() == nil && n.Right() == nil {
		return []string{fmt.Sprintf("%v", n.Value())}
	}

	var leftLines, rightLines []string
	if n.Left() != nil {
		leftLines = BuildLines(n.Left(), minDistBetweenSubtrees, shouldCenter)
	}
	if n.Right() != nil {
		rightLines = BuildLines(n.Right(), minDistBetweenSubtrees, shouldCenter)
	}

	if len(leftLines) == 0 && len(rightLines) == 0 {
		panic("BuildLines: encountered empty lines for both children")
	}

	roffset := 0
	if n.Left() != nil && n.Right() != nil {
		for i := range min(len(leftLines), len(rightLines)) {

			lend := len(leftLines[i])
			rstart := strings.IndexFunc(rightLines[i], func(r rune) bool {
				return r != ' '
			})

			// essentially, we need rstart > lend + minDistBetweenSubtrees on
			// any overlapping lines.
			roffset = max(roffset, lend-rstart+minDistBetweenSubtrees)
		}

		// we need to have a valid "middle" in order to satisfy 2, so adjust offset
		// in case parent ends up with a floating point x position.
		lstart := findCenter(leftLines[0], shouldCenter)
		rstart := findCenter(rightLines[0], shouldCenter)
		fmt.Println("lines", leftLines[0], rightLines[0], lstart, rstart, roffset)
		rstart += roffset
		if (lstart+rstart)%2 == 1 {
			roffset += 1
		}
		fmt.Println("fixed", roffset, lstart, rstart+1)
	} else if n.Right() != nil {
		rstart := findCenter(rightLines[0], shouldCenter)
		// node will need to be at position 0, so right child needs to start at
		// position at least 2 + pad in order to properly draw the connection.
		roffset = max(0, 2-rstart)
	}

	// Merge the left and right lines. Because we're calculating the offset based
	// on non-emtpy characters, the left subtree may actually overwrite part of the
	// right subtree's representation, e.g. if the left and right subtree both have
	// a very long label that is far enough from the start of the line. The distance
	// between these labels may only be minDistBetweenSubtrees and thus have the left
	// tree's label overwrites empty characters of the right tree's representation,
	// even with the offset shift. To relibaly write the final representation we
	// first write the right subtree with the offset, then overwrite those lines
	// left tree's content.
	out := []string{}
	offsetstr := strings.Repeat(" ", roffset)
	// write right tree with offset
	for _, line := range rightLines {
		out = append(out, offsetstr+line)
	}
	for i, line := range leftLines {
		// no overlapping right line
		if i >= len(rightLines) {
			out = append(out, line)
		} else {

			if !(len(out[i]) > len(line)) {
				panic("BuildLines: offset should have made all right subtree's lines longer")
			}

			// overwrite right line with left line
			out[i] = line + out[i][len(line):]
		}
	}

	// first line will either have one word or two words
	firstLine := out[0]
	fmt.Println("merged", firstLine)
	var start1, end1, start2 int
	start1, end1 = findCenterAndEnd(firstLine, shouldCenter)

	// there are two words
	if n.Left() != nil && n.Right() != nil {
		start2 = end1 + findCenter(firstLine[end1:], shouldCenter)
	}

	var pos int
	ret := []string{}
	// write the parent and its connecting lines
	if n.Left() != nil && n.Right() != nil {
		fmt.Println("after merge", start1, end1, start2)
		if (start1+start2)%2 != 0 {
			panic("BuildLines: cannot find integer middle position for parent")
		}
		pos = (start1 + start2) / 2
		parentLine := strings.Repeat(" ", pos) + fmt.Sprintf("%v", n.Value())
		toplines := []string{parentLine}
		diff := 1 // between left and right connection
		// draw a connection until you reach the left and right children
		for empty := pos - 1; empty > start1; empty-- {
			line := strings.Repeat(" ", empty)
			line += "/"
			line += strings.Repeat(" ", diff)
			diff += 2
			line += "\\"
			toplines = append(toplines, line)
		}

		ret = append(toplines, out...)
	}

	if n.Left() == nil {
		pos = start1 - 2
		parentLine := strings.Repeat(" ", pos) + fmt.Sprintf("%v", n.Value())
		secondline := strings.Repeat(" ", pos+1) + "\\"
		toplines := []string{parentLine, secondline}

		ret = append(toplines, out...)
	}

	if n.Right() == nil {
		pos = start1 + 2
		parentLine := strings.Repeat(" ", pos) + fmt.Sprintf("%v", n.Value())
		secondline := strings.Repeat(" ", pos-1) + "/"
		toplines := []string{parentLine, secondline}

		ret = append(toplines, out...)
	}

	if shouldCenter {
		// strip original whitelines
		ret[0] = ret[0][pos:]
		pad := getPad(n)
		if pad > pos {
			for i := 1; i < len(ret); i++ {
				ret[i] = strings.Repeat(" ", pad-pos) + ret[i]
			}
			return ret
		} else {
			ret[0] = strings.Repeat(" ", pos-pad) + ret[0]
		}
	}
	return ret

}

// getPad returns how much you need to move the node value to the left
// when printing so it's centered.
func getPad[T cmp.Ordered](n *RBTNode[T]) int {
	return (len(fmt.Sprintf("%v", n.Value())) - 1) / 2
}

func findCenter(line string, center bool) int {
	start := strings.IndexFunc(line, func(r rune) bool {
		return r != ' '
	})
	end := strings.IndexByte(line[start:], ' ')
	if end == -1 {
		end = len(line)
	} else {
		end += start
	}
	lenght := end - start
	if center {
		start += (lenght - 1) / 2
	}
	fmt.Println("called with line", line, "found", start)
	return start
}

func findCenterAndEnd(line string, center bool) (int, int) {
	start := strings.IndexFunc(line, func(r rune) bool {
		return r != ' '
	})
	end := strings.IndexByte(line[start:], ' ')
	if end == -1 {
		end = len(line)
	} else {
		end += start
	}
	lenght := end - start
	if center {
		start += (lenght - 1) / 2
	}
	return start, end
}
