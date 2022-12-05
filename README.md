# Advent of Code 2022
 
This year: Go!

# Daily Themes and Stars

1. `**` Parsing input, simple loops, arrays, finding the max and sum of a small data set
2. `**` Conditions, cases, small problems with many solutions (like FizzBuzz, where you want a clever one-liner but don't need it)
3. `**` ASCII, strings, set intersection, composable functions (`f(f(x,y),z)`)
4. `**` Ranges, comparisons (`<`, `<=`), tricky boundary conditions, splitting strings
5. `**` Stacks and queues, array copying, deep copies of object arrays, instructions and assembly languages

# Lessons Learned
* Go is a really low-level language and does not provide batteries for things like `sum()` and `Set()`.
* Go ranges are `[inclusive:exclusive]` (which Donovan & Kernighan call "half-open")
* [Unused variables cause a compile-time error in Go](https://go.dev/doc/faq#unused_variables_and_imports)
* Just like Java and JavaScript, `copy()` works fine on a 1D array but it isn't deep.
