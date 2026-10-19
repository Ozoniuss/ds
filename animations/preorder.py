from manim import *

KEYS = [10, 4, 14, 2, 7, 12, 18, 1, 3, 6, 8, 5]


class Node:
    def __init__(self, key):
        self.key = key
        self.left = None
        self.right = None


def insert(root, key):
    if root is None:
        return Node(key)
    if key < root.key:
        root.left = insert(root.left, key)
    else:
        root.right = insert(root.right, key)
    return root


def preorder(node, lo=None, hi=None):
    if node is None:
        return []
    return (
        [(node.key, lo, hi)]
        + preorder(node.left, lo, node.key)
        + preorder(node.right, node.key, hi)
    )


def bound_marker(graph, name, value, default_x, color):
    x = default_x if value is None else graph.vertices[value].get_center()[0]
    top = graph.get_top()[1] + 0.6
    bottom = graph.get_bottom()[1] - 0.4
    line = Line([x, bottom, 0], [x, top, 0], color=color, z_index=-1)
    text = "nil" if value is None else str(value)
    label = Text(f"{name} = {text}", font_size=28, color=color)
    label.next_to(line, UP, buff=0.15)
    return VGroup(line, label)


def bounds(graph, lo, hi):
    left = graph.get_left()[0] - 0.8
    right = graph.get_right()[0] + 0.8
    return VGroup(
        bound_marker(graph, "lo", lo, left, BLUE),
        bound_marker(graph, "hi", hi, right, RED),
    )


def edges(node):
    if node is None:
        return []
    out = []
    for child in (node.left, node.right):
        if child is not None:
            out.append((node.key, child.key))
            out += edges(child)
    return out


def positions(node, depth=0, ranks=None, out=None):
    if out is None:
        ranks, out = [0], {}
    if node is None:
        return out
    positions(node.left, depth + 1, ranks, out)
    out[node.key] = [ranks[0], -1.1 * depth, 0]
    ranks[0] += 1
    positions(node.right, depth + 1, ranks, out)
    return out


def build(keys):
    root = None
    for key in keys:
        root = insert(root, key)
    return root


def violates(key, lo, hi):
    return (lo is not None and key <= lo) or (hi is not None and key >= hi)


class BoundsScene(Scene):
    layout_by_value = False

    def build_root(self):
        return build(KEYS)

    def construct(self):
        root = self.build_root()
        keys = [key for key, _, _ in preorder(root)]

        layout = {key: np.array(pos) for key, pos in positions(root).items()}
        if self.layout_by_value:
            for key, pos in layout.items():
                pos[0] = key
        graph = Graph(
            keys,
            edges(root),
            layout=layout,
            layout_scale=1,
            labels={key: Text(str(key), font_size=34, color=BLACK) for key in keys},
        )
        graph.scale_to_fit_width(10.5)
        graph.move_to(DOWN * 0.3)
        self.add(graph)
        self.wait(0.5)

        markers = None
        for key, lo, hi in preorder(root):
            vertex = graph.vertices[key]
            new_markers = bounds(graph, lo, hi)
            if markers is None:
                markers = new_markers
                marker_anim = FadeIn(markers)
            else:
                marker_anim = Transform(markers, new_markers)
            failed = violates(key, lo, hi)
            self.play(
                vertex.animate.set_fill(RED if failed else YELLOW, family=False),
                marker_anim,
                run_time=1.2,
            )
            if failed:
                break
            self.wait(0.4)

        self.wait(1)


class PreorderTraversal(BoundsScene):
    pass


class PreorderTraversalFails(BoundsScene):
    layout_by_value = True

    def build_root(self):
        root = build(KEYS)
        root.left.right.left.left.key = 9
        return root
