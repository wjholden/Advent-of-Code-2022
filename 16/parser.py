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