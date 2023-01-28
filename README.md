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

# Lessons Learned
* Go is a really low-level language and does not provide batteries for things like `sum()` and `Set()`.
* `map[T]bool` works well as a set substitute.
* Go ranges are `[inclusive:exclusive]` (which Donovan & Kernighan call "half-open")
* [Unused variables cause a compile-time error in Go](https://go.dev/doc/faq#unused_variables_and_imports). (You can suppress the error with `_`).
* Just like Java and JavaScript, `copy()` works fine on a 1D array but it isn't deep.
* Go's object-orientation does not really accomodate control flow switched on types. In Java, you might have used the `instanceof` operator to downcast a variable to a more specific type. Go uses a `.` to test interface satisfiability, `type.(value)`, but I don't think you can downcast at all.
* Goroutines are great! Go has the easiest concurrency model of any language that I have used.
