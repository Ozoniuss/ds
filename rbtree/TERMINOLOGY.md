# Terminology

Note that the terminology I've used throughout this project may not
perfectly match academic terminology. Here I've only written the one
that isn't obvious.

Binary tree: structure defined on a set of nodes that either contains no
nodes, or is composed of three disjoint sets of nodes: a root node, a
binary tree called the left subtree and a binary tree called the right
subtree. The binary tree containing no nodes is called the empty or null
tree. Note that the position of the subtree (left or right) is part of
the tree's structure.

Left child: the root of a non-empty left subtree.

Right child: the root of a non-empty right subtree.

Binary-search-tree property: 

Simple path: a path where all nodes are different and each node is
connected to adjacent nodes via an edge. The length of a simple path is
the number of edges, not the number of nodes.

Depth of a node: the length of a simple path from the root to the node.
For example, the depth of 1 is 0, and the depth of 6 is 2.

    1
   / \
  2   13
 / \   \
4   6   19

height: does not include root
simple path: all nodes are different, with length being the number of edges
binary tree: not necessarily satisfying the bst property
ancestor, descendant: a node is its own ancestor/descendant (vs proper ancestor)
there is a unique simple path from the root to any node. root is everyone's ancestor

degree of x:  number of children of x
lenght of simple path from root to x: depth
level: all nodes of the same depth
height of node: number of edges of the longest simple path from the node to a leaf
height of tree: largest depth from all nodes

black height: number of black nodes from x (not including it) to the leaf (including tnil)
height: includes tnil
