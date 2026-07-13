package routing

import (
	"math"

	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

// edge is a directed hop from one node to an adjacent node, following a
// block of plain track between them.
type edge struct {
	To   types.Pos
	Dir  types.Dir // direction to head from the source node to take this edge
	Cost int       // number of tiles between the two nodes
}

type routeEntry struct {
	NextHop types.Dir
	Cost    int
}

// Manager builds a graph of routing nodes (signals/junctions/stations)
// connected by blocks, and answers "which direction should I go" queries for
// trains navigating toward a destination. It is rebuilt whenever blocks
// change, and computes routing tables lazily (cached per source node).
type Manager struct {
	Nodes map[types.Pos]*Node
	edges map[types.Pos][]edge
	// tables caches, per source node, the shortest-path next-hop to every
	// other reachable node. Cleared wholesale on Rebuild.
	tables map[types.Pos]map[types.Pos]routeEntry
}

func NewManager() *Manager {
	return &Manager{
		Nodes:  make(map[types.Pos]*Node),
		edges:  make(map[types.Pos][]edge),
		tables: make(map[types.Pos]map[types.Pos]routeEntry),
	}
}

// Rebuild recomputes the node graph from scratch and invalidates all cached
// routing tables. Cheap enough to call after any block change without
// bothering with incremental/dirty propagation.
func (m *Manager) Rebuild(w *world.World) {
	m.Nodes = make(map[types.Pos]*Node)
	m.edges = make(map[types.Pos][]edge)
	m.tables = make(map[types.Pos]map[types.Pos]routeEntry)

	for pos, tr := range w.Tracks {
		if station := w.StationAt(pos); station != nil {
			m.Nodes[pos] = &Node{Pos: pos, Kind: NodeStation, Name: station.Name}
		} else if tr.HasSignal() || tr.IsJunction() {
			m.Nodes[pos] = &Node{Pos: pos, Kind: classify(tr)}
		}
	}

	for pos := range m.Nodes {
		tr := w.Tracks[pos]
		for d := types.Dir(types.DirNorth); d <= types.DirWest; d <<= 1 {
			if tr.Direction&d == 0 {
				continue
			}
			to, cost, ok := m.walk(w, pos, d)
			if !ok {
				continue
			}
			m.edges[pos] = append(m.edges[pos], edge{To: to, Dir: d, Cost: cost})
		}
	}
}

// walk follows track from a node, starting in direction dir, until it
// reaches another node. Returns false if the track dead-ends, disconnects,
// or loops back on itself without ever reaching a node.
func (m *Manager) walk(w *world.World, from types.Pos, dir types.Dir) (types.Pos, int, bool) {
	curr := from
	d := dir
	cost := 0
	visited := map[types.Pos]bool{from: true}

	for {
		next := types.NextPos(curr, d)
		if visited[next] {
			return types.Pos{}, 0, false
		}

		nextTr := w.Tracks[next]
		if nextTr == nil || nextTr.Direction&types.OppositeDir(d) == 0 {
			return types.Pos{}, 0, false
		}

		cost++
		curr = next
		visited[curr] = true

		if _, isNode := m.Nodes[curr]; isNode {
			return curr, cost, true
		}

		if nextTr.Direction&d != 0 {
			continue // keep heading the same direction
		}

		found := false
		for dd := types.Dir(types.DirNorth); dd <= types.DirWest; dd <<= 1 {
			if dd == types.OppositeDir(d) {
				continue
			}
			if nextTr.Direction&dd != 0 {
				d = dd
				found = true
				break
			}
		}
		if !found {
			return types.Pos{}, 0, false
		}
	}
}

// NextHop returns the direction a train sitting at node pos should take to
// make progress toward dest. pos must be a node position (see m.Nodes) —
// callers should only ask at an actual decision point (a junction), since
// plain curves have only one non-backtracking direction anyway.
func (m *Manager) NextHop(pos types.Pos, dest types.Pos) (types.Dir, bool) {
	if pos == dest {
		return types.DirNone, false
	}

	entry, ok := m.routingTable(pos)[dest]
	if !ok {
		return types.DirNone, false
	}
	return entry.NextHop, true
}

func (m *Manager) routingTable(source types.Pos) map[types.Pos]routeEntry {
	if table, ok := m.tables[source]; ok {
		return table
	}
	table := dijkstra(source, m.edges)
	m.tables[source] = table
	return table
}

// dijkstra computes shortest-path cost and first-hop direction from source
// to every other reachable node. Graphs here (signals/junctions/stations)
// are small, so a simple O(V^2) implementation is plenty.
func dijkstra(source types.Pos, edges map[types.Pos][]edge) map[types.Pos]routeEntry {
	dist := map[types.Pos]int{source: 0}
	firstHop := map[types.Pos]types.Dir{}
	visited := map[types.Pos]bool{}

	for {
		curr, ok := closestUnvisited(dist, visited)
		if !ok {
			break
		}
		visited[curr] = true

		for _, e := range edges[curr] {
			nd := dist[curr] + e.Cost
			if d, ok := dist[e.To]; !ok || nd < d {
				dist[e.To] = nd
				if curr == source {
					firstHop[e.To] = e.Dir
				} else {
					firstHop[e.To] = firstHop[curr]
				}
			}
		}
	}

	table := make(map[types.Pos]routeEntry, len(dist))
	for pos, d := range dist {
		if pos == source {
			continue
		}
		table[pos] = routeEntry{NextHop: firstHop[pos], Cost: d}
	}
	return table
}

func closestUnvisited(dist map[types.Pos]int, visited map[types.Pos]bool) (types.Pos, bool) {
	best := math.MaxInt
	var pos types.Pos
	found := false
	for p, d := range dist {
		if visited[p] {
			continue
		}
		if d < best {
			best = d
			pos = p
			found = true
		}
	}
	return pos, found
}
