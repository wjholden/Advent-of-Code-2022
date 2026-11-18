#import "@preview/cetz:0.5.3"

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

= Top

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
x &= c - 2 cal(l) = c - "col" cal(l)\
y &= cal(l) - 1\
z &= r = r - 0 cal(l) = r - "row" cal(l)
$

In matrix form,

$
vec(x,y,z) = mat(
  0,1,-2,0;
  0,0,1,-1;
  1,0,0,0
) vec(r,c,cal(l),1)
$

Tested in Julia:

```julia
julia> [0 1 -2 0; 0 0 1 -1; 1 0 0 0] * [0 0 3 3; 8 11 8 11; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  3  0  3
 3  3  3  3
 0  0  3  3
```

= Front

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
vec(x,y,z) = mat(
  0,1,-2,0;
  -1,0,2,-1;
  0,0,1,-1
) vec(r,c,cal(l),1)
$

```julia
julia> [0 1 -2 0; -1 0 2 -1; 0 0 1 -1] * [4 4 7 7; 8 11 8 11; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  3  0  3
 3  3  0  0
 3  3  3  3
```

= Left

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
vec(x,y,z) = mat(
  0, 0, 0, 0;
  -1, 0, 2, -1;
  0, 1, -1, 0
) vec(r,c,cal(l),1)
$

```julia
julia> [0 0 0 0; -1 0 2 -1; 0 1 -1 0] * [4 4 7 7; 4 7 4 7; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  0  0  0
 3  3  0  0
 0  3  0  3
```

= Back

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
vec(x,y,z) = mat(
  0,-1,1,-1;
  -1,0,2,-1;
  0,0,0,0
) vec(r,c,cal(l),1)
$

```julia
julia> [0 -1 1 -1; -1 0 2 -1; 0 0 0 0] * [4 4 7 7; 0 3 0 3; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 3  0  3  0
 3  3  0  0
 0  0  0  0
```

= Under

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
vec(x,y,z) = mat(
  0,1,-2,0;
  0,0,0,0;
  -1,0,3,-1
) vec(r,c,cal(l),1)
$

```julia
julia> [0 1 -2 0; 0 0 0 0; -1 0 3 -1] * [8 8 11 11; 8 11 8 11; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 0  3  0  3
 0  0  0  0
 3  3  0  0
```

= Right

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
vec(x,y,z) = mat(
  0,0,1,-1;
  0,1,-3,0;
  -1,0,3,-1
) vec(r,c,cal(l),1)
$

```julia
julia> [0 0 1 -1; 0 1 -3 0; -1 0 3 -1] * [8 8 11 11; 12 15 12 15; 4 4 4 4; 1 1 1 1]
3×4 Matrix{Int64}:
 3  3  3  3
 0  3  0  3
 3  3  0  0
```
