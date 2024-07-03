package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/init-place/iplace/x/board/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queryServer struct {
	k Keeper
}

var _ types.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

func (q queryServer) GetAllBoards(ctx context.Context, request *types.QueryGetAllBoardsRequest) (*types.QueryGetAllBoardsResponse, error) {
	var boards []types.Board
	pageResponse, err := paginate(ctx, q.k.Boards, request.Pagination, func(key uint32) []byte {
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, key)
		return bytes
	}, func(key []byte) uint32 {
		return binary.BigEndian.Uint32(key)
	}, func(key uint32, value types.Board) error {
		boards = append(boards, value)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryGetAllBoardsResponse{
		Boards:     boards,
		Pagination: pageResponse,
	}, nil
}

func (q queryServer) GetBoard(ctx context.Context, request *types.QueryGetBoardRequest) (*types.QueryGetBoardResponse, error) {
	board, err := q.k.Boards.Get(ctx, request.Id)

	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryGetBoardResponse{Board: nil}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetBoardResponse{Board: &board}, nil
}

//func (q queryServer) GetAllSliceInfos(ctx context.Context, request *types.QueryGetAllSliceInfosRequest) (*types.QueryGetAllSliceInfosResponse, error) {
//	var slices []types.BoardSlice
//	pageResponse, err := paginate(ctx, q.k.Slices, request.Pagination, func(key string, value types.BoardSlice) error {
//		slices = append(slices, value)
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	if !request.WithPixels {
//		for i, slice := range slices {
//			slice.Pixels = []byte{}
//			slices[i] = slice
//		}
//	}
//
//	return &types.QueryGetAllSliceInfosResponse{
//		Slices:     slices,
//		Pagination: pageResponse,
//	}, nil
//}
//
//func (q queryServer) GetSlice(ctx context.Context, request *types.QueryGetSliceRequest) (*types.QueryGetSliceResponse, error) {
//	slice, err := q.k.Slices.Get(ctx, q.k.generateSliceId(request.BoardAddress, request.PositionX, request.PositionY))
//
//	if err != nil {
//		if errors.Is(err, collections.ErrNotFound) {
//			return &types.QueryGetSliceResponse{Slice: nil}, nil
//		}
//
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &types.QueryGetSliceResponse{Slice: &slice}, nil
//}

func (q queryServer) GetPixelInfo(ctx context.Context, request *types.QueryGetPixelInfoRequest) (*types.QueryGetPixelInfoResponse, error) {
	pixelStoreKey := q.k.generatePixelStoreKey(request.BoardId, request.PixelIndex)
	pixel, err := q.k.PixelInfos.Get(ctx, pixelStoreKey)

	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryGetPixelInfoResponse{Pixel: nil}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetPixelInfoResponse{Pixel: &pixel}, nil
}

func (q queryServer) GetPixels(ctx context.Context, request *types.QueryGetPixelsRequest) (*types.QueryGetPixelsResponse, error) {
	startPixelStoreKey := q.k.generatePixelStoreKey(request.BoardId, request.Start)
	endPixelStoreKey := q.k.generatePixelStoreKey(request.BoardId, request.End)

	iterRange := &collections.Range[[]byte]{}
	iterRange.StartInclusive(startPixelStoreKey)
	iterRange.EndExclusive(endPixelStoreKey)

	pixels := make([]byte, (request.End-request.Start)*4)
	if err := q.k.Pixels.Walk(ctx, iterRange, func(key []byte, value uint32) (stop bool, err error) {
		pixelIndex := binary.BigEndian.Uint32(key[5:])

		binary.BigEndian.PutUint32(pixels[(pixelIndex-request.Start)*4:], value)

		return false, nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetPixelsResponse{Pixels: pixels}, nil
}

func paginate[K, V any](
	ctx context.Context,
	collection collections.Map[K, V],
	pagination *query.PageRequest,
	serializeKey func(key K) []byte,
	deserializeKey func(key []byte) K,
	onResult func(key K, value V) error,
) (*query.PageResponse, error) {
	if pagination.Offset > 0 && pagination.Key != nil {
		return nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	var key K
	var offset uint64 = pagination.Offset
	var limit uint64 = 100
	var countTotal bool = pagination.CountTotal
	var reverse bool = pagination.Reverse

	if len(pagination.Key) != 0 {
		key = deserializeKey(pagination.Key)
	}
	if pagination.Limit > 0 {
		limit = pagination.Limit
	}

	ranger := &collections.Range[K]{}
	ranger.StartInclusive(key)

	if reverse {
		ranger.Descending()
	}

	iterator, err := collection.Iterate(ctx, ranger)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var count uint64
	var nextKey K

	if len(pagination.Key) != 0 {
		for ; iterator.Valid(); iterator.Next() {
			if count == limit {
				nextKey, err = iterator.Key()
				if err != nil {
					return nil, err
				}
				break
			}
			value, err := iterator.Value()
			if err != nil {
				return nil, err
			}
			err = onResult(key, value)
			if err != nil {
				return nil, err
			}

			count++
		}

		return &query.PageResponse{
			NextKey: serializeKey(nextKey),
		}, nil
	}

	end := offset + limit

	for ; iterator.Valid(); iterator.Next() {
		count++

		if count <= offset {
			continue
		}

		if count <= end {
			key, err := iterator.Key()
			if err != nil {
				return nil, err
			}
			value, err := iterator.Value()
			if err != nil {
				return nil, err
			}
			err = onResult(key, value)
			if err != nil {
				return nil, err
			}
		} else if count == end+1 {
			nextKey, err = iterator.Key()
			if err != nil {
				return nil, err
			}

			if !countTotal {
				break
			}
		}
	}

	res := &query.PageResponse{
		NextKey: serializeKey(nextKey),
	}

	if countTotal {
		res.Total = count
	}

	return res, nil
}
