# This will overflow on Gecode and Chuffed --- use OR-Tools!

x = 37503495108131
n = 20

print("array[1..20] of var -2..2: X;")
print()

print(
    f"""constraint {x} == {
        " + ".join(
            [f'{5**(i-1)} * X[{i}]' for i in range(1,n+1)]
        )
    };

solve satisfy;

output[
    "{"".join([f'\\(X[{n - i + 1}])' for i in range(1,n+1)])}"
]
"""
)