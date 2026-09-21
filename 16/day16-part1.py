import sys, parser, z3

if len(sys.argv) < 2:
    print("must provide input file as first argument")
    exit()

valves = parser.get_valves(sys.argv[1])
valve_index = parser.get_valve_indices(valves)

# def distances(src):
#     global valves
#     d = {}
#     q = queue.PriorityQueue()
#     q.put((0, src))
#     while not q.empty():
#         (distance, u) = q.get()
#         if not u in d:
#             d[u] = distance
#         for v in valves[u][1]:
#             if not v in d:
#                 q.put((distance + 1, v))
#     return d

# def opportunity(src, t):
#     global valves
#     vertex = None
#     value = -1
#     distance = distances(src)
#     for (u, dist) in distance.items():
#         payout = (t - dist - 1) * valves[u][0]
#         # print(f"Expected payout from {u} is {payout}")
#         if payout > value:
#             vertex = u
#             value = payout
#     # print(f"The greedy move is to select {k} with {v} payoff")
#     return vertex, value, distance[vertex]

# print(distances('AA'))
# opportunity('AA', 30)

# print("digraph {")
# for (valve, (_, neighbors)) in valves.items():
#     for neighbor in neighbors:
#         print(f"  {valve} -> {neighbor};")
# print("}")

# t = 30
# pressure = 0
# position = 'AA'
# while t > 0:
#     v, p, d = opportunity(position, t)
#     print(f"Move to {v} in {d} for payoff {p}")
#     t = t - d - 1
#     pressure += p
#     position = v
#     valves[position][0] = 0

# print("Part 1:", pressure) # 1790 too low.

# 1) Initialization
m_valves = len(valves)
n_times = 31
optimizer = z3.Optimize()

# 2) Decision variables
x = [[z3.Bool(f'at valve {v} at time {t}') for t in range(n_times)] for v in valves.keys()]
o = [z3.Int(f'valve {v} is open for this amount of time') for v in valves.keys()]

# 3) Constraints
# We start at position AA in the first time slot (t=0).
aa = valve_index['AA']
optimizer.add(x[aa][0])

# We are at exactly one valve for each position in time.
for t in range(n_times):
    optimizer.add(z3.AtMost(*[x[valve_index[v]][t] for v in valves.keys()], 1))

# Being at position x[i][t] constrains x[j][t+1] to staying in the same
# position or moving to one of the adjacent valves.
for t in range(n_times-1):
    for u, valve in enumerate(valves.keys()):
        adjacent = [valve_index[y] for y in valves[valve][1]]
        # print('Valve',valve,'is adjacent to',adjacent)
        adjacent_and_self = adjacent + [u]
        constraint = z3.Implies(x[u][t], z3.Or(*[x[v][t+1] for v in adjacent_and_self]))
        optimizer.add(constraint)

# Decide when a valve was opened, if at all.
for valve, (flow_rate, _) in valves.items():
    v = valve_index[valve]
    # Default to zero if never opened.
    # If there is no x[i][t] and x[i][t+1] that are both true, then o=0.
    c1 = z3.Implies(
        z3.Not(z3.Or(*[
            z3.And(x[v][t], x[v][t+1])
            for t in range(n_times-1)
        ])), o[v] == 0)
    optimizer.add(c1)
    # Otherwise, it starts at t+1.
    for t in range(n_times-1):
        c2 = z3.Implies(
            z3.And(x[v][t], x[v][t+1]),
            o[v] == n_times - (t+2)
        )
        optimizer.add(c2)
    # In all cases, the valve is open for nonnegative time.
    c3 = o[v] >= 0

# 4) Objective
obj = optimizer.maximize(z3.Sum(*[
    o[valve_index[i]] * flow_rate for
    (i, (flow_rate, _)) in valves.items()
]))

# 5) Solve
if optimizer.check() == z3.sat:
    model = optimizer.model()

    for t in range(n_times):
        for v in range(m_valves):
            if model[x[v][t]]:
                print(x[v][t])
    
    for o in o:
        print(o, model[o])
    print("Part 1:", obj.value())
    if sys.argv[1] in ["example.txt", r".\example.txt"]:
        assert 1651 == obj.value()
else:
    print("Not satisfiable")

