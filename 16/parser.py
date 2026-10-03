import re

def get_valves(file):
    valves = {}
    pattern = r'Valve (?P<valve>[A-Z]{2}) has flow rate=(?P<rate>\d+); tunnels? leads? to valves? (?P<neighbors>.+)'
    with open(file, "r") as f:
        for valve, flow_rate, neighbors in re.findall(pattern, f.read()):
            valves[valve] = (int(flow_rate), neighbors.split(", "))
    return valves

def get_valve_indices(valves):
    valve_index = {}
    for i, v in enumerate(valves.keys()):
        valve_index[v] = i
    return valve_index

def get_valves_i(file):
    valves = {}
    vi = []
    pattern = r'Valve (?P<valve>[A-Z]{2}) has flow rate=(?P<rate>\d+); tunnels? leads? to valves? (?P<neighbors>.+)'
    with open(file, "r") as f:
        for valve, flow_rate, neighbors in re.findall(pattern, f.read()):
            if valve not in vi:
                vi.append(valve)
            for neighbor in neighbors.split(', '):
                if neighbor not in vi:
                    vi.append(neighbor)
            v = vi.index(valve)
            adj = [vi.index(neighbor) for neighbor in neighbors.split(', ')]
            valves[v] = (int(flow_rate), adj)
    va = [None] * len(valves)
    for (v, (f, adj)) in valves.items():
        va[v] = (f, adj)
    return va, vi.index('AA')
    