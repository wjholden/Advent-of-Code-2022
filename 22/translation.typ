#import "@preview/cetz:0.5.3"

#title[Advent of Code 2022 Day 22]

#align(center)[
  #cetz.canvas(length: 4em, {
    import cetz.draw: *

    // Grid positions as (column, row) starting from bottom-left (0, 0)
    let squares = (
      (2, 2),                 // Top row, 3rd column
      (0, 1), (1, 1), (2, 1), // Second row, 1st–3rd columns
      (2, 0), (3, 0),         // Third row, 3rd–4th columns
    )

    for (x, y) in squares {
      rect((x, y), (x + 1, y + 1))
    }
  })
]

= Transforming rows and columns into $(x,y,z)$ coordiantes

== Top

The top face for the unfolded cube with side length $cal(l) = 4$.
In the unfolded pattern, $"row" = 0$ and $"col" = 2$.
$(r,c)$ pairs map to $(x,y,z)$ coordinates in the cube:

$
  (0, 8) &mapsto (0,3,0) \
  (0, 11) &mapsto (3,3,0) \
  (3, 8) &mapsto (0,3,3) \
  (3, 11) &mapsto (3,3,3) \
$

The functions to transform $(r,c)$ to $(x,y,z)$ are:

$
x &= c - 2 cal(l) = c - 2 cal(l)\
y &= cal(l) - 1\
z &= r = r - 0 cal(l) = r - cal(l)
$

In matrix form,

$
vec(x,y,z) = T v = mat(
  0,1,-2,0;
  0,0,1,-1;
  1,0,0,0
) vec(r,c,cal(l),1)
$

Alternatively,

$
  vec(x,y,z) = A x+b = mat(0, 1; 0, 0; 1, 0) vec(r,c) + vec(-2cal(l), cal(l)-1, 0)
$

