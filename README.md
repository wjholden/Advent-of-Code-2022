# Advent of Code 2022
 
This year: Go!

# Daily Themes and Stars

1. `**` parsing input, simple loops, arrays, finding the max and sum of a small data set
2. `**` conditions, cases, small problems with many solutions (like FizzBuzz, where you want a clever one-liner but don't need it)
3. `**` ASCII, strings, set intersection, composable functions (`f(f(x,y),z)`)
4. `**` ranges, comparisons (`<`, `<=`), tricky boundary conditions, splitting strings
5. `**` stacks and queues, array copying, deep copies of object arrays, instructions and assembly languages
6. `**` loops, strings and characters, substrings (slices), refactoring to more abstract/general code, equality, sets
7. `**` n-ary trees, types (objects, structs) with supertype (interface) and polymorphism, tree traversal, subproblems, file systems, command-line interfaces, string parsing, pipelines (such as generator functions, goroutines, and iterators)
8. `**` matrices/2D arrays, greedy algorithms, optimization (stopping early to not waste time), loops, extrema, refactoring
9. `**` mutable structs vs pure functions, pointers, recursion,  sets, distance/direction (vectors), absolute values
10. `**` string building, modular arithmetic, simple assembly languages, tricky off-by-one errors
11. `**` big integers (unless you can avoid it, which you can!), coprimes, parsers, interpreters, procedural programming, non-parallelizable problems, pointers/mutable structs
12. `**` graph searching, Dijkstra's algorithm (maybe A* for part 1), passing functions as arguments, reductions, brute force is not the answer

# Lessons Learned
* Go is a really low-level language and does not provide batteries for things like `sum()` and `Set()`.
* `map[T]bool` works well as a set substitute.
* Go ranges are `[inclusive:exclusive]` (which Donovan & Kernighan call "half-open")
* [Unused variables cause a compile-time error in Go](https://go.dev/doc/faq#unused_variables_and_imports). (You can suppress the error with `_`).
* Just like Java and JavaScript, `copy()` works fine on a 1D array but it isn't deep.
* Go's object-orientation does not really accomodate control flow switched on types. In Java, you might have used the `instanceof` operator to downcast a variable to a more specific type. Go uses a `.` to test interface satisfiability, `type.(value)`, but I don't think you can downcast at all.
* Goroutines are great! Go has the easiest concurrency model of any language that I have used.
* You'd think the simple rule "values are always passed by value" would be easy to remember, but this can surprise you if you've come from languages (like Java and Python) that pass object **references**. If you want to mutate an object in a function (or even `for` loop), then you need pointers.
* Fortunately, pointers are not terribly difficult in Go. The trickiest thing was getting the syntax just right for taking pointers to array members. Here is a working example:

```
var A []int = []int{1, 2, 3, 4}
var x *int = &A[0]
*x = 1234
fmt.Println(A)
```

* You also need pointers for fast code if you want to avoid lots of copying.
* You cannot get a pointer to a dictionary member.
* `slice = slice[:0]` is a neat trick to clear the contents of a slice.
* You can deserialize arbitrary JSON to an `interface{}`. From there, you can use [type assertions](https://go.dev/tour/methods/15) which is loosly comparable to a [cast](https://docs.oracle.com/javase/tutorial/java/IandI/subclasses.html) but maybe more restrictive.
* Go's `for` keyword replaces `while` and allows for some interesting syntax, including `for {}` with no condition.
* Assignments are legal in an `if` clause.
* Go's `panic` is very different from [`throw`/`catch`](https://docs.oracle.com/javase/tutorial/essential/exceptions/throwing.html).
