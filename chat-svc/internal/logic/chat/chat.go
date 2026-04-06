package chat

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/chat-svc/api/v1"
	"github.com/TsingpekTao/shopa/chat-svc/internal/dao"
	"github.com/TsingpekTao/shopa/chat-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/chat-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/chat-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// defaultPageSize 默认分页大小。
	defaultPageSize = 20
	// maxPageSize 最大分页限制，防止单次加载过多会话。
	maxPageSize = 100

	// readerTypeBuyer 用于区分买家的读取角色。
	readerTypeBuyer = 1
	// readerTypeSeller 用于区分卖家的读取角色。
	readerTypeSeller = 2

	// senderTypeBuyer 买家发送消息时的角色标识。
	senderTypeBuyer = 1
	// senderTypeSeller 卖家发送消息时的角色标识。
	senderTypeSeller = 2
	// senderTypeSystem 系统发送通知时的角色标识。
	senderTypeSystem = 3
)

// sChat 封装聊天业务逻辑，用于 service 层注册。
type sChat struct{}

// New 创建 chat 逻辑实例，目前为无状态实体。
func New() *sChat {
	return &sChat{}
}

func init() {
	// 注册聊天实现到 service 层，供 controller 统一调度。
	service.RegisterChat(New())
}

// CreateOrGetConversation 确保买家和指定店铺的特定场景下存在会话，必要时创建新会话。
func (s *sChat) CreateOrGetConversation(ctx context.Context, req *v1.CreateOrGetConversationReq) (*v1.CreateOrGetConversationRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	// 读取当前用户标识。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// 规范化场景码，空值默认预售。
	sceneCode := strings.ToUpper(strings.TrimSpace(req.GetSceneCode()))
	if sceneCode == "" {
		sceneCode = "PRE_SALE"
	}
	// 尝试根据场景/订单/锚定商品定位已存在的会话。
	conv, err := s.findConversationByScene(ctx, userID, req.GetShopNo(), sceneCode, req.GetOrderNo(), req.GetSubOrderNo(), req.GetAnchorSpuNo(), req.GetAnchorSkuNo())
	if err != nil {
		return nil, err
	}
	if conv != nil {
		// 已有会话直接返回并同步未读计数，以保持唯一性。
		unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeBuyer, userID)
		buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
		return &v1.CreateOrGetConversationRes{
			Conversation: toProtoConversation(conv, conversationPresentation{
				unread:           unread,
				viewerType:       readerTypeBuyer,
				buyerProfile:     loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId],
				shop:             loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)],
				buyerReadOffset:  buyerRead,
				sellerReadOffset: sellerRead,
				lastMessage:      loadLastMessage(ctx, conv.LastMessageNo),
			}),
		}, nil
	}

	// 如没有现有会话则创建新会话记录。
	convNo := generateBizNo("CV")
	_, err = dao.Conversation.Ctx(ctx).Data(do.Conversation{
		ConversationNo:     convNo,
		BuyerId:            userID,
		ShopNo:             strings.TrimSpace(req.GetShopNo()),
		SceneCode:          sceneCode,
		OrderNo:            strings.TrimSpace(req.GetOrderNo()),
		SubOrderNo:         strings.TrimSpace(req.GetSubOrderNo()),
		AnchorSpuNo:        strings.TrimSpace(req.GetAnchorSpuNo()),
		AnchorSkuNo:        strings.TrimSpace(req.GetAnchorSkuNo()),
		ConversationStatus: int(v1.ConversationStatus_CONVERSATION_STATUS_ACTIVE),
		Version:            1,
	}).Insert()
	if err != nil {
		// 可能因并发插入导致重复，重查已有会话避免冲突错误。
		conv, qErr := s.findConversationByScene(ctx, userID, req.GetShopNo(), sceneCode, req.GetOrderNo(), req.GetSubOrderNo(), req.GetAnchorSpuNo(), req.GetAnchorSkuNo())
		if qErr == nil && conv != nil {
			unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeBuyer, userID)
			buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
			return &v1.CreateOrGetConversationRes{
				Conversation: toProtoConversation(conv, conversationPresentation{
					unread:           unread,
					viewerType:       readerTypeBuyer,
					buyerProfile:     loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId],
					shop:             loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)],
					buyerReadOffset:  buyerRead,
					sellerReadOffset: sellerRead,
					lastMessage:      loadLastMessage(ctx, conv.LastMessageNo),
				}),
			}, nil
		}
		return nil, gerror.Wrap(err, "create conversation failed")
	}
	conv, err = s.getConversationByNo(ctx, convNo)
	if err != nil {
		return nil, err
	}
	// 创建成功后返回事后查询的会话实体，初始未读为 0。
	buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
	return &v1.CreateOrGetConversationRes{
		Conversation: toProtoConversation(conv, conversationPresentation{
			unread:           0,
			viewerType:       readerTypeBuyer,
			buyerProfile:     loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId],
			shop:             loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)],
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
			lastMessage:      loadLastMessage(ctx, conv.LastMessageNo),
		}),
	}, nil
}

