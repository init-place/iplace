package types

import "cosmossdk.io/errors"

var (
	ErrDuplicateId    = errors.Register(ModuleName, 1, "duplicate id")
	ErrUnknownBoardId = errors.Register(ModuleName, 2, "unknown board_id")
)
