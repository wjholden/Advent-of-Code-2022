import sys
from functools import cache
from parser import get_valves_i
from tqdm import tqdm

if len(sys.argv) < 2:
    print("must provide input file as first argument")
    exit()

(valves, aa) = get_valves_i(sys.argv[1])

def masks(n):
    for i in range(2**n):
        yield tuple(1 if valves[b][0] > 0 and (1<<b) & i else 0 for b in range(n))

def masks_2():
    in_use = [index for index, (flow, _) in enumerate(valves) if flow > 0]
    print(in_use)
    n = len(in_use)
    for i in range(2**n):
        t = [0] * len(valves)
        for b in range(n+1):
            print(i,b)
            if b in in_use and b & i:
                t[b] = 1
        print(t)
        yield tuple(t)

def rmask(v):
    if not v:
        yield []
        return
    flow_rate, _ = v[0]
    for tail in rmask(v[1:]):
        if flow_rate > 0:
            yield [1] + tail
        yield [0] + tail  


def antimask(mask):
    return tuple(0 if b and valves[b][0] > 0 else 1 for b in mask)

@cache
def flow(time, position, open_valves, mask):
    current_flow = sum(valves[v][0] * mask[v] for v in open_valves)
    if time >= 30:
        return current_flow
    # The first case is where we open the current valve. Never open valves that
    # with a flow rate of zero. Subtract that valve's flow rate from the current
    # step because it takes a full minute to open and we can't overcount.
    return current_flow + max(
        flow(time + 1, position, open_valves.union({position}), mask) - valves[position][0] if position not in open_valves and mask[position] > 0 and valves[position][0] > 0 else float('-inf'),
        *[flow(time+1, neighbor, open_valves, mask) for neighbor in valves[position][1]]
    )

def part2():
    n = len(list(rmask(valves)))
    # print(sum(1 if flow_rate > 0 else 0 for flow_rate, _ in valves), len(list(rmask(valves))))
    # print(len(set(masks(len(valves)))))
    i = 1
    for mask in tqdm(rmask(valves), total = n):
        mask = tuple(mask)
        #me = flow(4, aa, frozenset(), mask)
        # flow.cache_clear()
        #me += flow(4, aa, frozenset(), antimask(mask))
        if i % 200 == 0:
            flow.cache_clear()
        i += 1
        #yield me
        yield flow(4, aa, frozenset(), mask) + flow(4, aa, frozenset(), antimask(mask))
    return
        

print("Part 1:", flow(0, aa, frozenset(), tuple([1] * len(valves))))
print("Part 2:", max(part2())) # 3332 too high.