// SendMessage 提供买家向会话发送消息的入口，包含幂等处理与未读刷新。
func (s *sChat) SendMessage(ctx context.Context, req *v1.SendMessageReq) (*v1.SendMessageRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" || strings.TrimSpace(req.GetContentText()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no/content_text are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	msg, conv, replay, err := s.sendMessage(ctx, req.GetConversationNo(), senderTypeBuyer, userID, req.GetClientMessageNo(), req.GetMessageType(), req.GetContentText(), req.GetMediaAssetId(), req.GetExtJson())
	if err != nil {
		return nil, err
	}
	unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeBuyer, userID)
	buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
	buyerProfile := loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId]
	shop := loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)]
	return &v1.SendMessageRes{
		Conversation: toProtoConversation(conv, conversationPresentation{
			unread:           unread,
			viewerType:       readerTypeBuyer,
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
			lastMessage:      msg,
		}),
		Message: toProtoMessage(msg, messagePresentation{
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
		}),
		IdempotentReplay: replay,
	}, nil
}

// ListMyConversations 列出当前买家的会话，支持 seek 分页。
func (s *sChat) ListMyConversations(ctx context.Context, req *v1.ListMyConversationsReq) (*v1.ListMyConversationsRes, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.Conversation.Ctx(ctx).Where(dao.Conversation.Columns().BuyerId, userID).OrderDesc(dao.Conversation.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.Conversation.Columns().Id, cursorID)
	}
	var rows []entity.Conversation
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query conversations failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	shopSummaries := loadShopSummaries(ctx, conversationShopNos(rows))
	buyerProfiles := loadBuyerProfiles(ctx, conversationBuyerIDs(rows))
	list := make([]*v1.Conversation, 0, len(rows))
	for _, row := range rows {
		unread, _ := s.countUnread(ctx, row.ConversationNo, readerTypeBuyer, userID)
		buyerRead, sellerRead := loadReadOffsets(ctx, row.ConversationNo)
		list = append(list, toProtoConversation(&row, conversationPresentation{
			unread:           unread,
			viewerType:       readerTypeBuyer,
			buyerProfile:     buyerProfiles[row.BuyerId],
			shop:             shopSummaries[strings.TrimSpace(row.ShopNo)],
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
			lastMessage:      loadLastMessage(ctx, row.LastMessageNo),
		}))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMyConversationsRes{List: list, HasMore: hasMore, NextCursor: next}, nil
}

