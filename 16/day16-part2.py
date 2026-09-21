import sys, parser, z3

if len(sys.argv) < 2:
    print("must provide input file as first argument")
    exit()

valves = parser.get_valves(sys.argv[1])
valve_index = parser.get_valve_indices(valves)

# Initialize solver/optimizer and decision variables.

m_valves = len(valves)
n_times = 26 + 1
optimizer = z3.Optimize()

# We will model the valves the same way: an n*m matrix of binary decisions
# that each represent our positions.

my_positions = [[z3.Bool(f'I myself am at valve {v} at time {t}')
    for t in range(n_times)] for v in valves.keys()]

elephant_positions = [[z3.Bool(f'Elephant is at valve {v} at time {t}')
    for t in range(n_times)] for v in valves.keys()]

# I'll model the on/off state of the valves a little different this time.
# We'll have a separate decision variable for each valve at each time.

flow_v_t = [[z3.Int(f'flow at valve {v} at time {t}')
    for t in range(n_times)] for v in valves.keys()]

# Now for the constraints. First, all valves start with a flow of zero.

for v in range(m_valves):
    optimizer.add(flow_v_t[v][0] == 0)

# We must be at exactly one position at each time.
for t in range(n_times):
    optimizer.add(z3.AtMost(*[my_positions[v][t] for v in range(m_valves)], 1))
    optimizer.add(z3.AtMost(*[elephant_positions[v][t] for v in range(m_valves)], 1))
    optimizer.add(z3.Or(*[my_positions[v][t] for v in range(m_valves)]))
    optimizer.add(z3.Or(*[elephant_positions[v][t] for v in range(m_valves)]))

# The next position must be either adjacent or the same as the current.
for t in range(n_times-1):
    for valve, (_, neighbors) in valves.items():
        adj = [valve_index[u] for u in neighbors]
        v = valve_index[valve]
        adj.append(v)
        optimizer.add(
            z3.Implies(my_positions[v][t],
                       z3.Or(*[
                           my_positions[u][t+1] for u in adj
                       ])))
        optimizer.add(
                    z3.Implies(elephant_positions[v][t],
                               z3.Or(*[
                                   elephant_positions[u][t+1] for u in adj
                               ])))

# If I am at a valve for two moves in a row, then on the third move the flow
# takes the value from the associated flow_rate.

for (valve, (flow_rate, _)) in valves.items():
    v = valve_index[valve]
    # Skip the ones where flow rate is zero.
    if flow_rate == 0:
        # If the flow rate is zero then this valve can never contribute.
        for t in range(n_times):
            optimizer.add(flow_v_t[v][t] == 0)
    else:
        for t in range(n_times - 2):
            optimizer.add(
                z3.Implies(
                    z3.Or(
                        z3.And(my_positions[v][t], my_positions[v][t+1]),
                        z3.And(elephant_positions[v][t], elephant_positions[v][t+1])
                    ),
                    flow_v_t[v][t+2] == flow_rate   
                )
            )
        # Once we open a valve it stays open. Otherwise, the flow rate is zero.
        for t in range(n_times - 1):
            c = z3.Implies(
                                flow_v_t[v][t] > 0,
                                flow_v_t[v][t+1] == flow_rate,
                                flow_v_t[v][t+1] == 0
                            )
            optimizer.add(
                c
            )
            print(c)
        # The flow rate is either zero or the rate.
        for t in range(n_times):
            optimizer.add(z3.Or(flow_v_t[v][t] == 0, flow_v_t[v][t] == flow_rate))

# Last, our starting position. We both start at AA.
optimizer.add(my_positions[valve_index['AA']][0] == True)
optimizer.add(elephant_positions[valve_index['AA']][0] == True)

# Our objective is to maximize the sum of all of those flows.

objective = optimizer.maximize(
    z3.Sum(
        *[col for row in flow_v_t for col in row]
    )
)

if optimizer.check() == z3.sat:
    model = optimizer.model()

    # for t in range(n_times):
    #     for v in range(m_valves):
    #         if model[my_positions[v][t]]:
    #             print(my_positions[v][t])
    #         if model[elephant_positions[v][t]]:
    #             print(elephant_positions[v][t])

    for v in range(m_valves):
        for t in range(n_times):
            print(flow_v_t[v][t],"=",model[flow_v_t[v][t]])

    # print(model)
    # for row in my_positions:
    #     for col in row:
    #         print(model[col])
else:
    print("not satisfiable")