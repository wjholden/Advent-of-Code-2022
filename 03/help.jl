words = readlines("input.txt")
l = [w[1:length(w)÷2] for w in words]
r = [w[length(w)÷2+1:end] for w in words]
function priority(s)
	c = first(s)
	if 'A' <= c <= 'Z'
		Int(c) - Int('A') + 27
	else
		Int(c) - Int('a') + 1
	end
end
using Pipe
@pipe [priority.(intersect(Set(split(l[i], "")), Set(split(r[i], "")))) for i in 1:length(l)] |> sum
# find strings with no intersection -- this was my bug in Go
filter(isempty, intersect.(zip(Set.(split.(l, "")), Set.(split.(r, "")))))