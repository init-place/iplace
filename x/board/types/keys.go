package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "board"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	MaxBoardNameLen = 256
	//MaxSliceXSize   = 256
	//MaxSliceYSize   = 256
)

var (
	ParamsKey = collections.NewPrefix("params")
	BoardsKey = collections.NewPrefix("boards/")
	//SlicesKey = collections.NewPrefix("slices/")
	PixelInfosKey = collections.NewPrefix("pixel_infos/")
	PixelsKey     = collections.NewPrefix(0xa0)

	LastBoardIdKey = collections.NewPrefix("sequences/board_id")
	//LastSliceIdKey = collections.NewPrefix("sequences/slice_id")
)
