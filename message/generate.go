package message

// List of types to generate binary serializers for
// Run `make generate` or `go generate ./message` to regenerate all serializers.

// message package
//go:generate go run ../cmd/generator/main.go -type=ChatMessage -output=chat_message_gen.go
//go:generate go run ../cmd/generator/main.go -type=ChunksMessage -output=chunks_message_gen.go
//go:generate go run ../cmd/generator/main.go -type=GetChunksMessage -output=get_chunks_message_gen.go
//go:generate go run ../cmd/generator/main.go -type=LoginMessage -output=login_message_gen.go
//go:generate go run ../cmd/generator/main.go -type=InitialLoadMessage -output=initial_load_message_gen.go
//go:generate go run ../cmd/generator/main.go -type=WorldUpdateMessage -output=world_update_message_gen.go
//go:generate go run ../cmd/generator/main.go -type=TileUpdate -output=tile_update_gen.go
//go:generate go run ../cmd/generator/main.go -type=TrackUpdate -output=track_update_gen.go

// types package
//go:generate go run ../cmd/generator/main.go -type=TileType -dir=../types -output=tile_type_gen.go
//go:generate go run ../cmd/generator/main.go -type=Dir -dir=../types -output=dir_gen.go
//go:generate go run ../cmd/generator/main.go -type=Tile -dir=../types -output=tile_gen.go
//go:generate go run ../cmd/generator/main.go -type=Track -dir=../types -output=track_gen.go
//go:generate go run ../cmd/generator/main.go -type=Block -dir=../types -output=block_gen.go

// trains package
//go:generate go run ../cmd/generator/main.go -type=Train -dir=../trains -output=train_gen.go
//go:generate go run ../cmd/generator/main.go -type=CarType -dir=../trains -output=car_type_gen.go
//go:generate go run ../cmd/generator/main.go -type=TrainCar -dir=../trains -output=train_car_gen.go

// world package
//go:generate go run ../cmd/generator/main.go -type=Pos -dir=../types -output=pos_gen.go
//go:generate go run ../cmd/generator/main.go -type=Chunk -dir=../world -output=chunk_gen.go
