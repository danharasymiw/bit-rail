# Pathfinding with Routing Tables

## Blocks

- Tracks are divided into **blocks**
- Each block starts/ends anywhere a train can stop/turn
  - Signals
  - Junctions
  - Stations

- Each block stores:
  - ID
  - Start/end nodes (junctions/signals/stations)
  - Connecting blocks
  - Occupancy/reservation status
  - Travel cost (how much to go from one end to the other)

## Nodes/Junctions
- Nodes represent
  - Junctions
  - Signals
  - Stations
- Each node maintains a **routing table**
  - Maps destination -> next block/direction
  - Built lazily on demand
  - Entry includes:
    - next hop
    - cached cost
    - dirty flag if tracks have changed (these need to be propagated backwards)

## Routing Table Logic
1. Train requests path from current position to destination
2. Node checks routing table entry
  - If the entry exists return the next block/direction
  - If the entry does not exist or is flagged dirty, recalculate
3. Train potentially reserves a certain number of blocks ahead?
4. Move train forward each tick until it encounters another node/junction

## Handling track changes
- If a block or a junction changes
  - Identify nodes/junctions that are connected
  - Mark the entry in those tables as being dirty
    - Propagate the dirty flag backwards?

## Cost
- Each block has a cost
- Somehow apply a cost dynamically if a block is currently reserved or occupied?