// ListMessages 拉取指定会话的消息，返回倒序分页结果。
func (s *sChat) ListMessages(ctx context.Context, req *v1.ListMessagesReq) (*v1.ListMessagesRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	if err != nil {
		return nil, err
	}

	shopNo := shopNoFromContext(ctx)
	if shopNo != "" {
		if !strings.EqualFold(strings.TrimSpace(conv.ShopNo), strings.TrimSpace(shopNo)) {
			return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation not belongs to shop")
		}
	} else {
		userID, userErr := userIDFromContext(ctx)
		if userErr != nil {
			return nil, userErr
		}
		if conv.BuyerId != userID {
			return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation not belongs to buyer")
		}
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().ConversationNo, req.GetConversationNo()).OrderDesc(dao.ChatMessage.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.ChatMessage.Columns().Id, cursorID)
	}
	var rows []entity.ChatMessage
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query messages failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	buyerProfile := loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId]
	shop := loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)]
	buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
	list := make([]*v1.ChatMessage, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		list = append(list, toProtoMessage(&rows[i], messagePresentation{
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
		}))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMessagesRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

// MarkConversationRead 将指定消息及之前的消息标记为买家已读。
func (s *sChat) MarkConversationRead(ctx context.Context, req *v1.MarkConversationReadReq) (*v1.MarkConversationReadRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	if err != nil {
		return nil, err
	}
	if conv.BuyerId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation not belongs to buyer")
	}
	readNo, _, err := s.markRead(ctx, req.GetConversationNo(), readerTypeBuyer, userID, req.GetReadToMessageNo())
	if err != nil {
		return nil, err
	}
	return &v1.MarkConversationReadRes{ConversationNo: req.GetConversationNo(), ReadToMessageNo: readNo}, nil
}

// GetUnreadSummary 汇总买家未读会话数量及总未读条数。
func (s *sChat) GetUnreadSummary(ctx context.Context, req *v1.GetUnreadSummaryReq) (*v1.GetUnreadSummaryRes, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var rows []entity.Conversation
	if err = dao.Conversation.Ctx(ctx).Where(dao.Conversation.Columns().BuyerId, userID).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query conversations failed")
	}
	var totalConv uint64
	var totalMsg uint64
	for _, row := range rows {
		unread, _ := s.countUnread(ctx, row.ConversationNo, readerTypeBuyer, userID)
		if unread > 0 {
			totalConv++
			totalMsg += uint64(unread)
		}
	}
	return &v1.GetUnreadSummaryRes{
		TotalUnreadConversations: totalConv,
		TotalUnreadMessages:      totalMsg,
	}, nil
}

// ListShopConversations 为卖家展示其店铺下的会话列表，支持 seek 分页。
func (s *sChat) ListShopConversations(ctx context.Context, req *v1.ListShopConversationsReq) (*v1.ListShopConversationsRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.Conversation.Ctx(ctx).Where(dao.Conversation.Columns().ShopNo, req.GetShopNo()).OrderDesc(dao.Conversation.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.Conversation.Columns().Id, cursorID)
	}
	var rows []entity.Conversation
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query shop conversations failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	shopSummaries := loadShopSummaries(ctx, conversationShopNos(rows))
	buyerProfiles := loadBuyerProfiles(ctx, conversationBuyerIDs(rows))
	list := make([]*v1.Conversation, 0, len(rows))
	for _, row := range rows {
		unread, _ := s.countUnread(ctx, row.ConversationNo, readerTypeSeller, 0)
		buyerRead, sellerRead := loadReadOffsets(ctx, row.ConversationNo)
		list = append(list, toProtoConversation(&row, conversationPresentation{
			unread:           unread,
			viewerType:       readerTypeSeller,
			buyerProfile:     buyerProfiles[row.BuyerId],
			shop:             shopSummaries[strings.TrimSpace(row.ShopNo)],
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
			lastMessage:      loadLastMessage(ctx, row.LastMessageNo),
		}))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListShopConversationsRes{List: list, HasMore: hasMore, NextCursor: next}, nil
}

// SendMessageAsSeller 卖家在指定会话中发送消息。
func (s *sChat) SendMessageAsSeller(ctx context.Context, req *v1.SendMessageAsSellerReq) (*v1.SendMessageAsSellerRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" || strings.TrimSpace(req.GetContentText()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no/content_text are required")
	}
	userID, err := userIDFromContextOptional(ctx)
	if err != nil {
		return nil, err
	}
	msg, conv, replay, err := s.sendMessage(ctx, req.GetConversationNo(), senderTypeSeller, userID, req.GetClientMessageNo(), req.GetMessageType(), req.GetContentText(), req.GetMediaAssetId(), req.GetExtJson())
	if err != nil {
		return nil, err
	}
	unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeSeller, userID)
	buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
	buyerProfile := loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId]
	shop := loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)]
	return &v1.SendMessageAsSellerRes{
		Conversation: toProtoConversation(conv, conversationPresentation{
			unread:           unread,
			viewerType:       readerTypeSeller,
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
			lastMessage:      msg,
		}),
		Message: toProtoMessage(msg, messagePresentation{
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
		}),
		IdempotentReplay: replay,
	}, nil
}

