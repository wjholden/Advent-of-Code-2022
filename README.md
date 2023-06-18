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
13. `**` JSON parsing, nested arrays, comparators, sorting, 1-indexed arrays, edge cases, unclear specifications
14. `**` mutable state, if/elseif chains, large problems, coordinate systems, paths
15. `* ` might need to compare correctness with Julia for this one...

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
* Here's a little one-liner to measure execution time: `defer func(start time.Time) { fmt.Println(time.Since(start)) }(time.Now())`.
* Assignments are allowed in `if` conditions. For example, `if _, ok := d[x]; !ok { ... }`.
* Everything implements `interface{}`, so you can upcast anything to that for hillbilly generics. (I haven't learned Go's actual generics yet).
* Go *really* doesn't want you to use inheritance. It turns out that type assertions (`t.(x)`) are only for interfaces, not structs. I think it's impossible to upcast an *embedding* struct to its *embedded* struct (`type A struct { ... }` and `type B struct { A ... }`).
* Initializing a struct literal with an embedded struct requires more syntax than you'd expect.
* Interfaces define only member methods, not member variables.
* You cannot check if a channel is closed.
* The [`sync.WaitGroup`](https://pkg.go.dev/sync#WaitGroup) is really cool, but if you're passing it to a function then make sure you pass a pointer. Otherwise, you'll pass by copy and then deadlock. This is different from [channels](https://go.dev/tour/concurrency/2), which you can somehow pass by value. The inconsistency can be confusing. Maybe people usually just declare their `WaitGroup` as a private global variable. Another option might be to wrap the target function of your [goroutine](https://go.dev/tour/concurrency/1) in an anonymous function. For example:

```
var wg sync.WaitGroup
for _, x := range X {
    wg.Add(1)
    
    go f(x, wg) // Wrong! wg needs to be a pointer.

    go f(x, &wg) // Better, but kinda ugly that the function knows about wg.

    go func() { // Best, a general solution for arbitrary functions.
        f(x)
        wg.Done()
    }()
}
wg.Wait()
```

* Iterate over the key/value pairs of a `map` using `range`.
* Go can `switch` on channels. If you have some concurrent subscriber that needs to dynamically dispatch on inputs from many publishers, then use multiple channels. Don't try to use a single channel with many interfaces. You cannot (at least, not easily/idiomatically) switch on `instanceof` like you could in Java.
* `go run` is a nice way to quickly test programs written in Go. The compiler is so fast that this feels like a scripting language.
* `go get` and `go install` from Github is difficult and frustrating to get right. It works great with the [gopl.io examples](https://github.com/adonovan/gopl.io/), but I can't get it to work correctly for myself.

# References

* [The Go Programming Language](https://www.amazon.com/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440/)
* [Google I/O 2012 - Meet the Go Team](https://www.youtube.com/watch?v=sln-gJaURzk)
* [How to start a Go project in 2023](https://boyter.org/posts/how-to-start-go-project-2023/)
