package chat

import (
	"context"
	"strings"

	chatv1 "github.com/TsingpekTao/shopa/chat-svc/api/chat/v1"
	pb "github.com/TsingpekTao/shopa/chat-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) CreateOrGetConversation(ctx context.Context, req *chatv1.CreateOrGetConversationReq) (*chatv1.CreateOrGetConversationRes, error) {
	return c.chat.CreateOrGetConversation(ctx, &req.CreateOrGetConversationReq)
}

func (c *ControllerV1) CreateOrGetConversationAlias(ctx context.Context, req *chatv1.CreateOrGetConversationAliasReq) (*chatv1.CreateOrGetConversationAliasRes, error) {
	return c.chat.CreateOrGetConversation(ctx, &req.CreateOrGetConversationReq)
}

func (c *ControllerV1) SendMessage(ctx context.Context, req *chatv1.SendMessageReq) (*chatv1.SendMessageRes, error) {
	return c.chat.SendMessage(ctx, &req.SendMessageReq)
}

func (c *ControllerV1) ListMyConversations(ctx context.Context, req *chatv1.ListMyConversationsReq) (*chatv1.ListMyConversationsRes, error) {
	return c.chat.ListMyConversations(ctx, &req.ListMyConversationsReq)
}

func (c *ControllerV1) ListMyConversationsAlias(ctx context.Context, req *chatv1.ListMyConversationsAliasReq) (*chatv1.ListMyConversationsAliasRes, error) {
	return c.chat.ListMyConversations(ctx, &req.ListMyConversationsReq)
}

func (c *ControllerV1) ListMessages(ctx context.Context, req *chatv1.ListMessagesReq) (*chatv1.ListMessagesRes, error) {
	return c.chat.ListMessages(ctx, &pb.ListMessagesReq{
		ConversationNo: req.ConversationNo,
		PageSize:       req.PageSize,
		NextCursor:     req.NextCursor,
	})
}

func (c *ControllerV1) MarkConversationRead(ctx context.Context, req *chatv1.MarkConversationReadReq) (*chatv1.MarkConversationReadRes, error) {
	return c.chat.MarkConversationRead(ctx, &pb.MarkConversationReadReq{
		ConversationNo:  req.ConversationNo,
		ReadToMessageNo: req.ReadToMessageNo,
		IdempotencyKey:  req.IdempotencyKey,
	})
}

func (c *ControllerV1) MarkConversationReadAlias(ctx context.Context, req *chatv1.MarkConversationReadAliasReq) (*chatv1.MarkConversationReadAliasRes, error) {
	return c.chat.MarkConversationRead(ctx, &req.MarkConversationReadReq)
}

func (c *ControllerV1) GetUnreadSummary(ctx context.Context, req *chatv1.GetUnreadSummaryReq) (*chatv1.GetUnreadSummaryRes, error) {
	return c.chat.GetUnreadSummary(ctx, &pb.GetUnreadSummaryReq{})
}

func (c *ControllerV1) GetUnreadSummaryAlias(ctx context.Context, req *chatv1.GetUnreadSummaryAliasReq) (*chatv1.GetUnreadSummaryAliasRes, error) {
	return c.chat.GetUnreadSummary(ctx, &pb.GetUnreadSummaryReq{})
}

func (c *ControllerV1) ListShopConversations(ctx context.Context, req *chatv1.ListShopConversationsReq) (*chatv1.ListShopConversationsRes, error) {
	return c.chat.ListShopConversations(ctx, &req.ListShopConversationsReq)
}

func (c *ControllerV1) ListShopConversationsAlias(ctx context.Context, req *chatv1.ListShopConversationsAliasReq) (*chatv1.ListShopConversationsAliasRes, error) {
	return c.chat.ListShopConversations(ctx, &pb.ListShopConversationsReq{
		ShopNo:     req.ShopNo,
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
	})
}

func (c *ControllerV1) ListMessagesAsSeller(ctx context.Context, req *chatv1.ListMessagesAsSellerReq) (*chatv1.ListMessagesAsSellerRes, error) {
	if r := g.RequestFromCtx(ctx); r != nil && strings.TrimSpace(req.ShopNo) != "" {
		r.Header.Set("X-Shop-No", strings.TrimSpace(req.ShopNo))
	}
	return c.chat.ListMessages(ctx, &pb.ListMessagesReq{
		ConversationNo: req.ConversationNo,
		PageSize:       req.PageSize,
		NextCursor:     req.NextCursor,
	})
}

func (c *ControllerV1) SendMessageAsSeller(ctx context.Context, req *chatv1.SendMessageAsSellerReq) (*chatv1.SendMessageAsSellerRes, error) {
	return c.chat.SendMessageAsSeller(ctx, &req.SendMessageAsSellerReq)
}

func (c *ControllerV1) MarkConversationReadAsSeller(ctx context.Context, req *chatv1.MarkConversationReadAsSellerReq) (*chatv1.MarkConversationReadAsSellerRes, error) {
	return c.chat.MarkConversationReadAsSeller(ctx, &pb.MarkConversationReadAsSellerReq{
		ConversationNo:  req.ConversationNo,
		ReadToMessageNo: req.ReadToMessageNo,
		IdempotencyKey:  req.IdempotencyKey,
	})
}

func (c *ControllerV1) MarkConversationReadAsSellerAlias(ctx context.Context, req *chatv1.MarkConversationReadAsSellerAliasReq) (*chatv1.MarkConversationReadAsSellerAliasRes, error) {
	return c.chat.MarkConversationReadAsSeller(ctx, &req.MarkConversationReadAsSellerReq)
}

func (c *ControllerV1) PublishSystemNotice(ctx context.Context, req *chatv1.PublishSystemNoticeReq) (*chatv1.PublishSystemNoticeRes, error) {
	return c.chat.PublishSystemNotice(ctx, &req.PublishSystemNoticeReq)
}

func (c *ControllerV1) PublishSystemNoticeAlias(ctx context.Context, req *chatv1.PublishSystemNoticeAliasReq) (*chatv1.PublishSystemNoticeAliasRes, error) {
	return c.chat.PublishSystemNotice(ctx, &req.PublishSystemNoticeReq)
}

func (c *ControllerV1) GetConversationSnapshot(ctx context.Context, req *chatv1.GetConversationSnapshotReq) (*chatv1.GetConversationSnapshotRes, error) {
	return c.chat.GetConversationSnapshot(ctx, &pb.GetConversationSnapshotReq{ConversationNo: req.ConversationNo})
}
