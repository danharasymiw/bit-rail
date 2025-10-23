package engine

import (
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

type blockManager struct {
	w *world.World
}

func newBlockManager(w *world.World) *blockManager {
	return &blockManager{
		w: w,
	}
}

func (bm *blockManager) calculateBlock(x, y int, track *types.Track) *types.Block {
	maxDistance := 250

	type QueueItem struct {
		track       *types.Track
		pos         world.Pos
		enteredFrom types.Dir
		distance    int
	}

	var queue []QueueItem
	queue = append(queue,
		QueueItem{
			track:       track,
			pos:         world.Pos{X: x, Y: y},
			enteredFrom: types.OppositeDir(track.SignalDir),
			distance:    0,
		})

	var (
		tracksInFlood []*types.Track
		foundBlock    *types.Block
		visited       = map[*types.Track]bool{}
	)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if visited[curr.track] || curr.distance > maxDistance {
			continue
		}
		visited[curr.track] = true

		if curr.track.Block != nil {
			foundBlock = curr.track.Block
			break
		}
		if curr.track.HasSignal && curr.track.SignalDir&curr.enteredFrom != 0 {
			if foundBlock == nil {
				foundBlock = types.NewBlock()
			}
			tracksInFlood = append(tracksInFlood, curr.track)
			continue
		}

		for d := types.Dir(types.DirNorth); d <= types.DirWest; d <<= 1 {
			if curr.track.Direction&d == 0 {
				continue
			}

			nextPos := nextPos(curr.pos, d)
			neighbourTile := bm.w.TileAt(nextPos)
			// TODO: do we really need this check? We can just check if the track is in the map?
			if neighbourTile.Type != types.TileTrack {
				continue
			}
			if neighbour := bm.w.Tracks[nextPos]; neighbour.Direction&types.OppositeDir(d) != 0 {
				queue = append(queue, QueueItem{
					track:       neighbour,
					pos:         nextPos,
					enteredFrom: types.OppositeDir(d),
					distance:    curr.distance + 1,
				})
			}
		}

	}
	for _, track := range tracksInFlood {
		track.Block = foundBlock
	}
	return foundBlock
}
