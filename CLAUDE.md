# bit-rail

Terminal-based multiplayer railway builder (Dwarf Fortress-style). Players build track networks, place signals, and watch trains navigate them with block-based collision avoidance.

## Running the project

```bash
make local          # server + client in one process (most common for dev)
make local-debug    # same, with debug logging and block colorization
make server         # headless server only
make client         # client only (expects server on localhost:2977)
make generate       # regenerate all binary serializers (run after changing types)
```

## Package structure

```
cmd/bit-rail/       entry point; -server / -local / (default: client) flags
cmd/generator/      code generator: reads Go AST, emits binary Read*/Write* funcs
engine/             game loop, block manager, WebSocket server
  engine.go         tick loop; handles player messages; owns dirty track updates
  block_manager.go  BFS-based block computation; ProcessDirty() called each tick
  network.go        WebSocket server; per-player read/write goroutines
client/             TUI client (tcell)
  client.go         main loop; camera; chunk loading; message dispatch
  renderer.go       SimpleRenderer; debug mode colors tracks by block ID
  network.go        WebSocket client; read/write loops
trains/             train movement logic
  train.go          Tick(); moveCars(); signal checking; occupancy updates
  train_world_view.go  interface the train uses (decouples from *world.World)
types/              shared value types (no deps on other packages)
  pos.go            Pos{X,Y int}, Dir bitmask, NextPos, OppositeDir
  track.go          Track{Direction, SignalDir, Block *Block}; signal helpers
  tiles.go          TileType enum + Tile struct
  block.go          Block{ID uuid, OccupiedBy uuid}
world/              World struct + chunk system + Perlin generator
  world.go          Tiles [][]*, Tracks map[Pos]*Track, Trains, Occupied map
  generator.go      Perlin noise terrain (elevation + tree overlay)
  test_worlds/      hardcoded worlds for dev/testing (IntersectingLoops is default)
message/            binary protocol types + generated serializers
  messages.go       message type constants + struct definitions
  binary.go         WriteMessage/GetMessageType; writeString/readUUID helpers
  generate.go       go:generate directives (run `make generate` after type changes)
  *_gen.go          generated Read*/Write* functions — DO NOT edit by hand
```

## Key concepts

**Blocks** — contiguous sections of track between signals. Only one train allowed per block at a time. Represented as a `*Block` shared among all `*Track` pointers in the block. The `blockManager` uses BFS to compute connected components; signal tiles act as boundaries. `MarkDirty(pos)` + `ProcessDirty()` keeps rebuilds localized.

**Directions** — `Dir` is a bitmask (`DirNorth=1, DirEast=2, DirSouth=4, DirWest=8`). A track's `Direction` field is a bitmask of the two (or more) directions it connects. Combined constants (`DirNorthSouth`, `DirFourWay`, etc.) are in `types/pos.go`.

**Signals** — a `Track` with `SignalDir != 0` is a signal tile. `SignalDir` is a bitmask of directions the signal faces. A train approaching from the signal direction must check `Block.OccupiedBy` before entering.

**Chunks** — the world is divided into 64×64 tile chunks. The client requests chunks around the camera as it moves. The server sends `ChunksMessage` in response and `WorldUpdateMessage` each tick.

**Message union structs** — both `incomingMessage` and `outgoingMessage` (in engine and client) use tagged-union style: only one pointer field is non-nil per message. Switch on which field is set to dispatch.

## Binary serialization

All wire types have generated `Read<Type>` and `Write<Type>` functions in `message/*_gen.go`. These are produced by `cmd/generator/main.go` which parses Go AST directly. To add a new serializable type:
1. Define the struct in the appropriate package
2. Add a `//go:generate` line to `message/generate.go`
3. Run `make generate`

Struct field tags: `` `binary:"uint16"` `` on an `int` field serializes it as `uint16` on the wire. Pointer fields are written as `bool` presence flag + value.

**Important**: `Block` is serialized per-track (pointer → value copy). The client must call `deduplicateBlocks()` after receiving tracks to restore shared pointers by matching `Block.ID`.

## Architecture notes

- Engine runs a single goroutine tick loop. All game state mutation happens here (no locks on world state). Network I/O uses channels to cross into/out of the tick goroutine.
- `broadcastCh` is read by `broadcastLoop` which fans out to per-player `outgoingCh` channels (non-blocking sends).
- Train movement: `firstCar` determines the front based on `IsReversing`. Cars shift position chain-style in `moveCars`. Occupancy map is updated in `moveCars`.
- Server binds on `-bind` flag (default `:2977`). Client dials `-connect` flag (default `ws://localhost:2977/ws`).

## Planned features (see docs/routing.md)

- Routing tables at junctions/signals/stations for pathfinding
- Stations as destinations; trains navigate autonomously
- Track building/editing by players at runtime
