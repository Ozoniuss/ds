# Constraint-based testing

Before you read this, I want to emphasize that this article isn't some practical guide for how you should write library tests, or testing tips in general. In fact, knowledge of basic testing practices is a prerequisite. It's rather about how I went around some challenges when writing my Red-Black tree library. If you came here for a tutorial, sorry but this isn't for you.

This article also isn't about the theory and implementation behind Red-Black trees. No need to leave if that scares you;) Ok, let's get into it then.

For a while I've been working on a Red-Black tree library, and I kind of wanted to talk about one part of building it that I found annoying, interesting and challenging at the same time. Of all things, when you're looking an efficient data structure module, arguably what you care about most is correctness. I mean, ergonomics and performance don't really mean much when the code is not working as expected. A good part of my development effort went into the testing suite. Actualy, most of it really, since the algorithms are a [textbook](https://www.cs.mcgill.ca/~akroit/math/compsci/Cormen%20Introduction%20to%20Algorithms.pdf) rip-off, though you can find them pretty much everywhere.

Since you want correctness, it means you want _a lot_ of tests. And that came with a problem.

## Writing down a tree

Let's start with how you'd typically use a binary tree. I'd say most cases are summarized by:

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

	// work with the tree, i.e. go through nodes
	explore(t.Root())

	// maybe more writes, more work
}
```

Well, at least kinda. This is obviously oversimplified and your code most likely never looked like this. What I wanted to show here is the lifecycle of the object. 

However, there's actually a good chance you used a tree-like data structure in your code, just not directly. In "real world" use-cases, they're a common tool powering other abstractions. The [std::map](https://en.cppreference.com/cpp/container/map) and [std::set](https://en.cppreference.com/cpp/container/set) C++ containers are implemented as Red-Black trees and I would argue there is production code out there with sets and maps. Those have a different API, but the operations happening behind these abstractions use a "tree-like" API. Just have a look at C++'s [std::set](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_set.h) implementation which uses this [Red-Black tree](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_set.h#L137) right here to see for yourself. Of course, that's assuming you can even read that shit. I kept looking at the set file for a while and wondering, "where's the fucking set, though?"

> To be "that guy", it's actually `gcc`'s implementation and it's actually [this](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_tree.h) Red-Black tree. Have a look at the [insertion code](https://github.com/gcc-mirror/gcc/blob/b1263bca56ea82e555ddccf6729042a25805580c/libstdc%2B%2B-v3/include/bits/stl_tree.h#L764). That's a total of 307 underscores!

As a side note, I personally have use cases where I use trees directly, which is when I do coding problems and [trees](https://github.com/Ozoniuss/Algorithms/commit/5d8d7f809c220913afab2c51a5feb266cd53e2ad) happen to show up. In those cases my code does look similar to the snippet, but typically only implementing the "work with the tree, i.e. go through nodes" part. Yeah, I find a lot of joy in figuring out if I want to go left or right. What I don't find a lot of joy in though (at least with leetcode), is the "do some writes on the tree with values coming from somewhere" part, because getting those values from their inputs to test my stuff locally has to be at least as complicated as writing the actual solution to the problem.

Most of the time, your input looks something like this:

![alt text](image-3.png)

Here, it comes from [this](https://leetcode.com/problems/count-nodes-equal-to-average-of-subtree/) problem. This is perfectly reasonable, as long as you represent it with an image. That's not the case in your code, though. Leetcode weren't kind enough to provide the implementation of an `Insert` method like I have in my sample. Instead, it ends up looking something like this:

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

This is only 6 nodes, by the way. And you don't even need to specify the color of a node, like in a Red-Black tree.

I think at this point you're starting to realize the direction I'm going to. Even for relatively small trees, this representation totally lacks clarity. Even though you could argue having to write this may not be a pain in the ass anymore now with AI, there's still no way you can give someone that snippet and have them understand how the tree looks like. You can imagine why for tests, this becomes quite a significant downside. Or, you could of course argue that reading code (particularly tests) isn't a problem anymore either nowadays.

Admittedly, there's ways to simplify writing these inputs. One example is writing it as a list, which is how LeetCode provides them:

```
root = [4,8,5,0,1,null,6]
```

Also, to give them some credit, you can now edit the inputs and see the tree with a visualizer:

![alt text](image-5.png)

This representation is much more compact, but it still has the same downside of lacking clarity. I was actually surprised the first time I saw how `[4,null,5,null,6]` looks: 

![alt text](image-6.png)

I suppose eventually, you'd get used to whatever representation you use most commonly. But there's no denying that images are simply a much, much better representation of a binary tree input. I actually ended up [writing one](https://github.com/Ozoniuss/ds/blob/8ca58b28055dfc3b02b78dea8a271187d7cab0f8/rbtree/format.go#L73) so that I can have visual tree representations in my source code! Here's a snippet from how my [test suite](<https://github.com/Ozoniuss/ds/blob/8ca58b28055dfc3b02b78dea8a271187d7cab0f8/rbtree/rbt_visual_test.go)>) looked like at some point:

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
//
```

