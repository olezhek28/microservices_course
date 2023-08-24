package chat_v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/olezhek28/microservices_course/week_8/chat/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	chatChan <- req.GetMessage()

	return &emptypb.Empty{}, nil
}