// MarkConversationReadAsSeller 卖家标记会话为已读，用于后台处理。
func (s *sChat) MarkConversationReadAsSeller(ctx context.Context, req *v1.MarkConversationReadAsSellerReq) (*v1.MarkConversationReadAsSellerRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	userID, err := userIDFromContextOptional(ctx)
	if err != nil {
		return nil, err
	}
	readNo, _, err := s.markRead(ctx, req.GetConversationNo(), readerTypeSeller, userID, req.GetReadToMessageNo())
	if err != nil {
		return nil, err
	}
	return &v1.MarkConversationReadAsSellerRes{ConversationNo: req.GetConversationNo(), ReadToMessageNo: readNo}, nil
}

// PublishSystemNotice 向会话发送系统通知消息。
func (s *sChat) PublishSystemNotice(ctx context.Context, req *v1.PublishSystemNoticeReq) (*v1.PublishSystemNoticeRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" || strings.TrimSpace(req.GetContentText()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no/content_text are required")
	}
	msg, _, _, err := s.sendMessage(ctx, req.GetConversationNo(), senderTypeSystem, 0, "", v1.MessageType_MESSAGE_TYPE_SYSTEM_NOTICE, req.GetContentText(), 0, req.GetExtJson())
	if err != nil {
		return nil, err
	}
	return &v1.PublishSystemNoticeRes{Message: toProtoMessage(msg, messagePresentation{})}, nil
}

// GetConversationSnapshot 返回会话当前状态及最新一条消息。
func (s *sChat) GetConversationSnapshot(ctx context.Context, req *v1.GetConversationSnapshotReq) (*v1.GetConversationSnapshotRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	if err != nil {
		return nil, err
	}
	var msg entity.ChatMessage
	if conv.LastMessageNo != "" {
		_ = dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().MessageNo, conv.LastMessageNo).Scan(&msg)
	}
	if msg.Id == 0 {
		_ = dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().ConversationNo, conv.ConversationNo).OrderDesc(dao.ChatMessage.Columns().Id).Scan(&msg)
	}
	buyerRead, sellerRead := loadReadOffsets(ctx, conv.ConversationNo)
	buyerProfile := loadBuyerProfiles(ctx, []uint64{conv.BuyerId})[conv.BuyerId]
	shop := loadShopSummaries(ctx, []string{conv.ShopNo})[strings.TrimSpace(conv.ShopNo)]
	return &v1.GetConversationSnapshotRes{
		Conversation: toProtoConversation(conv, conversationPresentation{
			unread:           0,
			viewerType:       readerTypeBuyer,
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
			lastMessage:      &msg,
		}),
		LastMessage: toProtoMessage(&msg, messagePresentation{
			buyerProfile:     buyerProfile,
			shop:             shop,
			buyerReadOffset:  buyerRead,
			sellerReadOffset: sellerRead,
		}),
	}, nil
}

