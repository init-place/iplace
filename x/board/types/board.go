package types

import (
	"fmt"
)

func (board Board) Validate() (err error) {
	//if board.SizeX-1 < board.ReferenceX || board.SizeY-1 < board.ReferenceY {
	//	return errors.New("reference position cannot larger than board size")
	//}

	if len(board.Name) > MaxBoardNameLen {
		return fmt.Errorf("board name cannot exceed %d characters", MaxBoardNameLen)
	}

	return nil
}

//func (slice BoardSlice) GetSliceId() string {
//	return fmt.Sprintf("%s-%d.%d", slice.BoardAddress, slice.PositionX, slice.PositionY)
//}
//
//func (slice *BoardSlice) InitializePixels() {
//	pixels := make([]byte, slice.SizeX*slice.SizeY*colorByteSize)
//
//	for i := 0; i < len(pixels); i++ {
//		pixels[i] = 0
//	}
//
//	slice.Pixels = pixels
//}
//
//func (slice *BoardSlice) SetPixel(x uint32, y uint32, color uint32) {
//	binary.BigEndian.PutUint32(slice.Pixels[y*slice.SizeX+x:], color)
//}
//
//func (slice BoardSlice) ReadPixel(x uint32, y uint32) (color uint32) {
//	return binary.BigEndian.Uint32(slice.Pixels[y*slice.SizeX+x:])
//}
//
//func (slice BoardSlice) Validate(board Board) (err error) {
//	if slice.BoardAddress != board.Address {
//		return errors.New("board address mismatch")
//	}
//
//	if (slice.PositionX < 0 && uint32(-slice.PositionX) > board.ReferenceX) ||
//		(slice.PositionX >= 0 && uint32(slice.PositionX) > board.SizeX-board.ReferenceX-1) {
//		return errors.New("reference position X out of range")
//	}
//
//	if (slice.PositionY < 0 && uint32(-slice.PositionY) > board.ReferenceY) ||
//		(slice.PositionY >= 0 && uint32(slice.PositionY) > board.SizeY-board.ReferenceY-1) {
//		return errors.New("reference position Y out of range")
//	}
//
//	if len(slice.Pixels) != int(slice.SizeX*slice.SizeY*colorByteSize) {
//		return errors.New("invalid number of pixels")
//	}
//
//	return nil
//}

func (pixel PixelInfo) GetPixelId() string {
	return fmt.Sprintf("%d-0-%d", pixel.BoardId, pixel.PixelIndex)
}

func (pixel PixelInfo) Validate(board Board) (err error) {
	//if slice.PositionX != pixel.SlicePositionX || slice.PositionY != pixel.SlicePositionY {
	//	return errors.New("invalid slice position")
	//}
	//
	//if pixel.PixelPositionX > slice.SizeX-1 {
	//	return errors.New("pixel position X out of range")
	//}
	//
	//if pixel.PixelPositionY > slice.SizeY-1 {
	//	return errors.New("pixel position Y out of range")
	//}

	if pixel.PixelIndex > board.SizeX*board.SizeY-1 {
		return fmt.Errorf("pixel index %d is out of bounds", pixel.PixelIndex)
	}

	return nil
}
