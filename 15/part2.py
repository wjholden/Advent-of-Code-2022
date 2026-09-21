import re

import z3

pattern = r'Sensor at x=(?P<sx>-?\d+), y=(?P<sy>-?\d+): closest beacon is at x=(?P<bx>-?\d+), y=(?P<by>-?\d+)'
with open("day15.txt", "r") as f: # use example.txt for the example...obviously.
    sensors_and_beacons = re.findall(pattern, f.read())

sensors_and_beacons = [[int(x) for x in t] for t in sensors_and_beacons]

def distance(sx, sy, bx, by):
    return abs(sx - bx) + abs(sy - by)

left_boundry = sensors_and_beacons[0][0]
right_boundry = left_boundry

for i in range(0, len(sensors_and_beacons)):
    [sx, sy, bx, by] = sensors_and_beacons[i]
    d = distance(sx, sy, bx, by)
    left_boundry = min(left_boundry, sx - d)
    right_boundry = max(right_boundry, sx + d)
    sensors_and_beacons[i].append(d)

# cannot_contain_beacon = 0
# y = 2000000 # 10 in the example
# for x in range(left_boundry, right_boundry + 1):
#     # If (x,y) is inside the range of a given sensor/beacon distance then
#     # it cannot possibly contain a beacon.
#     for [sx, sy, bx, by, d1] in sensors_and_beacons:
#         d2 = distance(sx, sy, x, y)
#         if (bx, by) == (x, y):
#             break
#         if d2 <= d1:
#             cannot_contain_beacon += 1
#             break

# print("Part 1:", cannot_contain_beacon)

x = z3.Int('x')
y = z3.Int('y')

solver = z3.Solver()

solver.add(0 <= x)
solver.add(x <= 4000000) # 20 in the example
solver.add(0 <= y)
solver.add(y <= 4000000) # 20 in the example

for [sx, sy, _bx, _by, d] in sensors_and_beacons:
    solver.add(z3.Sum(z3.Abs(x - sx), z3.Abs(y - sy)) > d)

if solver.check() == z3.sat:
    print("Part 2:", solver.model().eval(x * 4000000 + y))