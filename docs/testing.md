# Constraint-based testing

> Before you read this, I want to emphasize that this article isn't about the theory behind Red-Black trees. No need to leave if that scares you;)

For a while I've been working on a Red-Black tree library, and I kind of wanted to talk about my journey of coming up with a test suite for it. I found this part interesting, challenging and quite annoying all at the same time, as it got me stuck in so many trains of thought wondering if I have the right approach. Of all things, when you're implementing an efficient data structure module, arguably what you care about most is correctness. I mean, ergonomics and performance don't really mean much when the code is not working as expected. In fact, most of my effort so far went into the testing suite, since the algorithms are a [textbook](https://www.cs.mcgill.ca/~akroit/math/compsci/Cormen%20Introduction%20to%20Algorithms.pdf) rip-off (though you can find them pretty much everywhere online).

When you want correctness, it means you want to write _a lot_ of tests. And that came with problems.

## Writing down a tree

Let's start with how you'd typically use a binary tree. I'd say most cases are summarized by:

```go

func main() {

	// define a tree
	t := rbtree.NewRBT[int]()

	// get values from somewhere
	t.Insert(2)
	t.Insert(6)
	t.Insert(13)
	t.Insert(1)
	t.Insert(4)
	t.Insert(19)

	// work with the tree, i.e. go through nodes
	explore(t.Root())

	// maybe more writes, more work
}
```

Well, at least kinda. This is obviously oversimplified and your code most likely never looked like this. What I wanted to show here is the lifecycle of the object. 

