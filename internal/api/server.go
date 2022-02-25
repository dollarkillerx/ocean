package api

import (
	"context"

	"github.com/dollarkillerx/ocean/rpc/ocean"
	"google.golang.org/grpc"
)

type Server struct{}

func (s Server) Ask(ctx context.Context, in *ocean.AskRequest, opts ...grpc.CallOption) (*ocean.AskResponse, error) {

	switch in.Action {

	}

	return &ocean.AskResponse{
		Code: "0",
	}, nil
}