// findConversationByScene 按买家、店铺、订单与锚定商品匹配已有会话。
func (s *sChat) findConversationByScene(ctx context.Context, buyerID uint64, shopNo, sceneCode, orderNo, subOrderNo, anchorSpuNo, anchorSkuNo string) (*entity.Conversation, error) {
	var row entity.Conversation
	err := dao.Conversation.Ctx(ctx).
		Where(dao.Conversation.Columns().BuyerId, buyerID).
		Where(dao.Conversation.Columns().ShopNo, strings.TrimSpace(shopNo)).
		Where(dao.Conversation.Columns().SceneCode, strings.ToUpper(strings.TrimSpace(sceneCode))).
		Where(dao.Conversation.Columns().OrderNo, strings.TrimSpace(orderNo)).
		Where(dao.Conversation.Columns().SubOrderNo, strings.TrimSpace(subOrderNo)).
		Where(dao.Conversation.Columns().AnchorSpuNo, strings.TrimSpace(anchorSpuNo)).
		Where(dao.Conversation.Columns().AnchorSkuNo, strings.TrimSpace(anchorSkuNo)).
		Scan(&row)
	if err != nil {
		if isNoRowsErr(err) {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "query conversation failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

// getConversationByNo 根据 conversation_no 读取会话实体。
func (s *sChat) getConversationByNo(ctx context.Context, conversationNo string) (*entity.Conversation, error) {
	var row entity.Conversation
	if err := dao.Conversation.Ctx(ctx).Where(dao.Conversation.Columns().ConversationNo, conversationNo).Scan(&row); err != nil {
		if isNoRowsErr(err) {
			return nil, gerror.NewCode(gcode.CodeNotFound, "conversation not found")
		}
		return nil, gerror.Wrap(err, "query conversation by no failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "conversation not found")
	}
	return &row, nil
}

// sendMessage 在事务中插入消息并同步会话的最新消息字段，同时支持幂等校验。
func (s *sChat) sendMessage(ctx context.Context, conversationNo string, senderType int, senderUserID uint64, clientMessageNo string, messageType v1.MessageType, contentText string, mediaAssetID uint64, extJSON string) (*entity.ChatMessage, *entity.Conversation, bool, error) {
	conv, err := s.getConversationByNo(ctx, conversationNo)
	if err != nil {
		return nil, nil, false, err
	}
	normalizedExtJSON, err := normalizeExtJSON(extJSON)
	if err != nil {
		return nil, nil, false, err
	}
	if senderType == senderTypeBuyer && senderUserID > 0 && conv.BuyerId != senderUserID {
		return nil, nil, false, gerror.NewCode(gcode.CodeNotAuthorized, "conversation not belongs to buyer")
	}
	// 客户端消息号用于幂等检查，避免重复插入相同内容。
	if strings.TrimSpace(clientMessageNo) != "" {
		var existed entity.ChatMessage
		err = dao.ChatMessage.Ctx(ctx).
			Where(dao.ChatMessage.Columns().ConversationNo, conversationNo).
			Where(dao.ChatMessage.Columns().SenderType, senderType).
			Where(dao.ChatMessage.Columns().SenderUserId, senderUserID).
			Where(dao.ChatMessage.Columns().ClientMessageNo, strings.TrimSpace(clientMessageNo)).
			Scan(&existed)
		if err == nil && existed.Id > 0 {
			unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeBuyer, conv.BuyerId)
			return &existed, toConversationWithUnread(conv, unread), true, nil
		}
		if err != nil {
			if isNoRowsErr(err) {
				err = nil
			} else {
				return nil, nil, false, gerror.Wrap(err, "query message idempotency failed")
			}
		}
	}
	if messageType == v1.MessageType_MESSAGE_TYPE_UNSPECIFIED {
		messageType = v1.MessageType_MESSAGE_TYPE_TEXT
	}
	msgNo := generateBizNo("MSG")
	now := gtime.Now()
	// 通过事务先插入消息，再更新会话的最新消息预览和时间。
	if err = dao.ChatMessage.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model(dao.ChatMessage.Table()).Data(do.ChatMessage{
			MessageNo:       msgNo,
			ConversationNo:  conversationNo,
			SenderType:      senderType,
			SenderUserId:    senderUserID,
			ClientMessageNo: strings.TrimSpace(clientMessageNo),
			MessageType:     int(messageType),
			ContentText:     strings.TrimSpace(contentText),
			MediaAssetId:    mediaAssetID,
			ExtJson:         normalizedExtJSON,
			SentAt:          now,
		}).Insert(); e != nil {
			return gerror.Wrap(e, "insert chat message failed")
		}
		if _, e := tx.Model(dao.Conversation.Table()).
			Where(dao.Conversation.Columns().ConversationNo, conversationNo).
			Data(do.Conversation{
				LastMessageNo:      msgNo,
				LastMessagePreview: strings.TrimSpace(contentText),
				LastMessageAt:      now,
			}).Update(); e != nil {
			return gerror.Wrap(e, "update conversation last message failed")
		}
		return nil
	}); err != nil {
		return nil, nil, false, err
	}
	var msg entity.ChatMessage
	if err = dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().MessageNo, msgNo).Scan(&msg); err != nil {
		return nil, nil, false, gerror.Wrap(err, "query inserted message failed")
	}
	conv.LastMessageNo = msgNo
	conv.LastMessagePreview = strings.TrimSpace(contentText)
	conv.LastMessageAt = now
	return &msg, conv, false, nil
}