I'll be covering how it works in a future article. I find this a significant step up in terms of clarity but still, wasn't good enough for me to be satisfied with how much it proves correctness. On top of it, you actually need to, well, write these inputs:) No wonder AI generates all your test cases now. When you start to think a bit more about it, I'm basically paying a bunch big-ass corporate monopolies to write test cases for my open source library. I prefer not to think about it.


Let's dive into correctness a bit more. So far, I kind of implied that your testing approach is something like this:

- Start with a set of inputs;
- Apply the operation on that set of inputs;
- Evaluate the result of the operation against an expected result.

This is very basic but also a good description of how most unit test cases are set up. Let's be honest, you probably just want to pass that coverage requirement, don't you? While I find it suitable for reasonably simple logic, it quickly becomes tricky for Red-Black tree operations. To put it into perspective, here's how the insert implementation [looks like](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree#Insertion). The book that I've been reading dedicates 6 pages just to explain the recoloring procedure you may or may not need to execute at the end of the insertion. Here's a snippet. ![alt text](image-4.png)

Sure, having to read through this isn't the worst thing in the world (remember those C++ snippets from earlier?). Nonetheless, it's complex enough to wrap your head around and inevitably doesn't give me enough confidence that I have a sufficient number of test cases to cover all scenarios. In theory, those are finite (textbook classifies it in 3 cases) but not in the sense that each test input will pass through one of these 3 cases. Rather, the intermediary state during an insertion can be one of these 3 cases. There are insertion operations that can pass all cases, some multiple times. Even if you re-write the code so that you can test each insertion step, there still is the question of what your inputs would be to the parent insertion operation.

As if this weren't enough, there's some additional problems which I won't cover in depth, but raise important considerations about how you should design your test suite:

- The visual trees add clarity in representation. They don't necessarily add clarity in understanding what the hell happened after the insert. (add image) [what happened there](https://github.com/Ozoniuss/ds/blob/8ca58b28055dfc3b02b78dea8a271187d7cab0f8/rbtree/rbt_visual_test.go#L579)
- I'm not sure if every structurally valid Red-Black tree can be reproduced using a random sequence of insertions and deletions, at least with the algorithm that I'm using. One might argue an input like that does not constitute a valid test case, though personally I think it should.

But perhaps the one worth pointing out separately is that, while the actual implementation of the insertion operation is deterministic, there are multiple valid rebalancing implementations. The same input can produce different valid trees under different rebalancing algorithms, meaning that a test suite checking the output of a certain input isn't portable if you want to change your implementation say, to a more efficient one.

At this point, I suggest to start thinking about these tests differently. Particularly, the Red-Black tree definition only defines a set of invariants (constraints) that must hold after every operation.


So for this reason, it's time to take a step back and look at what it actually means to be a Red-Black tree:

0. It must satisfy the binary search tree property
1. The root is black
2. Every node is colored either red or black
3. Every leaf (TNIL) is black.
4. If a node is red, then both its children are black
5. For each node, all simple paths from the node to descendant leaves contain the same number of black nodes.

Like I said, this article isn't about Red-Black trees, so it's not actually worth spending time going into details about each one. The only important part is the result satisfying these constraints provides: _the tree is approximately balanced_. This effectively allows for implementations of write and read methods that run in logarithmic time with respect to the number of nodes in the tree. On a higher level, it offers a means to efficiently implement a collection of sorted elements.

