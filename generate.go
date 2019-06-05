package scarab

//go:generate protoc --go_out=paths=source_relative:. proto/generic/generic.proto
//go:generate protoc --go_out=paths=source_relative:. proto/state/state.proto
