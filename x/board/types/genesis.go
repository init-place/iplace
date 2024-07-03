package types

func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	boards := make(map[uint32]Board)
	for _, board := range gs.Boards {
		if _, ok := boards[board.Id]; ok {
			return ErrDuplicateId
		}
		if err := board.Validate(); err != nil {
			return err
		}

		boards[board.Id] = board
	}

	//sliceIds := make(map[string]bool)
	//for _, slice := range gs.Slices {
	//	if _, ok := sliceIds[slice.GetSliceId()]; ok {
	//		return ErrDuplicateId
	//	}
	//	if board, ok := boards[slice.BoardAddress]; ok {
	//		if err := slice.Validate(board); err != nil {
	//			return err
	//		}
	//	} else {
	//		return ErrUnknownBoardAddress
	//	}
	//
	//	sliceIds[slice.GetSliceId()] = true
	//}

	pixels := make(map[uint32]map[uint32]PixelInfo)
	for _, pixel := range gs.Pixels {
		if boardPixels, ok := pixels[pixel.BoardId]; ok {
			if _, ok := boardPixels[pixel.PixelIndex]; ok {
				return ErrDuplicateId
			}
		} else {
			pixels[pixel.BoardId] = make(map[uint32]PixelInfo)
		}

		if board, ok := boards[pixel.BoardId]; ok {
			if err := pixel.Validate(board); err != nil {
				return err
			}
		} else {
			return ErrUnknownBoardId
		}

		pixels[pixel.BoardId][pixel.PixelIndex] = pixel
	}

	return nil
}