(Coincidentally, $A$ is exactly equal to the heading matrices we'll see below.)

With this, all we need to reverse this is to solve $ A^(-1) (vec(x,y,z) - b) = A^(-1) A x = x = vec(r,c) $

Tested in Julia:

```julia
julia> [0 1 -2 0; 0 0 1 -1; 1 0 0 0] * [0 0 3 3; 8 11 8 11; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  3  0  3
 3  3  3  3
 0  0  3  3
``` 

== Front

$
  "row" &= 1 \
  "col" &= 2 \
  (4, 8) &mapsto (0,3,3) \
  (4, 11) &mapsto (3,3,3) \
  (7, 8) &mapsto (0,0,3) \
  (7, 11) &mapsto (3,0,3) \
  x &= c - 2cal(l) \
  y &= -r + 2cal(l) - 1 \
  z &= cal(l) - 1
$

$
vec(x,y,z) = F v = mat(
  0,1,-2,0;
  -1,0,2,-1;
  0,0,1,-1
) vec(r,c,cal(l),1)
$

Alternatively,

$
  vec(x,y,z) = A x+b = mat(0, 1; -1, 0; 0, 0) vec(r,c) + vec(-2cal(l), 2cal(l)-1, cal(l) -1)
$

```julia
julia> [0 1 -2 0; -1 0 2 -1; 0 0 1 -1] * [4 4 7 7; 8 11 8 11; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  3  0  3
 3  3  0  0
 3  3  3  3
```

== Left

$
  "row" &= 1 \
  "col" &= 1 \
  (4,4) &mapsto (0,3,0) \
  (4,7) &mapsto (0,3,3) \
  (7,4) &mapsto (0,0,0) \
  (7,7) &mapsto (0,0,3) \
  x &= 0 \
  y &= -r + 2cal(l) - 1 \
  z &= c - cal(l)
$

$
vec(x,y,z) = L v = mat(
  0, 0, 0, 0;
  -1, 0, 2, -1;
  0, 1, -1, 0
) vec(r,c,cal(l),1)
$

Alternatively,

$
  vec(x,y,z) = A x+b = mat(0, 0; -1, 0; 0, 1) vec(r,c) + vec(0, 2cal(l)-1,-cal(l))
$

```julia
julia> [0 0 0 0; -1 0 2 -1; 0 1 -1 0] * [4 4 7 7; 4 7 4 7; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  0  0  0
 3  3  0  0
 0  3  0  3
```

== Back

$
  "row" &= 1 \
  "col" &= 0 \
  (4,0) &mapsto (3,3,0) \
  (4,3) &mapsto (0,3,0) \
  (7,0) &mapsto (3,0,0) \
  (7,3) &mapsto (0,0,0) \
  x &= -c + cal(l) - 1 \
  y &= -r + 2cal(l) - 1 \
  z &= 0
$

$
vec(x,y,z) = B v = mat(
  0,-1,1,-1;
  -1,0,2,-1;
  0,0,0,0
) vec(r,c,cal(l),1)
$

Alternatively,

$
  vec(x,y,z) = A x+b = mat(0, -1;-1, 0; 0, 0) vec(r,c) + vec(cal(l)-1,2cal(l)-1,0)
$

```julia
julia> [0 -1 1 -1; -1 0 2 -1; 0 0 0 0] * [4 4 7 7; 0 3 0 3; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 3  0  3  0
 3  3  0  0
 0  0  0  0
```

== Under

$
  "row" &= 2 \
  "col" &= 2 \
  (8,8) &mapsto (0,0,3) \
  (8,11) &mapsto (3,0,3) \
  (11,8) &mapsto (0,0,0) \
  (11,11) &mapsto (3,0,0) \
  x &= c - 2cal(l) \
  y &= 0 \
  z &= -r + 3cal(l) - 1
$

$
vec(x,y,z) = U v = mat(
  0,1,-2,0;
  0,0,0,0;
  -1,0,3,-1
) vec(r,c,cal(l),1)
$

Alternatively,

$
  vec(x,y,z) = A x+b = mat(0,1;0,0;-1,0) vec(r,c) + vec(-2cal(l),0,3cal(l)-1)
$

```julia
julia> [0 1 -2 0; 0 0 0 0; -1 0 3 -1] * [8 8 11 11; 8 11 8 11; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  3  0  3
 0  0  0  0
 3  3  0  0
```

== Right

$
  "row" &= 2 \
  "col" &= 3 \
  (8,12) &mapsto (3,0,3) \
  (8,15) &mapsto (3,3,3) \
  (11,12) &mapsto (3,0,0) \
  (11,15) &mapsto (3,3,0) \
  x &= cal(l) - 1 \
  y &= c - 3cal(l) \
  z &= -r + 3cal(l) - 1
$

$
vec(x,y,z) = R v = mat(
  0,0,1,-1;
  0,1,-3,0;
  -1,0,3,-1
) vec(r,c,cal(l),1)
$

Alternatively,

$
  vec(x,y,z) = A x+b = mat(0,0;0,1;-1,0) vec(r,c) + vec(cal(l)-1,-3cal(l),3cal(l)-1)
$

```julia
julia> [0 0 1 -1; 0 1 -3 0; -1 0 3 -1] * [8 8 11 11; 12 15 12 15; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 3  3  3  3
 0  3  0  3
 3  3  0  0
```

= $dif r$ and $dif c$

We need to translate $dif r$ and $dif c$ into a three-dimensional heading. This should be easy. Note that the unfolded cube has an inverted vertical axis (the top of the unfolded pattern has $r=0$), but everything below assumes a $+y$-up world.

#columns(2)[
== Top

$
  bold(H) &= cases(
  dif x &= dif c,
  dif y &= 0,
  dif z &= dif r,
  ) = mat(0,1;0,0;1,0) vec(dif r, dif c) \
  bold(N) &= vec(0, 1, 0)
$

== Front

$
  bold(H) &= cases(
    dif x &= dif c,
    dif y &= -dif r,
    dif z &= 0,
  ) = mat(0, 1; -1, 0; 0, 0) vec(dif r, dif c) \
  bold(N) &= vec(0, 0, 1)
$

== Left

$
  bold(H) &= cases(
  dif x &= 0,
  dif y &= -dif r,
  dif z &= dif c,
  ) = mat(0, 0; -1, 0; 0, 1) vec(dif r, dif c) \
  bold(N) &= vec(-1, 0, 0)
$

#colbreak()

== Back

$
  bold(H) &= cases(
  dif x &= -dif c,
  dif y &= -dif r,
  dif z &= 0,
  ) = mat(0, -1; -1, 0; 0, 0) vec(dif r, dif c) \
  bold(N) &= vec(0, 0, -1)
$

== Under

$
  bold(H) &= cases(
  dif x &= dif c,
  dif y &= 0,
  dif z &= -dif r,
  ) = mat(0, 1; 0, 0; -1, 0) vec(dif r, dif c) \
  bold(N) &= vec(0, -1, 0)
$

== Right

$
  bold(H) &= cases(
  dif x &= 0,
  dif y &= dif c,
  dif z &= -dif r,
  ) = mat(0, 0; 0, 1; -1, 0) vec(dif r, dif c) \
  bold(N) &= vec(1, 0, 0)
$

]