// markRead 更新读取偏移，记录当前读者已看到的最新消息。
func (s *sChat) markRead(ctx context.Context, conversationNo string, readerType int, readerUserID uint64, readToMessageNo string) (string, uint64, error) {
	var readMsg entity.ChatMessage
	// 优先使用客户端指定的消息号，否则后面会取最新消息作为阅读范围。
	if strings.TrimSpace(readToMessageNo) != "" {
		if err := dao.ChatMessage.Ctx(ctx).
			Where(dao.ChatMessage.Columns().ConversationNo, conversationNo).
			Where(dao.ChatMessage.Columns().MessageNo, strings.TrimSpace(readToMessageNo)).
			Scan(&readMsg); err != nil {
			if isNoRowsErr(err) {
				readMsg = entity.ChatMessage{}
			} else {
				return "", 0, gerror.Wrap(err, "query read message failed")
			}
		}
	}
	if readMsg.Id == 0 {
		if err := dao.ChatMessage.Ctx(ctx).
			Where(dao.ChatMessage.Columns().ConversationNo, conversationNo).
			OrderDesc(dao.ChatMessage.Columns().Id).
			Scan(&readMsg); err != nil {
			if isNoRowsErr(err) {
				return "", 0, nil
			}
			return "", 0, gerror.Wrap(err, "query latest message failed")
		}
	}
	if readMsg.Id == 0 {
		return "", 0, nil
	}
	var offset entity.ChatReadOffset
	// 读取读者当前的已读偏移，后续用于判断是否需更新。
	if err := dao.ChatReadOffset.Ctx(ctx).
		Where(dao.ChatReadOffset.Columns().ConversationNo, conversationNo).
		Where(dao.ChatReadOffset.Columns().ReaderType, readerType).
		Where(dao.ChatReadOffset.Columns().ReaderUserId, readerUserID).
		Scan(&offset); err != nil {
		if isNoRowsErr(err) {
			offset = entity.ChatReadOffset{}
		} else {
			return "", 0, gerror.Wrap(err, "query read offset failed")
		}
	}
	if offset.Id == 0 {
		_, err := dao.ChatReadOffset.Ctx(ctx).Data(do.ChatReadOffset{
			ConversationNo:  conversationNo,
			ReaderType:      readerType,
			ReaderUserId:    readerUserID,
			ReadToMessageNo: readMsg.MessageNo,
			ReadToMessageId: readMsg.Id,
			ReadAt:          gtime.Now(),
		}).Insert()
		if err != nil {
			return "", 0, gerror.Wrap(err, "insert read offset failed")
		}
		return readMsg.MessageNo, readMsg.Id, nil
	}
	if readMsg.Id <= offset.ReadToMessageId {
		return offset.ReadToMessageNo, offset.ReadToMessageId, nil
	}
	_, err := dao.ChatReadOffset.Ctx(ctx).
		Where(dao.ChatReadOffset.Columns().Id, offset.Id).
		Data(do.ChatReadOffset{
			ReadToMessageNo: readMsg.MessageNo,
			ReadToMessageId: readMsg.Id,
			ReadAt:          gtime.Now(),
		}).Update()
	if err != nil {
		return "", 0, gerror.Wrap(err, "update read offset failed")
	}
	return readMsg.MessageNo, readMsg.Id, nil
}

// countUnread 统计指定阅读者在会话内未读的消息数量。
func (s *sChat) countUnread(ctx context.Context, conversationNo string, readerType int, readerUserID uint64) (uint32, error) {
	var offset entity.ChatReadOffset
	_ = dao.ChatReadOffset.Ctx(ctx).
		Where(dao.ChatReadOffset.Columns().ConversationNo, conversationNo).
		Where(dao.ChatReadOffset.Columns().ReaderType, readerType).
		Where(dao.ChatReadOffset.Columns().ReaderUserId, readerUserID).
		Scan(&offset)
	model := dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().ConversationNo, conversationNo)
	// 仅统计对方或系统发送的消息，以免包含自己已读的发送记录。
	if readerType == readerTypeBuyer {
		model = model.WhereIn(dao.ChatMessage.Columns().SenderType, []int{senderTypeSeller, senderTypeSystem})
	} else {
		// 卖家视角需过滤掉自己消息，只统计买家与系统发送的记录。
		model = model.WhereIn(dao.ChatMessage.Columns().SenderType, []int{senderTypeBuyer, senderTypeSystem})
	}
	if offset.ReadToMessageId > 0 {
		model = model.WhereGT(dao.ChatMessage.Columns().Id, offset.ReadToMessageId)
	}
	count, err := model.Count()
	if err != nil {
		return 0, gerror.Wrap(err, "count unread failed")
	}
	return uint32(count), nil
}

