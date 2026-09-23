import sys, parser
from z3 import Optimize, Int, Bool, And, Implies, Or, AtMost, AtLeast, Sum, If, sat, PbEq

if len(sys.argv) < 2:
    print("must provide input file as first argument")
    exit()

valves = parser.get_valves(sys.argv[1])
valve_index = parser.get_valve_indices(valves)

m_valves = len(valves)
n_times = 26 + 1
opt = Optimize()

# m*n matrices of bools that represent if we or the elephant was at a position.
me = [[Bool(f'Me: {v} at t={t:02}') for t in range(n_times)] for v in valves.keys()]
el = [[Bool(f'El: {v} at t={t:02}') for t in range(n_times)] for v in valves.keys()]

# Starting position is 'AA' at t=0.
opt.add(me[valve_index['AA']][0])
opt.add(el[valve_index['AA']][0])

# You can occupy exactly one position at a time.
for t in range(n_times):
    opt.add(AtLeast(*[me[v][t] for v in range(m_valves)], 1))
    opt.add(AtLeast(*[el[v][t] for v in range(m_valves)], 1))
    opt.add(AtMost(*[me[v][t] for v in range(m_valves)], 1))
    opt.add(AtMost(*[el[v][t] for v in range(m_valves)], 1))

# You can only either stay in place or advance to an adjacent position.
for t in range(n_times-1):
    for valve, (_, adj) in valves.items():
        v = valve_index[valve]
        neighbors = [valve_index[u] for u in adj]
        opt.add(Implies(me[v][t], Or(me[v][t+1], *[me[u][t+1] for u in neighbors])))
        opt.add(Implies(el[v][t], Or(el[v][t+1], *[el[u][t+1] for u in neighbors])))

# Another m*n matrix of integers. These represent flows.
flows = [[Bool(f'{v} open at t={t:02}') for t in range(n_times)] for v in valves.keys()]

# All valves are initially closed.
for v in range(m_valves):
    opt.add(flows[v][0] == False)

# If we've occupied a position for two turns in a row then the valve opens on the
# next turn.
for valve in valves.keys():
    v = valve_index[valve]
    flow_rate = valves[valve][0]
    for t in range(1, n_times):
        # An effort to use fewer constraints. Doesn't really fix the combinatorial
        # complexity.
        if flow_rate == 0:
            opt.add(flows[v][t] == False)
            continue

        # This was my bug. If we aren't careful to constraint the "1th" position
        # then the value will be unconstrained and the solver will be able to assign
        # a meaningless value.
        if t == 1:
            opt.add(flows[v][t] == flows[v][t-1])
        else:
            opt.add(If(Or(And(me[v][t-1], me[v][t-2]), And(el[v][t-1], el[v][t-2])),
                    flows[v][t] == True,
                    flows[v][t] == flows[v][t-1]))
            # Another great suggestion from Gemini.
            # opened = Or(And(me[v][t-2], me[v][t-1]), And(el[v][t-2], el[v][t-1]))
            # opt.add(flows[v][t] == Or(flows[v][t-1], opened))

# The objective is to maximize the total of the flows.
obj = opt.maximize(Sum([
    If(flows[valve_index[valve]][t], valves[valve][0], 0)
    for t in range(n_times) for valve in valves.keys()
]))

if opt.check() == sat:
    model = opt.model()
    for t in range(n_times):
        for v in range(m_valves):
            # print(flows[v][t], model[flows[v][t]])
            if model[me[v][t]]:
                print(me[v][t], model[me[v][t]])
            if model[el[v][t]]:
                print(el[v][t], model[el[v][t]])
            if model[flows[v][t]]:
                print(flows[v][t], model[flows[v][t]])
        print()
    print(obj.value())
else:
    print("not satisfiable")