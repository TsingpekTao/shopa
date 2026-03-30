// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/chat-svc/api/v1"
)

type (
	IChat interface {
		CreateOrGetConversation(ctx context.Context, req *v1.CreateOrGetConversationReq) (*v1.CreateOrGetConversationRes, error)
		SendMessage(ctx context.Context, req *v1.SendMessageReq) (*v1.SendMessageRes, error)
		ListMyConversations(ctx context.Context, req *v1.ListMyConversationsReq) (*v1.ListMyConversationsRes, error)
		ListMessages(ctx context.Context, req *v1.ListMessagesReq) (*v1.ListMessagesRes, error)
		MarkConversationRead(ctx context.Context, req *v1.MarkConversationReadReq) (*v1.MarkConversationReadRes, error)
		GetUnreadSummary(ctx context.Context, req *v1.GetUnreadSummaryReq) (*v1.GetUnreadSummaryRes, error)
		ListShopConversations(ctx context.Context, req *v1.ListShopConversationsReq) (*v1.ListShopConversationsRes, error)
		SendMessageAsSeller(ctx context.Context, req *v1.SendMessageAsSellerReq) (*v1.SendMessageAsSellerRes, error)
		MarkConversationReadAsSeller(ctx context.Context, req *v1.MarkConversationReadAsSellerReq) (*v1.MarkConversationReadAsSellerRes, error)
		PublishSystemNotice(ctx context.Context, req *v1.PublishSystemNoticeReq) (*v1.PublishSystemNoticeRes, error)
		GetConversationSnapshot(ctx context.Context, req *v1.GetConversationSnapshotReq) (*v1.GetConversationSnapshotRes, error)
	}
)

var (
	localChat IChat
)

func Chat() IChat {
	if localChat == nil {
		panic("implement not found for interface IChat, forgot register?")
	}
	return localChat
}

func RegisterChat(i IChat) {
	localChat = i
}