// toConversationWithUnread 拷贝会话实体并附加未读计数。
func toConversationWithUnread(row *entity.Conversation, unread uint32) *entity.Conversation {
	if row == nil {
		return nil
	}
	copied := *row
	_ = unread
	return &copied
}

// toProtoConversation 将 Conversation 实体转换为 proto 格式。
func toProtoConversation(row *entity.Conversation, meta conversationPresentation) *v1.Conversation {
	if row == nil {
		return nil
	}
	latestReadByPeer := false
	var latestPeerReadAt *timestamppb.Timestamp
	if row.LastMessageNo != "" && meta.lastMessage != nil {
		switch {
		case meta.viewerType == readerTypeBuyer && meta.lastMessage.SenderType == senderTypeBuyer:
			if meta.sellerReadOffset != nil && meta.sellerReadOffset.ReadToMessageId >= meta.lastMessage.Id {
				latestReadByPeer = true
				latestPeerReadAt = toProtoTs(meta.sellerReadOffset.ReadAt)
			}
		case meta.viewerType == readerTypeSeller && meta.lastMessage.SenderType == senderTypeSeller:
			if meta.buyerReadOffset != nil && meta.buyerReadOffset.ReadToMessageId >= meta.lastMessage.Id {
				latestReadByPeer = true
				latestPeerReadAt = toProtoTs(meta.buyerReadOffset.ReadAt)
			}
		}
	}
	return &v1.Conversation{
		ConversationNo:          row.ConversationNo,
		BuyerId:                 row.BuyerId,
		ShopNo:                  row.ShopNo,
		SceneCode:               row.SceneCode,
		OrderNo:                 row.OrderNo,
		SubOrderNo:              row.SubOrderNo,
		AnchorSpuNo:             row.AnchorSpuNo,
		AnchorSkuNo:             row.AnchorSkuNo,
		LastMessagePreview:      row.LastMessagePreview,
		LastMessageAt:           toProtoTs(row.LastMessageAt),
		ConversationStatus:      v1.ConversationStatus(row.ConversationStatus),
		UnreadCount:             meta.unread,
		CreatedAt:               toProtoTs(row.CreatedAt),
		UpdatedAt:               toProtoTs(row.UpdatedAt),
		ShopName:                defaultShopName(meta.shop),
		BuyerDisplayName:        defaultBuyerName(meta.buyerProfile),
		BuyerAvatarUrl:          strings.TrimSpace(meta.buyerProfile.avatarURL),
		ShopAvatarUrl:           strings.TrimSpace(meta.shop.avatarURL),
		BuyerReadToMessageNo:    readToMessageNo(meta.buyerReadOffset),
		SellerReadToMessageNo:   readToMessageNo(meta.sellerReadOffset),
		LatestMessageReadByPeer: latestReadByPeer,
		LatestMessagePeerReadAt: latestPeerReadAt,
	}
}

