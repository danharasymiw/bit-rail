package engine

import (
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/google/uuid"
)

type blockManager struct {
	w *world.World

	// tiles modified since game tick
	dirty map[types.Pos]struct{}
}

func newBlockManager(w *world.World) *blockManager {
	return &blockManager{
		w:     w,
		dirty: make(map[types.Pos]struct{}, 256),
	}
}

// MarkDirty should be called whenever a track/signal edit occurs at pos.
func (bm *blockManager) MarkDirty(pos types.Pos) {
	bm.dirty[pos] = struct{}{}
}

// ProcessDirty rebuilds blocks for only the region impacted by recent edits.
func (bm *blockManager) ProcessDirty() {
	if len(bm.dirty) == 0 {
		return
	}

	changed := make([]types.Pos, 0, len(bm.dirty))
	for p := range bm.dirty {
		changed = append(changed, p)
	}
	clear(bm.dirty)

	bm.rebuildBlocksAround(changed)
}

// RebuildAll recomputes every block from scratch.
func (bm *blockManager) RebuildAll() {
	for _, tr := range bm.w.Tracks {
		tr.Block = nil
	}

	visited := make(map[types.Pos]bool, len(bm.w.Tracks))

	type item struct {
		pos types.Pos
		tr  *types.Track
	}

	for startPos, startTr := range bm.w.Tracks {
		if visited[startPos] {
			continue
		}

		queue := []item{{pos: startPos, tr: startTr}}
		visited[startPos] = true

		component := make([]*types.Track, 0, 32)

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			if curr.tr.Block != nil {
				continue
			}
			component = append(component, curr.tr)

			for d := types.Dir(types.DirNorth); d <= types.DirWest; d <<= 1 {
				if curr.tr.Direction&d == 0 {
					continue
				}

				np := types.NextPos(curr.pos, d)
				neigh, ok := bm.w.Tracks[np]
				if !ok {
					continue
				}

				if neigh.Direction&types.OppositeDir(d) == 0 {
					continue
				}

				if isSignalBoundary(curr.tr, neigh, d) {
					continue
				}

				if !visited[np] {
					visited[np] = true
					queue = append(queue, item{pos: np, tr: neigh})
				}
			}
		}

		if len(component) == 0 {
			continue
		}

		block := types.NewBlock()
		for _, tr := range component {
			tr.Block = block
		}
	}
}

// isSignalBoundary reports whether the BFS should stop crossing from curr to
// neighbour in dir — i.e. whether a signal faces the direction of travel.
// Blocks are contiguous sections of track between signals; junctions and
// stations do NOT split a block (see types.IsNode for the separate,
// broader notion of a routing decision point used by the routing package).
func isSignalBoundary(curr, neighbour *types.Track, dir types.Dir) bool {
	return types.IsSignalBoundary(curr, neighbour, dir)
}

// rebuildBlocksAround performs a localized rebuild:
// - identifies blocks that might be split/merged due to edits
// - clears ONLY those tracks' Block pointers
// - recomputes connected components (respecting signal boundaries) starting from local seeds
func (bm *blockManager) rebuildBlocksAround(changed []types.Pos) {
	// Helper: is there a track at pos?
	isTrack := func(pos types.Pos) bool {
		if bm.w.TileAt(pos).Type != types.TileTrack {
			return false
		}
		return bm.w.Tracks[pos] != nil
	}

	// 1) Determine impacted block IDs (from changed + 4-neighbors).
	impacted := make(map[uuid.UUID]struct{}, 32)

	addIfHasBlock := func(pos types.Pos) {
		tr := bm.w.Tracks[pos]
		if tr == nil || tr.Block == nil {
			return
		}
		impacted[tr.Block.ID] = struct{}{}
	}

	for _, p := range changed {
		addIfHasBlock(p)
		for d := types.Dir(types.DirNorth); d <= types.DirWest; d <<= 1 {
			addIfHasBlock(types.NextPos(p, d))
		}
	}

	// 2) Clear blocks for tracks in impacted blocks.
	// Also build "seeds": changed + 4-neighbors that are tracks.
	seeds := make([]types.Pos, 0, len(changed)*5)

	for _, p := range changed {
		if isTrack(p) {
			seeds = append(seeds, p)
		}
		for d := types.Dir(types.DirNorth); d <= types.DirWest; d <<= 1 {
			np := types.NextPos(p, d)
			if isTrack(np) {
				seeds = append(seeds, np)
			}
		}
	}

	if len(impacted) > 0 {
		for _, tr := range bm.w.Tracks {
			if tr == nil || tr.Block == nil {
				continue
			}
			if _, ok := impacted[tr.Block.ID]; ok {
				tr.Block = nil
			}
		}
	}

	// 3) Rebuild components starting from seeds (only for nil-block tracks).
	type item struct {
		pos types.Pos
		tr  *types.Track
	}

	seedVisited := make(map[types.Pos]bool, len(seeds))

	for _, start := range seeds {
		if seedVisited[start] {
			continue
		}
		seedVisited[start] = true

		startTr := bm.w.Tracks[start]
		if startTr == nil || startTr.Block != nil {
			continue
		}

		queue := []item{{pos: start, tr: startTr}}
		component := make([]*types.Track, 0, 32)

		// Local visited for this component (prevents revisiting while still nil).
		visited := map[types.Pos]bool{start: true}

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			if curr.tr == nil || curr.tr.Block != nil {
				continue
			}
			component = append(component, curr.tr)

			for d := types.Dir(types.DirNorth); d <= types.DirWest; d <<= 1 {
				if curr.tr.Direction&d == 0 {
					continue
				}

				np := types.NextPos(curr.pos, d)
				neigh := bm.w.Tracks[np]
				if neigh == nil || neigh.Block != nil {
					continue
				}

				if neigh.Direction&types.OppositeDir(d) == 0 {
					continue
				}

				if isSignalBoundary(curr.tr, neigh, d) {
					continue
				}

				if !visited[np] {
					visited[np] = true
					queue = append(queue, item{pos: np, tr: neigh})
				}
			}
		}

		if len(component) == 0 {
			continue
		}

		block := types.NewBlock()
		for _, tr := range component {
			tr.Block = block
		}
	}
}
