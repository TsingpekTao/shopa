package v1

import (
	pb "github.com/TsingpekTao/shopa/chat-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateOrGetConversationReq struct {
	g.Meta `path:"/v1/chat/buyer/conversations:create-or-get" method:"post" tags:"Chat-Buyer" summary:"Create or get conversation"`
	pb.CreateOrGetConversationReq
}

type CreateOrGetConversationRes = pb.CreateOrGetConversationRes

type CreateOrGetConversationAliasReq struct {
	g.Meta `path:"/v1/chat/buyer/conversations:get-or-create" method:"post" tags:"Chat-Buyer" summary:"Create or get conversation (alias)"`
	pb.CreateOrGetConversationReq
}

type CreateOrGetConversationAliasRes = pb.CreateOrGetConversationRes

type SendMessageReq struct {
	g.Meta `path:"/v1/chat/buyer/messages:send" method:"post" tags:"Chat-Buyer" summary:"Buyer send message"`
	pb.SendMessageReq
}

type SendMessageRes = pb.SendMessageRes

type ListMyConversationsReq struct {
	g.Meta `path:"/v1/chat/buyer/conversations:list" method:"get" tags:"Chat-Buyer" summary:"List buyer conversations"`
	pb.ListMyConversationsReq
}

type ListMyConversationsRes = pb.ListMyConversationsRes

type ListMyConversationsAliasReq struct {
	g.Meta `path:"/v1/chat/buyer/conversations" method:"get" tags:"Chat-Buyer" summary:"List buyer conversations (alias)"`
	pb.ListMyConversationsReq
}

type ListMyConversationsAliasRes = pb.ListMyConversationsRes

type ListMessagesReq struct {
	g.Meta         `path:"/v1/chat/buyer/conversations/{conversation_no}/messages" method:"get" tags:"Chat-Buyer" summary:"List conversation messages"`
	ConversationNo string `json:"conversation_no" v:"required#conversation_no is required"`
	PageSize       int32  `json:"page_size"`
	NextCursor     string `json:"next_cursor"`
}

type ListMessagesRes = pb.ListMessagesRes

type MarkConversationReadReq struct {
	g.Meta          `path:"/v1/chat/buyer/conversations/{conversation_no}:mark-read" method:"post" tags:"Chat-Buyer" summary:"Mark conversation read"`
	ConversationNo  string `json:"conversation_no" v:"required#conversation_no is required"`
	ReadToMessageNo string `json:"read_to_message_no"`
	IdempotencyKey  string `json:"idempotency_key"`
}

type MarkConversationReadRes = pb.MarkConversationReadRes

type MarkConversationReadAliasReq struct {
	g.Meta `path:"/v1/chat/buyer/conversations:mark-read" method:"post" tags:"Chat-Buyer" summary:"Mark conversation read (alias)"`
	pb.MarkConversationReadReq
}

type MarkConversationReadAliasRes = pb.MarkConversationReadRes

type GetUnreadSummaryReq struct {
	g.Meta `path:"/v1/chat/buyer/unread:summary" method:"get" tags:"Chat-Buyer" summary:"Get unread summary"`
}

type GetUnreadSummaryRes = pb.GetUnreadSummaryRes

type GetUnreadSummaryAliasReq struct {
	g.Meta `path:"/v1/chat/buyer/unread-summary" method:"get" tags:"Chat-Buyer" summary:"Get unread summary (alias)"`
}

type GetUnreadSummaryAliasRes = pb.GetUnreadSummaryRes

type ListShopConversationsReq struct {
	g.Meta `path:"/v1/chat/seller/conversations:list" method:"get" tags:"Chat-Seller" summary:"List shop conversations"`
	pb.ListShopConversationsReq
}

type ListShopConversationsRes = pb.ListShopConversationsRes

type ListShopConversationsAliasReq struct {
	g.Meta     `path:"/v1/chat/seller/shops/{shop_no}/conversations" method:"get" tags:"Chat-Seller" summary:"List shop conversations (alias)"`
	ShopNo     string `json:"shop_no" v:"required#shop_no is required"`
	PageSize   int32  `json:"page_size"`
	NextCursor string `json:"next_cursor"`
}

type ListShopConversationsAliasRes = pb.ListShopConversationsRes

type ListMessagesAsSellerReq struct {
	g.Meta         `path:"/v1/chat/seller/shops/{shop_no}/conversations/{conversation_no}/messages" method:"get" tags:"Chat-Seller" summary:"List shop conversation messages"`
	ShopNo         string `json:"shop_no" v:"required#shop_no is required"`
	ConversationNo string `json:"conversation_no" v:"required#conversation_no is required"`
	PageSize       int32  `json:"page_size"`
	NextCursor     string `json:"next_cursor"`
}

type ListMessagesAsSellerRes = pb.ListMessagesRes

type SendMessageAsSellerReq struct {
	g.Meta `path:"/v1/chat/seller/messages:send" method:"post" tags:"Chat-Seller" summary:"Seller send message"`
	pb.SendMessageAsSellerReq
}

type SendMessageAsSellerRes = pb.SendMessageAsSellerRes

type MarkConversationReadAsSellerReq struct {
	g.Meta          `path:"/v1/chat/seller/conversations/{conversation_no}:mark-read" method:"post" tags:"Chat-Seller" summary:"Seller mark conversation read"`
	ConversationNo  string `json:"conversation_no" v:"required#conversation_no is required"`
	ReadToMessageNo string `json:"read_to_message_no"`
	IdempotencyKey  string `json:"idempotency_key"`
}

type MarkConversationReadAsSellerRes = pb.MarkConversationReadAsSellerRes

type MarkConversationReadAsSellerAliasReq struct {
	g.Meta `path:"/v1/chat/seller/conversations:mark-read" method:"post" tags:"Chat-Seller" summary:"Seller mark conversation read (alias)"`
	pb.MarkConversationReadAsSellerReq
}

type MarkConversationReadAsSellerAliasRes = pb.MarkConversationReadAsSellerRes

type PublishSystemNoticeReq struct {
	g.Meta `path:"/v1/chat/internal/system-notice:publish" method:"post" tags:"Chat-Internal" summary:"Publish system notice"`
	pb.PublishSystemNoticeReq
}

type PublishSystemNoticeRes = pb.PublishSystemNoticeRes

type PublishSystemNoticeAliasReq struct {
	g.Meta `path:"/v1/chat/internal/system-notices:publish" method:"post" tags:"Chat-Internal" summary:"Publish system notice (alias)"`
	pb.PublishSystemNoticeReq
}

type PublishSystemNoticeAliasRes = pb.PublishSystemNoticeRes

type GetConversationSnapshotReq struct {
	g.Meta         `path:"/v1/chat/internal/conversations/{conversation_no}/snapshot" method:"get" tags:"Chat-Internal" summary:"Get conversation snapshot"`
	ConversationNo string `json:"conversation_no" v:"required#conversation_no is required"`
}

type GetConversationSnapshotRes = pb.GetConversationSnapshotRes