However, there's actually a good chance you used a tree-like data structure in your code, just not directly. In "real world" use cases, they're a common tool powering other abstractions. The [std::map](https://en.cppreference.com/cpp/container/map) and [std::set](https://en.cppreference.com/cpp/container/set) C++ containers are implemented as Red-Black trees and I would argue we have lots of production code out there with sets and maps. Those have a different API, but the operations happening behind these abstractions use a "tree-like" API. Have a look at C++'s [std::set](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_set.h) implementation, which uses this [Red-Black tree](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_set.h#L137) to see for yourself. Of course, that's assuming you can even read that shit. I kept looking at the set file for a while and wondering, "where's the fucking set, though?"

> To be "that guy", it's actually `gcc`'s implementation and it's actually [this](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_tree.h) Red-Black tree. Here's its [insertion code](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_tree.h#L764). That's a total of 307 __underscores!

As a side note, I personally have use cases where I use trees directly, which is when I do coding problems and [trees](https://github.com/Ozoniuss/Algorithms/commit/5d8d7f809c220913afab2c51a5feb266cd53e2ad) happen to show up. Then my code can end up looking similar to the snippet, with most effort going into the "work with the tree, i.e. go through nodes" part. Yeah, I find a lot of joy in figuring out if I want to go left or right. What I don't find a lot of joy in (at least with leetcode problems) is the "get values from somewhere" part. Getting those values from their inputs to test my stuff locally can be at least as complicated as writing the actual solution to the problem.

Most of the time, your input looks something like this:

![alt text](image-3.png)

Here, I took one input from [problem 2265](https://leetcode.com/problems/count-nodes-equal-to-average-of-subtree/). This is perfectly fine, as long as you represent it with an image. That's not the case in your source code, though. Instead, it ends up looking something like this:

```go
func() Tree[int] {
	root := &BSTNode[int]{value: 4}
	l := &BSTNode[int]{value: 8, parent: root}
	r := &BSTNode[int]{value: 5, parent: root}
	ll := &BSTNode[int]{value: 0, parent: l}
	lr := &BSTNode[int]{value: 1, parent: l}
	rr := &BSTNode[int]{value: 6, parent: r}
	root.left = l
	root.right = r
	l.left = ll
	l.right = lr
	r.right = rr
	return &BST[int]{root: root, size: 6}
}()
```

This is only 6 nodes, by the way.

I think at this point you're starting to see the direction I'm going in. Even for relatively small trees, this representation totally lacks clarity. Sure, I can see arguments coming from a mile away that this isn't painful anymore now with AI, but still, there's no way you can give that snippet to someone and have them understand what the tree looks like. You can imagine why for tests, it becomes quite a significant downside. Or, you could of course argue that reading code (particularly tests) isn't a problem anymore either. Whatever floats your boat.

Admittedly, there are ways to simplify writing these inputs. One example is writing it as a list, which is how LeetCode does it:

```
root = [4,8,5,0,1,null,6]
```

> To give them some credit, you can now edit the inputs and see the tree with a visualizer:

![alt text](image-5.png)

This representation is much more compact, but it still has the same downside of lacking clarity. I was actually surprised the first time I saw how `[4,null,5,null,6]` looks: 

![alt text](image-6.png)

Another is what I did in the beginning: create the tree with a sequence of inserts. That's essentially the same as the representation above, since that list is a sequence of inserts. LeetCode has `null` elements because some of their test inputs aren't valid binary search trees, like the first image.

I suppose eventually, you'd get used to whatever representation you use most commonly. But there's no denying that images are simply a much, much better representation of a binary tree input. I actually ended up [writing a renderer](https://github.com/Ozoniuss/ds/blob/8ca58b28055dfc3b02b78dea8a271187d7cab0f8/rbtree/format.go#L73) just so I can see tree drawings in my source code! Here's a snippet from my [test suite](https://github.com/Ozoniuss/ds/blob/8ca58b28055dfc3b02b78dea8a271187d7cab0f8/rbtree/rbt_visual_test.go):

>  Here I'm testing transplant(u, v) which overwrites node u with node v, using table-driven testing. Note that transplant can leave the tree in an invalid state. 

```go
func TestTransplant(t *testing.T) {

	type testcase struct {
		// testcase definition
	}

	cases := []testcase{
		// ...
		{
			name: "u is the right child of its parent and v is an internal subtree",
			u:    30,
			v:    5,
			initialTree: `
       20[B]
     v---+---v
   10[B]   30[B]
 v---+---v
5[R]   15[R]`,
			transplantResult: `
      20[B]
     v--+--v
   10[B]  5[R]
 v---+---v
5[R]   15[R]`,
		},
		// ...	
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// check initial tree looks as expected after transplant
		})
	}
}
```

I'll be covering how it works in a future article, since it was tricky to get right. I find this a significant step up in terms of clarity but still, wasn't good enough for me to be satisfied with proving correctness. On top of it, you actually need to, well, write these inputs:) No wonder AI generates all your test cases now. When you start to think a bit more about it, I'm basically paying a bunch of big-ass corporate monopolies to write test cases for my open source library. I prefer not to think about it.

## Testing for correctness

Let's dive into correctness a bit more. So far, I kind of implied that your testing approach looks something like this:

- Start with a set of inputs;
- Apply the operation on each of those inputs;
- Evaluate the result of each operation against an expected result.

This is very basic but also a great description of how most unit test cases are set up. Let's be honest, you probably just want to pass that coverage requirement, don't you? While I find it suitable for reasonably simple logic, it becomes tricky for Red-Black tree operations. To put it into perspective, here's what the insert implementation [looks like](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree#Insertion). This is without the underscores! The book that I've been reading dedicates 6 pages just to explain the rebalancing procedure you may or may not need to execute at the end of the insertion. Here's a snippet. ![alt text](image-4.png)

Sure, having to read through this isn't the worst thing in the world (remember those C++ snippets from earlier?). Nonetheless, it's complex enough to wrap your head around and inevitably doesn't give me enough confidence that I have a sufficient number of test cases to cover all branches (pun unintended). In theory, those are finite (textbook describes 3 cases) but not in the sense that each test input will pass through one of these 3 cases. Rather, the intermediary state during an insertion can be one of these 3 cases. Rebalancing is done in a loop, so there are insertion operations that can fall through all cases, some multiple times. Even if you re-write the code so that you can individually test every rebalancing case, it still leaves the open question of what inputs are suitable for complete insertion tests.

There are some additional considerations here [1,2], but perhaps the one worth pointing out separately is that, while the implementation of the insertion operation is (or at least should definitely be) deterministic, there are multiple valid rebalancing implementations, with different efficiency considerations. The same input can produce different valid Red-Black trees under different rebalancing algorithms, meaning that a test suite like this isn't portable.

At this point, I suggest we take a step back and think about these tests from a different perspective. Let's briefly mention what a Red-Black tree actually is. All you need to become one is to satisfy a set of invariants (constraints):

0. It must satisfy the binary search tree property
1. Every node is colored either red or black.
2. The root is black.
3. Every leaf (TNIL) is black.
4. If a node is red, then both its children are black
5. For each node, all simple paths from the node to TNIL leaves contain the same number of black nodes.

If that's you, congrats! You're a Red-Black tree now 🌞

Like I promised, I won't go into details about Red-Black trees. The important bit is the result they provide: _the tree is approximately balanced_. This effectively allows for implementations of write and read methods that run in logarithmic time with respect to the number of nodes in the tree. On a higher level, it offers a means to implement a collection of sorted elements efficiently.

The nice thing about constraints is that they must hold after every write on a Red-Black tree. This is a substantial benefit when designing a test suite, because _we no longer care what the tree looks like_. [3] In the beginning, I showed how you generally use a tree:

```go
func main() {

	// define a tree
	t := rbtree.NewRBT[int]()

	// do some writes on the tree with values coming from somewhere
	t.Insert(2)
	t.Insert(6)
	t.Insert(13)
	t.Insert(1)
	t.Insert(4)
	t.Insert(19)

	// ...
}
```

Focusing on constraints enables us to employ a similar pattern in our tests. After every operation, we just assert that the resulting structure is a Red-Black tree, without having to visualize its shape. We could come up with a testing sequence that looks something like this:

- Generate a random list of (unique) elements;
- Insert every element into the tree;
- Check the constraints hold after every insertion;
- Remove every element from the tree, maybe in a different order;
- Check the constraints hold after every removal.

Turns out, you can validate the constraints in a single tree traversal. Below, you can find the algorithm that I'm using:

```go
func assertRBTProperties[T cmp.Ordered](t *testing.T, tree *RBT[T]) {
	t.Helper()

	if tree == nil {
		panic("assert called on nil tree")
	}

	// nothing to assert here
	if tree.Size() == 0 {
		return
	}

	if tree.Root() == nil {
		panic("nil root")
	}

	if tree.Root().color != _COLOR_BLACK {
		t.Fatal("property 2 violation: root should be black")
	}

	var walk func(n *RBTNode[T], bnCount int, lo, hi *T)
	bnCountTarget := -1 // unset
	walk = func(n *RBTNode[T], bncount int, lo, hi *T) {

		if n.color != _COLOR_RED && n.color != _COLOR_BLACK {
			t.Fatal("property 1 violation: found a node with an invalid color")
		}

		if n.isSentinel() {
			if n.color != _COLOR_BLACK {
				t.Fatal("property 3 violation: sentinel node should be black")
			}
			// if we haven't explored any simple path to TNIL yet, set it as the
			// target to compare to. Otherwise, compare it against an existing
			// target since the black node count must match.
			if bnCountTarget == -1 {
				bnCountTarget = bncount
			}
			if bnCountTarget != bncount {
				t.Fatal("property 5 violation: found two paths with different black length")
			}
			return
		}

		v := n.Value()
		if lo != nil && v <= *lo {
			t.Errorf("bst violation: %v is in the right subtree of %v", v, *lo)
		}
		if hi != nil && v >= *hi {
			t.Errorf("bst violation: %v is in the left subtree of %v", v, *hi)
		}

		if n.color == _COLOR_RED &&
			(n.left.color != _COLOR_BLACK || n.right.color != _COLOR_BLACK) {
			t.Fatal("property 4: found red node whose children aren't both black")
		}

		if n.color == _COLOR_BLACK {
			bncount += 1
		}

		walk(n.left, bncount, lo, &v)
		walk(n.right, bncount, &v, hi)
	}

	walk(tree.Root(), 0, nil, nil)
}
```

And also, one of my tests:

```go
func TestRBTPropertiesDeterministicPermutations(t *testing.T) {
	rng := rand.New(rand.NewPCG(6, 7))

	for i := range 100 {
		n := rng.IntN(1_000) + 1
		values := rng.Perm(1_000)[:n]

		t.Run(fmt.Sprintf("%03d/n=%d", i, n), func(t *testing.T) {
			tr := NewRBT[int]()
			cnt := 0
			for _, v := range values {
				tr.Insert(v)
				cnt += 1
				assertRBTProperties(t, tr)
				if tr.Size() != cnt {
					t.Fatalf("tree size is not %d after inserts: %d", cnt, tr.Size())
				}
			}
			for _, v := range values {
				tr.Delete(v)
				cnt -= 1
				assertRBTProperties(t, tr)
				if tr.Size() != cnt {
					t.Fatalf("tree size is not %d after deletes: %d", cnt, tr.Size())
				}
			}
		})
	}
}
```

And there you go! With randomized inputs, we're obviously able to cover a lot more ground than specifying test cases by hand. Knowing that every operation produces a valid Red-Black tree gives me a high degree of confidence that the implementation is correct. Or, at least it's doing a good job in making it easier to lie to myself. It was important to me to have coverage, but I found little research about what kind of distribution randomized tree inserts provide in Red-Black trees. Admittedly, I didn't put a lot of effort into researching this and just went YOLO with `go test -cover` [4]. Anyway, for me it was good enough to prove correctness.

Or, was it?

## The meaning of correctness

At first glance, I would say the approach above is quite convincing that my tree is correct. At least, I was convinced in the beginning. Though we are overlooking a few things if we just leave it there.

Earlier, I said Red-Black trees allow for efficient implementations. I want to heavily stress the word "allow". On their own, constraints only define valid tree structures on an abstract level. Modelling that representation in code **does not automatically mean that your operations run efficiently!** To give an analogy, even if you have a list of sorted numbers, nothing stops you from searching through it sequentially. I will not discuss performance tests more in depth, but you should definitely consider adding them to your test suite.

Well, that kinda sucks, but I can't really be bummed out about it since now I have a reason to write benchmarks, which is quite fun. But at least, my tree is correct, right? Turns out, the answer isn't quite that simple. I like to think there are different ways to view "correctness". Originally, the way I looked at it was being able to say "yeah this thing is fast!" In that regard, satisfying the Red-Black tree properties does give me plenty of confidence, especially with benchmarks in place.

However, that's slightly different from proving that the internal representation of my data structure and the API I'm exposing behave correctly. The angle I use to approach this kind of correctness is through the lens of the question, "Can my program ever reach an invalid state?" There's a difference between the theoretical set of constraints aimed towards efficiency and the practical set of invariants your internal representation must hold at all times, even if they may overlap. [5] For example, my function above does not test that every node has a valid link to a parent, or that the count of internal nodes matches the number of unique elements inserted. And I do expose a `Parent()` method on the node and a `Size()` method on the tree, plus the write algorithms change quite a few links. It's definitely problematic when the internal representation no longer looks like you expect.

There are various layers to it on the implementation side. Arguably, what you care about most is the public API. You could expose it as an interface and test via the interface alone, which allows you to swap internal implementations (for example, some tree implementations don't hold parent pointers). At the same time, I find it pretty important to also validate your specific internal representation. It has the downside of coupling the tests to your internal representation which requires changes to the test suite if you switch it, but in this case it would only require changes in my validation function. It's nowhere near as bad in terms of portability as having to potentially change an entire test suite's outputs, such as if you tested explicit trees and changed the rebalancing implementation. [6]

Honestly, to me it's a matter of taste. I think if you chase the "perfect" abstractions, you may end up writing C++ template code. So I just try to find an approach that gives me enough confidence but is also reasonably cheap to refactor. One of the hardest things I had to accept when switching from maths to engineering is that things weren't so black and white anymore (I guess the're black and red). To me engineering is all about making the right trade-offs, and in this particular case it means finding the right combination of validating the abstract constraints, the additional internal ones and the public API. I find it perfectly reasonable to combine all of them in your tests, so be creative and have some fun:)


## Additional considerations

1. Coming up with valid test cases and understanding what happened is still difficult for Red-Black trees, even with visual tests. Can you easily tell if the input tree is valid and what the hell happened during the insertion?

```go
		{
			name: `the recoloring moves the violation up to a red parent whose own
			uncle is black, so the fixup loops and then rotates`,
			value: 5,
			tree: `
            50[B]
          v---+---v
        25[R]   75[B]
      v---+---v
    15[B]   35[B]
  v---+---v
10[R]   20[R]`,
			beforeFixup: `
               50[B]
             v---+---v
           25[R]   75[B]
         v---+---v
       15[B]   35[B]
     v---+---v
   10[R]   20[R]
    /
  5[R]`,
			insertResult: `
               25[B]
         v-------+-------v
       15[R]           50[R]
     v---+---v       v---+---v
   10[B]   20[B]   35[B]   75[B]
    /
  5[R]`,
		},
```

2. I'm not sure if every structurally valid Red-Black tree can be reproduced using a sequence of insertions and deletions, at least with the algorithm that I'm using. One might argue an "unreachable" input isn't a valuable test case, which I can understand, though personally I think it should be if the operation itself isn't different. The algorithm is designed based on the constraints alone, which are still satisfied. Still, it's something I worried about when writing inputs.

3. We still care about inputs, particularly how we generate them. For example, if your Red-Black tree doesn't support duplicates, you have to decide if you want to avoid them (and simplify e.g. testing the size) or include them (and e.g. test insertion errors in the case your API returns one).

4. `go test -cover` measures statement coverage. It tells you every line ran at least once, but not which combinations of rebalancing cases ran, or how many loop iterations. So it doesn't fully answer the question from earlier. In my case, both fixup functions are at 100%. The only lines the randomized test misses are inserting a duplicate and deleting a missing value (my inputs are unique), plus some panic guards.

5. If you read [The Power of 10: Rules for Developing Safety-Critical Code](https://en.wikipedia.org/wiki/The_Power_of_10:_Rules_for_Developing_Safety-Critical_Code) you will find advice to actually test invariants at runtime. I won't do that in my implementation, but I do `panic` in some cases that I know would break the code.

6. I still kept a reasonably small number of visual tests. That was mainly so I can better understand what is going on in the operations myself, but it also adds a slight (though not necessarily rational) increase in my confidence in my test suite. Plus, coding is pay to win these days anyway, so refactoring is "cheap".



What I love about constraints is that you can use them to infer other constraints. For Red-Black trees, one particularly interesting one is that a subtree rooted at a node `x` contains at least `2^bh(x) - 1` internal nodes, where `bh` is the black height of the node. In turn, this implies that the height of any Red-Black tree is at most `2*lg(n+1)`, where `n` is the number of nodes. That gives you more possibilities when deciding which sets of constraints you might want to validate in your tests, because some may have more efficient algorithms to check or may simply be easier to write. Just be careful when testing a different set of constraints than the original ones, as **they must imply the original set of constraints**. Just satisfying the height property alone **does not** imply a tree is a Red-Black tree.

> I like math a lot, so I like to think of the starting constraints as some axioms that define a space where you can play around and explore its properties. It's just amazing how many new things you can deduce starting from just a few statements that we label as true. Strictly speaking, there's more to it than just the starting 6 axioms. You could go even further and ask yourself, "What does a tree, a node or a child actually mean?" but here I find it important to be practical. As long as the underlying language used to define these "axioms" is the same for everyone (e.g. it is clear what a child node means), then the constraints are a good starting point to explore.
