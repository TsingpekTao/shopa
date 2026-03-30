package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/chat-svc/api/v1"
	"github.com/TsingpekTao/shopa/chat-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerChatServiceServer
	v1.UnimplementedSellerChatServiceServer
	v1.UnimplementedInternalChatServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerChatServiceServer(s.Server, ctrl)
	v1.RegisterSellerChatServiceServer(s.Server, ctrl)
	v1.RegisterInternalChatServiceServer(s.Server, ctrl)
}

func (*Controller) CreateOrGetConversation(ctx context.Context, req *v1.CreateOrGetConversationReq) (res *v1.CreateOrGetConversationRes, err error) {
	return service.Chat().CreateOrGetConversation(ctx, req)
}

func (*Controller) SendMessage(ctx context.Context, req *v1.SendMessageReq) (res *v1.SendMessageRes, err error) {
	return service.Chat().SendMessage(ctx, req)
}

func (*Controller) ListMyConversations(ctx context.Context, req *v1.ListMyConversationsReq) (res *v1.ListMyConversationsRes, err error) {
	return service.Chat().ListMyConversations(ctx, req)
}

func (*Controller) ListMessages(ctx context.Context, req *v1.ListMessagesReq) (res *v1.ListMessagesRes, err error) {
	return service.Chat().ListMessages(ctx, req)
}

func (*Controller) MarkConversationRead(ctx context.Context, req *v1.MarkConversationReadReq) (res *v1.MarkConversationReadRes, err error) {
	return service.Chat().MarkConversationRead(ctx, req)
}

func (*Controller) GetUnreadSummary(ctx context.Context, req *v1.GetUnreadSummaryReq) (res *v1.GetUnreadSummaryRes, err error) {
	return service.Chat().GetUnreadSummary(ctx, req)
}

func (*Controller) ListShopConversations(ctx context.Context, req *v1.ListShopConversationsReq) (res *v1.ListShopConversationsRes, err error) {
	return service.Chat().ListShopConversations(ctx, req)
}

func (*Controller) SendMessageAsSeller(ctx context.Context, req *v1.SendMessageAsSellerReq) (res *v1.SendMessageAsSellerRes, err error) {
	return service.Chat().SendMessageAsSeller(ctx, req)
}

func (*Controller) MarkConversationReadAsSeller(ctx context.Context, req *v1.MarkConversationReadAsSellerReq) (res *v1.MarkConversationReadAsSellerRes, err error) {
	return service.Chat().MarkConversationReadAsSeller(ctx, req)
}

func (*Controller) PublishSystemNotice(ctx context.Context, req *v1.PublishSystemNoticeReq) (res *v1.PublishSystemNoticeRes, err error) {
	return service.Chat().PublishSystemNotice(ctx, req)
}

func (*Controller) GetConversationSnapshot(ctx context.Context, req *v1.GetConversationSnapshotReq) (res *v1.GetConversationSnapshotRes, err error) {
	return service.Chat().GetConversationSnapshot(ctx, req)
}