// toProtoMessage 将 ChatMessage 实体转换为 proto 输出。
func toProtoMessage(row *entity.ChatMessage, meta messagePresentation) *v1.ChatMessage {
	if row == nil || row.Id == 0 {
		return nil
	}
	var (
		senderDisplayName string
		senderAvatarURL   string
	)
	switch row.SenderType {
	case senderTypeBuyer:
		senderDisplayName = defaultBuyerName(meta.buyerProfile)
		senderAvatarURL = strings.TrimSpace(meta.buyerProfile.avatarURL)
	case senderTypeSeller:
		senderDisplayName = defaultShopName(meta.shop)
		senderAvatarURL = strings.TrimSpace(meta.shop.avatarURL)
	case senderTypeSystem:
		senderDisplayName = "系统通知"
	}
	peerRead, peerOffset := peerReadForMessage(row, meta.buyerReadOffset, meta.sellerReadOffset)
	return &v1.ChatMessage{
		MessageNo:         row.MessageNo,
		ConversationNo:    row.ConversationNo,
		SenderType:        v1.SenderType(row.SenderType),
		SenderUserId:      row.SenderUserId,
		MessageType:       v1.MessageType(row.MessageType),
		ClientMessageNo:   row.ClientMessageNo,
		ContentText:       row.ContentText,
		MediaAssetId:      row.MediaAssetId,
		ExtJson:           row.ExtJson,
		SentAt:            toProtoTs(row.SentAt),
		SenderDisplayName: senderDisplayName,
		SenderAvatarUrl:   senderAvatarURL,
		PeerRead:          peerRead,
		PeerReadAt:        readAtTs(peerOffset),
	}
}

// userIDFromContext 尝试从 HTTP 头或 gRPC metadata 中读取必要的 x-user-id。
func readToMessageNo(offset *entity.ChatReadOffset) string {
	if offset == nil {
		return ""
	}
	return offset.ReadToMessageNo
}

func readAtTs(offset *entity.ChatReadOffset) *timestamppb.Timestamp {
	if offset == nil {
		return nil
	}
	return toProtoTs(offset.ReadAt)
}

func conversationBuyerIDs(rows []entity.Conversation) []uint64 {
	seen := make(map[uint64]struct{}, len(rows))
	out := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row.BuyerId == 0 {
			continue
		}
		if _, ok := seen[row.BuyerId]; ok {
			continue
		}
		seen[row.BuyerId] = struct{}{}
		out = append(out, row.BuyerId)
	}
	return out
}

func conversationShopNos(rows []entity.Conversation) []string {
	seen := make(map[string]struct{}, len(rows))
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		shopNo := strings.TrimSpace(row.ShopNo)
		if shopNo == "" {
			continue
		}
		if _, ok := seen[shopNo]; ok {
			continue
		}
		seen[shopNo] = struct{}{}
		out = append(out, shopNo)
	}
	return out
}

func userIDFromContext(ctx context.Context) (uint64, error) {
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil && uid > 0 {
					return uid, nil
				}
			}
		}
	}
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		if val := md.Get(key); val != nil {
			uid := gconv.Uint64(val)
			if uid > 0 {
				return uid, nil
			}
		}
	}
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

// userIDFromContextOptional 读取可选 x-user-id，未提供时返回 0。
func userIDFromContextOptional(ctx context.Context) (uint64, error) {
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil {
					return uid, nil
				}
				return 0, gerror.Wrap(err, "parse x-user-id failed")
			}
		}
	}
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		if val := md.Get(key); val != nil {
			return gconv.Uint64(val), nil
		}
	}
	return 0, nil
}

func shopNoFromContext(ctx context.Context) string {
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-shop-no", "X-Shop-No", "shop_no", "shopNo"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				return raw
			}
			if raw := strings.TrimSpace(r.Get(key).String()); raw != "" {
				return raw
			}
		}
	}
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-shop-no", "shop_no", "shopno"} {
		if val := md.Get(key); val != nil {
			if shopNo := strings.TrimSpace(gconv.String(val)); shopNo != "" {
				return shopNo
			}
		}
	}
	return ""
}

// normalizePageSize 统一规整分页大小，确保在安全范围内。
func normalizePageSize(reqSize int32) int {
	size := int(reqSize)
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return size
}

// parseCursor 将 next_cursor 转换为 uint64，用于 seek 分页。
func parseCursor(cursor string) (uint64, error) {
	cursor = strings.TrimSpace(cursor)
	if cursor == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	return id, nil
}

// toProtoTs 把 gtime.Timestamp 转成 protobuf Timestamp。
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

// generateBizNo 生成业务序号，结合前缀和时间戳。
func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}

func isNoRowsErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "no rows in result set")
}

func normalizeExtJSON(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(trimmed)) {
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "ext_json must be valid json")
	}
	return trimmed, nil
}
