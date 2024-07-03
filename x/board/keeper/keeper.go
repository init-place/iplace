package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	sdkerrors "cosmossdk.io/errors"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/init-place/iplace/x/board/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Schema collections.Schema
	Params collections.Item[types.Params]
	Boards collections.Map[uint32, types.Board]
	//Slices      collections.Map[string, types.BoardSlice]
	PixelInfos  collections.Map[[]byte, types.PixelInfo]
	Pixels      collections.Map[[]byte, uint32]
	LastBoardId collections.Item[uint32]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		Params:       collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Boards:       collections.NewMap(sb, types.BoardsKey, "boards", collections.Uint32Key, codec.CollValue[types.Board](cdc)),
		//Slices:       collections.NewMap(sb, types.SlicesKey, "slices", collections.StringKey, codec.CollValue[types.BoardSlice](cdc)),
		PixelInfos:  collections.NewMap(sb, types.PixelInfosKey, "pixel_infos", collections.BytesKey, codec.CollValue[types.PixelInfo](cdc)),
		Pixels:      collections.NewMap(sb, types.PixelsKey, "pixels", collections.BytesKey, collections.Uint32Value),
		LastBoardId: collections.NewItem(sb, types.LastBoardIdKey, "last_board_id", collections.Uint32Value),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

//func (k Keeper) generateBoardAddress(ctx context.Context) (sdk.AccAddress, error) {
//	lastBoardId, err := k.LastBoardId.Get(ctx)
//
//	if err != nil {
//		if nerrors.Is(err, collections.ErrNotFound) {
//			lastBoardId = 0
//		} else {
//			return nil, errors.Wrap(err, "failed to generate board id")
//		}
//	}
//
//	myBoardId := lastBoardId + 1
//	err = k.LastBoardId.Set(ctx, myBoardId)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to generate board id")
//	}
//
//	boardIdBytes := make([]byte, 8)
//	binary.BigEndian.PutUint64(boardIdBytes, myBoardId)
//
//	return addresstypes.Module(types.ModuleName, boardIdBytes)[:types.BoardAddressLen], nil
//}

func (k Keeper) nextBoardId(ctx context.Context) (uint32, error) {
	lastBoardId, err := k.LastBoardId.Get(ctx)

	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			lastBoardId = 0
		} else {
			return lastBoardId, sdkerrors.Wrap(err, "failed to get last board id")
		}
	}

	return lastBoardId + 1, nil
}

func (k Keeper) generatePixelStoreKey(boardId uint32, pixelIndex uint32) []byte {
	var sliceId uint8 = 0

	key := make([]byte, 9)

	binary.BigEndian.PutUint32(key, boardId)
	key[4] = sliceId
	binary.BigEndian.PutUint32(key[5:], pixelIndex)

	return key
}