I want to highlight the part where I say "allows for implementations". On their own, these constraints only define valid tree structures on an abstract level. Just because you model that representation in code does not automatically mean that your operations actually run efficiently. To give an example, you may have a list of elements that is sorted in increasing order. You can use binary search to look for an element efficienty, but nothing stops you from reading the list sequentially until you find the element. This is important because testing those constraints alone does not automatically mean that your implementation is efficient.

Still, given these constraints do allow efficient implementations, it's extremely important to validate them. The great part about constraints is that they must hold on any state of the Red-Black tree object to create. This is a substantial benefit when designing a test suite, because _we don't actually care anymore about test inputs_. Let's go back to the beginning of the article where I showed how you generally use a tree:

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

	// work with the tree, i.e. go through nodes
	explore(t.Root())

	// maybe more writes, more work
}
```

If we switch from testing specific inputs to testing that constraints hold, we can now effectively define tests that follow the same pattern. By generating randomized lists of write operations and validating the constraints after every write, we gain enough confidence that the structure of the tree is valid.

Of course, this isn't a substitution for testing the public API. Combined with constraint-based testing, we can now include additional tests using the public interface's methods alone. To give a concrete example of how you can achieve this, you could:


- Generate a random list of elements;
- Insert every element that has an even position in the list;
- Check the constraints hold after every insertion;
- Assert that searching all elements with even position in the input list returns a valid node;
- Assert that searching all elements with an odd position in the input list returns a null node.

This will highly depend on the definition of your public interface, though, particularly the write and read method's input and output parameters. In my case, insertions don't return anything and reads return a node pointer, so the setup above makes sense. Consider appropriate adjustments to the testing setup based on the API you want to expose. If the tree accepted duplicates and the insertion method returns the number of elements in the tree, your test input may have duplicates and after each insertion, you could test that it actually returns the proper value. There's just so many ways to play around and explore with the design of your interface, so have some fun with it:)

Another nice things about constraints is that you can use them to derive other constraints. For Red-Black trees, one particularly interesting one is that a subtree rooted at a node `x` contains at least `2^bh(x) - 1` internal nodes. In turn, this implies that the height of any tree satisfying it is at most `2*lg(n+1)`, where `n` is the number of nodes. That gives you more possiblities when deciding which sets of constraints you might want to validate in your tests, because some may have a much more efficient algorithm to actually validate. Just be mindful that, if you decide to test a different set of constraints than the original ones, it's important they imply the original set of constraints. The second property I stated above _does not_ imply the former.

I like math a lot, so I like to think of the first 6 constraints I'm starting from as some axioms which define a space that I can use to keep playing around and exploring its properties. The amount of things we can deduce starting from a very small set of properties is amazing.

> Strictly speaking, there's more to it than just the starting 6 axioms. You could go even further and ask yourself, "What does a tree, a node or a child actually mean?" but here I find it important to be practical. As long as the underlying language used to define these "axions" is the same for everyone (e.g. it is clear what a child node means), then the constraints are a good starting point to explore.

Taking a step back to the Red-Black tree library tests. There is an efficient way to validate the 6 original constraints using a single tree traversal. This gives us a lot of flexibility in our tests, since we now can write tests on large randomized sets of inputs. Validating the public interface works as expected combined with the fact that the shape of the tree is correct gives me plenty of confidence in the testing suite.

If you're interested, here is the algorithm to validate the constraints:

(insert code)

Also, this is the shape of the tests:

(insert code)



https://en.wikipedia.org/wiki/The_Power_of_10:_Rules_for_Developing_Safety-Critical_Code taking a next step further?




----

Mostly yes. The "why" is clear, and you cover two real limitations. The section is missing the most important limitation, though, and a few "be careful" points. Some of them close loops you opened in the first half.

Why it's useful
The core point lands: constraints hold no matter what the input was, so you can throw random write sequences at the tree. Reusing the main() snippet works well for that.

What's missing is a tie back to the problems you raised earlier. You'll add portability, but the other two deserve a sentence each:

Writing down trees. With constraints, there's no expected tree to write down at all. That answers the whole first half of the article, and you never say it outright.
Case coverage (line 139). Random sequences hit many fixup case combinations, including multiple passes through the loop, without you having to design inputs that do it.
Limitations
What you have:

Constraints don't guarantee efficiency (line 162).
Constraints don't replace testing the public API (line 190).
What's missing:

A valid tree can still hold the wrong elements. This is the biggest gap. An Insert that does nothing keeps the empty tree perfectly valid. So does one that silently drops a value or a whole subtree. This is the real reason the public API tests are needed, so say it explicitly. Your even/odd test catches it partially. An in-order traversal compared against the sorted set of inserted values (plus a size check) catches it fully.

Your efficiency example targets the wrong thing. For a Red-Black tree, the constraints do bound the height, so a normal search on a valid tree is logarithmic. What they don't bound is the cost of writes: an Insert that rebuilds the whole tree from scratch passes every constraint check. That's a sharper example than the sorted list.

Internal helpers break constraints mid-operation. You already noted that transplant can leave the tree invalid (line 88). So constraints only hold between public operations, which means helpers like rotate and transplant can't be constraint-tested. They need example tests, which is exactly why the visual suite still exists. This is your natural bridge to the visual tests article, and it fixes the "any state" wording in line 164.

Coverage isn't guaranteed. Your old notes (lines 234 and 276) worried that random sequences might miss branches. That worry disappeared from the new section. You raised it in line 139, so readers expect an answer. Either the random tests plus a coverage report show every branch is hit, or you added explicit cases. The cover/profile.go links at the bottom suggest you have something here.

Reachability (line 144) is still unanswered. Random insert/delete sequences only produce trees that the algorithm itself can reach. The constraint checker works on hand-built trees too, though, so the two approaches complement each other. One sentence closes it.

What to be careful about
What you have: derived constraints must imply the originals (line 201).

What's missing:

Test the checker. If the validator has a bug, every test passes. Feed it trees that break each property, one per test.
Inputs still matter. "We don't actually care anymore about test inputs" (line 164) overstates it. Uniform random inserts rarely produce some patterns: sorted or reverse-sorted runs, duplicates, deleting the root, deleting missing values, long insert/delete interleavings. You don't care about specific inputs, but you do care about the distribution.
Log the seed. A failing random test is useless if you can't reproduce it.
Smaller content issues
Line 201, "the second property … does not imply the former", is ambiguous. I think you mean the height bound doesn't imply the Red-Black properties. An AVL tree is a good concrete example: it meets the height bound and doesn't even have colors.
Line 201's efficiency argument is weak. You already check all the original constraints in one O(n) pass, and checking a whole tree can't be cheaper than that. Derived constraints are better framed as extra checks that directly test what you care about (balance). Checking a consequence can only find bugs. Checking a set of constraints that implies the originals is what proves the tree valid.
Line 207 says "the shape of the tree is correct". Given your argument, it should be "valid". "Correct" is exactly the word the article argues against.
Name-drop property-based testing. This idea is usually called property-based testing (QuickCheck, testing/quick, rapid). One line naming it helps readers who know the term and fits your "it's an idea, not a guide" framing. You can keep calling it constraint-based.
I'd leave the axioms aside (lines 203–205) as it is. It's your voice and costs nothing. Just fix "axions"

---

To be honest I think these are avlid (particularly 1, for 2 I said that writes CAN be implemented efficiently), but we're already besides the point for the others. You can't cover everything with a test suite but this applies to any test suite. (reachability also actually does not matter tbh). I think those are better pointed out in an appendix. I think I want my article to get the point that you can reach a high level of correctness if you just prove the constraints. obivoulsy the means in which you do prove correctness can, in turn, be defective. but tthat will just create an endless loop -- now we prove the means we used to prove right?

---
That works. Here's where I'd put each point:

- **Keep #1 in the main body.** It's the one that directly challenges your thesis: a tree can satisfy every constraint and still be missing an element you inserted. You already answer it with the public API tests. One sentence on *why* those tests are needed makes "constraints + API tests = high confidence" hold together.
- **#2: you're right.** Line 162 does say the constraints *allow* efficient implementations. I'll drop it.
- **Everything else (coverage, seeds, input distribution, internal helpers) goes in the appendix.** They're practical notes, not part of the argument.

On the endless loop, I agree you can't prove the prover forever. There's a way out of it, though, and it actually supports your thesis. Checking the constraints is a single, simple traversal. Rebalancing is six pages of textbook. Constraint-based testing moves your trust from the complicated code to the simple code. You still trust *something*, but it's something you can read and verify by eye. That's why you can stop there.

I'd add that as one sentence near the end of the section: "the checker is much simpler than the code it checks, and that's the whole point." It answers the loop before readers bring it up.