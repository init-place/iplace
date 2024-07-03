package keeper

import (
	"context"

	"github.com/init-place/iplace/x/board/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}

	for _, board := range data.Boards {
		if err := k.Boards.Set(ctx, board.Id, board); err != nil {
			return err
		}
	}

	//for _, slice := range data.Slices {
	//	if err := k.Slices.Set(ctx, slice.GetSliceId(), slice); err != nil {
	//		return err
	//	}
	//}

	for _, pixel := range data.Pixels {
		key := k.generatePixelStoreKey(pixel.BoardId, pixel.PixelIndex)
		if err := k.PixelInfos.Set(ctx, key, pixel); err != nil {
			return err
		}
		if err := k.Pixels.Set(ctx, key, pixel.Color); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var boards []types.Board
	if err = k.Boards.Walk(ctx, nil, func(boardId uint32, board types.Board) (bool, error) {
		boards = append(boards, board)
		return false, nil
	}); err != nil {
		return nil, err
	}

	//var slices []types.BoardSlice
	//if err = k.Slices.Walk(ctx, nil, func(index string, slice types.BoardSlice) (bool, error) {
	//	slices = append(slices, slice)
	//	return false, nil
	//}); err != nil {
	//	return nil, err
	//}

	var pixels []types.PixelInfo
	if err = k.PixelInfos.Walk(ctx, nil, func(key []byte, pixel types.PixelInfo) (bool, error) {
		pixels = append(pixels, pixel)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Params: params,
		Boards: boards,
		//Slices: slices,
		Pixels: pixels,
	}, nil
}