= New Directions

Updating heading is easy:

$
  bold(H)' = -bold(N)
$

We're just heading in the opposite direction of our old normal vector.
Likewise, your new normal vector is simply your old heading.

$
  bold(N)' = bold(H)
$

#pagebreak()
= My puzzle input

#align(center)[
  #cetz.canvas(length: 3em, {
    import cetz.draw: *

    // (row, column) from top to bottom, left to right.
    // Note that these row/column conventions are different from the other cetz diagram above.
    let squares = (
      (0, 1) /* top */, (0, 2) /* right */,
      (1, 1) /* front */,
      (2, 0) /* left */, (2, 1) /* under */,
      (3, 0) /* back */,
    )

    for (r, c) in squares {
      rect((c, 3 - r), (c + 1, 3 - r + 1))
    }
  })
]

#columns(2)[

== Top

#table(
  columns: 2,
  table.header($(r,c)$, $(x,y,z)$),
  $(4,0)$, $(0,3,0)$,
  $(4,3)$, $(3,3,0)$,
  $(7,0)$, $(0,3,3)$,
  $(7,3)$, $(3,3,3)$
)

$
  vec(x,y,z) = A x + b = mat(0, 1; 0, 0; 1, 0) vec(r, c) + vec(-cal(l),cal(l)-1,0)
$

== Left

#table(
  columns: 2,
  table.header($(r,c)$, $(x,y,z)$),
  $(8,0)$, $(0,3,3)$,
  $(8,3)$, $(0,0,3)$,
  $(11,0)$, $(0,3,0)$,
  $(11,3)$, $(0,0,0)$
)

$
  vec(x,y,z) = A x + b = mat(0,0;0,-1;-1,0) vec(r, c) + vec(0,cal(l)-1,3cal(l)-1)
$

== Front

#table(
  columns: 2,
  table.header($(r,c)$, $(x,y,z)$),
  $(4,4)$, $(0,3,3)$,
  $(4,7)$, $(3,3,3)$,
  $(7,4)$, $(0,0,3)$,
  $(7,7)$, $(3,0,3)$
)

$
  vec(x,y,z) = A x + b = mat(0, 1; -1, 0; 0, 0) vec(r, c) + vec(-cal(l),2cal(l)-1,cal(l)-1)
$

#colbreak()

== Right

#table(
  columns: 2,
  table.header($(r,c)$, $(x,y,z)$),
  $(0,8)$, $(3,3,0)$,
  $(0,11)$, $(3,0,0)$,
  $(3,8)$, $(3,3,3)$,
  $(3,11)$, $(3,0,3)$
)

$
  vec(x,y,z) = A x + b = mat(0,0;0,-1;1,0) vec(r, c) + vec(cal(l-1),3cal(l)-1,0)
$

== Under

#table(
  columns: 2,
  table.header($(r,c)$, $(x,y,z)$),
  $(8,4)$, $(0,0,3)$,
  $(11,4)$, $(0,0,0)$,
  $(8,7)$, $(3,0,3)$,
  $(11,7)$, $(3,0,0)$
)

$
  vec(x,y,z) = A x + b = mat(0,1;0,0;-1,0) vec(r, c) + vec(-cal(l),0,3cal(l)-1)
$

== Back

#table(
  columns: 2,
  table.header($(r,c)$, $(x,y,z)$),
  $(12,0)$, $(0,3,0)$,
  $(12,3)$, $(0,0,0)$,
  $(15,0)$, $(3,3,0)$,
  $(15,3)$, $(3,0,0)$
)

$
  vec(x,y,z) = A x + b = mat(1,0;0,-1;0,0) vec(r, c) + vec(-3cal(l),cal(l)-1,0)
$

]