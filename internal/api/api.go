package api

import (
	"context"
	"errors"

	"github.com/dollarkillerx/ocean/internal/storage"
	"github.com/dollarkillerx/ocean/rpc/ocean"
	"google.golang.org/grpc"
)

type Server struct {
	storage storage.Interface
}

func (s Server) Ask(ctx context.Context, in *ocean.AskRequest, opts ...grpc.CallOption) (*ocean.AskResponse, error) {

	switch in.Action {
	case ocean.Action_ACTION_MIGRATE:
		err := s.migrate(in.Index, in.Payload)
		if err != nil {
			return &ocean.AskResponse{
				Code:    "-1",
				Message: err.Error(),
			}, nil
		}
		return &ocean.AskResponse{
			Code: "0",
		}, nil
	case ocean.Action_ACTION_INSERT:
	case ocean.Action_ACTION_SEARCH:
	case ocean.Action_ACTION_DELETE:
	default:
		return nil, errors.New("illegal parameter")
	}

	return &ocean.AskResponse{
		Code: "0",
	}, nil
}
