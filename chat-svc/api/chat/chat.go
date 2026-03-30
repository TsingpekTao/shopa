package chat

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/chat-svc/api/chat/v1"
)

type IChatV1 interface {
	CreateOrGetConversation(ctx context.Context, req *v1.CreateOrGetConversationReq) (res *v1.CreateOrGetConversationRes, err error)
	CreateOrGetConversationAlias(ctx context.Context, req *v1.CreateOrGetConversationAliasReq) (res *v1.CreateOrGetConversationAliasRes, err error)
	SendMessage(ctx context.Context, req *v1.SendMessageReq) (res *v1.SendMessageRes, err error)
	ListMyConversations(ctx context.Context, req *v1.ListMyConversationsReq) (res *v1.ListMyConversationsRes, err error)
	ListMyConversationsAlias(ctx context.Context, req *v1.ListMyConversationsAliasReq) (res *v1.ListMyConversationsAliasRes, err error)
	ListMessages(ctx context.Context, req *v1.ListMessagesReq) (res *v1.ListMessagesRes, err error)
	MarkConversationRead(ctx context.Context, req *v1.MarkConversationReadReq) (res *v1.MarkConversationReadRes, err error)
	MarkConversationReadAlias(ctx context.Context, req *v1.MarkConversationReadAliasReq) (res *v1.MarkConversationReadAliasRes, err error)
	GetUnreadSummary(ctx context.Context, req *v1.GetUnreadSummaryReq) (res *v1.GetUnreadSummaryRes, err error)
	GetUnreadSummaryAlias(ctx context.Context, req *v1.GetUnreadSummaryAliasReq) (res *v1.GetUnreadSummaryAliasRes, err error)
	ListShopConversations(ctx context.Context, req *v1.ListShopConversationsReq) (res *v1.ListShopConversationsRes, err error)
	ListShopConversationsAlias(ctx context.Context, req *v1.ListShopConversationsAliasReq) (res *v1.ListShopConversationsAliasRes, err error)
	SendMessageAsSeller(ctx context.Context, req *v1.SendMessageAsSellerReq) (res *v1.SendMessageAsSellerRes, err error)
	MarkConversationReadAsSeller(ctx context.Context, req *v1.MarkConversationReadAsSellerReq) (res *v1.MarkConversationReadAsSellerRes, err error)
	MarkConversationReadAsSellerAlias(ctx context.Context, req *v1.MarkConversationReadAsSellerAliasReq) (res *v1.MarkConversationReadAsSellerAliasRes, err error)
	PublishSystemNotice(ctx context.Context, req *v1.PublishSystemNoticeReq) (res *v1.PublishSystemNoticeRes, err error)
	PublishSystemNoticeAlias(ctx context.Context, req *v1.PublishSystemNoticeAliasReq) (res *v1.PublishSystemNoticeAliasRes, err error)
	GetConversationSnapshot(ctx context.Context, req *v1.GetConversationSnapshotReq) (res *v1.GetConversationSnapshotRes, err error)
}
