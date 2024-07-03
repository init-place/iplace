package board

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	boardv1 "github.com/init-place/iplace/api/iplace/board/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: boardv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "GetAllBoards",
					Use:       "get-all-boards",
					Short:     "Get all boards",
				},
				{
					RpcMethod: "GetBoard",
					Use:       "get-board [id]",
					Short:     "Get board",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				//{
				//	RpcMethod: "GetSlice",
				//	Use:       "get-slice [board_address] [position_x] [position_y]",
				//	Short:     "Get slice",
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{
				//		{ProtoField: "board_address"},
				//		{ProtoField: "position_x"},
				//		{ProtoField: "position_y"},
				//	},
				//},
				{
					RpcMethod: "GetPixelInfo",
					Use:       "get-pixel-info [board_id] [pixel_index]",
					Short:     "Get pixel info",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "board_id"},
						{ProtoField: "pixel_index"},
					},
				},
				{
					RpcMethod: "GetPixels",
					Use:       "get-pixels [board_id] [start] [end]",
					Short:     "Get pixels",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "board_id"},
						{ProtoField: "start"},
						{ProtoField: "end"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: boardv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateBoard",
					Use:       "create-board [name] [size_x] [size_y] [creator]",
					Short:     "create board",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
						{ProtoField: "size_x"},
						{ProtoField: "size_y"},
						{ProtoField: "creator"},
					},
				},
				{
					RpcMethod: "SetPixel",
					Use:       "set-pixel [board_id] [pixel_index] [color] [setter]",
					Short:     "set pixel",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "board_id"},
						//{ProtoField: "slice_id"},
						{ProtoField: "pixel_index"},
						{ProtoField: "color"},
						{ProtoField: "setter"},
					},
				},
			},
		},
	}
}
