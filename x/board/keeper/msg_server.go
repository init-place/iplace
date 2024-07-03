package keeper

import (
	"context"
	sdkerrors "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/init-place/iplace/x/board/types"
)

type msgServer struct {
	keeper Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (m msgServer) CreateBoard(ctx context.Context, msg *types.MsgCreateBoard) (*types.MsgCreateBoardResponse, error) {
	boardId, err := m.keeper.nextBoardId(ctx)

	if err != nil {
		return nil, sdkerrors.Wrap(err, "error generating board address")
	}

	board := types.Board{
		Id:      boardId,
		Name:    msg.Name,
		SizeX:   msg.SizeX,
		SizeY:   msg.SizeY,
		Creator: msg.Creator,
		Admin:   msg.Creator,
	}

	if err := board.Validate(); err != nil {
		return nil, err
	}

	if err := m.keeper.Boards.Set(ctx, board.Id, board); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateBoard,
		sdk.NewAttribute(types.AttributeKeyBoardId, fmt.Sprintf("%d", board.Id)),
	))

	return &types.MsgCreateBoardResponse{Board: &board}, nil
}

func (m msgServer) SetPixel(ctx context.Context, msg *types.MsgSetPixel) (*types.MsgSetPixelResponse, error) {
	board, err := m.keeper.Boards.Get(ctx, msg.BoardId)
	if err != nil {
		return nil, err
	}

	pixelStoreKey := m.keeper.generatePixelStoreKey(msg.BoardId, msg.PixelIndex)

	pixel := types.PixelInfo{
		BoardId:    msg.BoardId,
		PixelIndex: msg.PixelIndex,
		Color:      msg.Color,
		Setter:     msg.Setter,
	}

	if err = pixel.Validate(board); err != nil {
		return nil, err
	}

	if err = m.keeper.PixelInfos.Set(ctx, pixelStoreKey, pixel); err != nil {
		return nil, err
	}

	if err = m.keeper.Pixels.Set(ctx, pixelStoreKey, pixel.Color); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSetPixel,
		sdk.NewAttribute(types.AttributeKeyBoardId, fmt.Sprintf("%d", board.Id)),
		sdk.NewAttribute(types.AttributeKeyPixelIndex, fmt.Sprintf("%d", msg.PixelIndex)),
		sdk.NewAttribute(types.AttributeKeyColor, fmt.Sprintf("%d", msg.Color)),
	))

	return &types.MsgSetPixelResponse{}, nil
}
