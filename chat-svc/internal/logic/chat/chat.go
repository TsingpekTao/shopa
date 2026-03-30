package chat

import (
	"context"
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
	defaultPageSize = 20
	maxPageSize     = 100

	readerTypeBuyer  = 1
	readerTypeSeller = 2

	senderTypeBuyer  = 1
	senderTypeSeller = 2
	senderTypeSystem = 3
)

type sChat struct{}

func New() *sChat {
	return &sChat{}
}

func init() {
	service.RegisterChat(New())
}

func (s *sChat) CreateOrGetConversation(ctx context.Context, req *v1.CreateOrGetConversationReq) (*v1.CreateOrGetConversationRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sceneCode := strings.ToUpper(strings.TrimSpace(req.GetSceneCode()))
	if sceneCode == "" {
		sceneCode = "PRE_SALE"
	}
	conv, err := s.findConversationByScene(ctx, userID, req.GetShopNo(), sceneCode, req.GetOrderNo(), req.GetSubOrderNo(), req.GetAnchorSpuNo(), req.GetAnchorSkuNo())
	if err != nil {
		return nil, err
	}
	if conv != nil {
		unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeBuyer, userID)
		return &v1.CreateOrGetConversationRes{Conversation: toProtoConversation(conv, unread)}, nil
	}

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
		conv, qErr := s.findConversationByScene(ctx, userID, req.GetShopNo(), sceneCode, req.GetOrderNo(), req.GetSubOrderNo(), req.GetAnchorSpuNo(), req.GetAnchorSkuNo())
		if qErr == nil && conv != nil {
			unread, _ := s.countUnread(ctx, conv.ConversationNo, readerTypeBuyer, userID)
			return &v1.CreateOrGetConversationRes{Conversation: toProtoConversation(conv, unread)}, nil
		}
		return nil, gerror.Wrap(err, "create conversation failed")
	}
	conv, err = s.getConversationByNo(ctx, convNo)
	if err != nil {
		return nil, err
	}
	return &v1.CreateOrGetConversationRes{Conversation: toProtoConversation(conv, 0)}, nil
}

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
	return &v1.SendMessageRes{
		Conversation:     toProtoConversation(conv, unread),
		Message:          toProtoMessage(msg),
		IdempotentReplay: replay,
	}, nil
}

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
	list := make([]*v1.Conversation, 0, len(rows))
	for _, row := range rows {
		unread, _ := s.countUnread(ctx, row.ConversationNo, readerTypeBuyer, userID)
		list = append(list, toProtoConversation(&row, unread))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMyConversationsRes{List: list, HasMore: hasMore, NextCursor: next}, nil
}

func (s *sChat) ListMessages(ctx context.Context, req *v1.ListMessagesReq) (*v1.ListMessagesRes, error) {
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
	list := make([]*v1.ChatMessage, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		list = append(list, toProtoMessage(&rows[i]))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMessagesRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

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
	list := make([]*v1.Conversation, 0, len(rows))
	for _, row := range rows {
		unread, _ := s.countUnread(ctx, row.ConversationNo, readerTypeSeller, 0)
		list = append(list, toProtoConversation(&row, unread))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListShopConversationsRes{List: list, HasMore: hasMore, NextCursor: next}, nil
}

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
	return &v1.SendMessageAsSellerRes{
		Conversation:     toProtoConversation(conv, unread),
		Message:          toProtoMessage(msg),
		IdempotentReplay: replay,
	}, nil
}

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

func (s *sChat) PublishSystemNotice(ctx context.Context, req *v1.PublishSystemNoticeReq) (*v1.PublishSystemNoticeRes, error) {
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" || strings.TrimSpace(req.GetContentText()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no/content_text are required")
	}
	msg, _, _, err := s.sendMessage(ctx, req.GetConversationNo(), senderTypeSystem, 0, "", v1.MessageType_MESSAGE_TYPE_SYSTEM_NOTICE, req.GetContentText(), 0, req.GetExtJson())
	if err != nil {
		return nil, err
	}
	return &v1.PublishSystemNoticeRes{Message: toProtoMessage(msg)}, nil
}

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
	return &v1.GetConversationSnapshotRes{
		Conversation: toProtoConversation(conv, 0),
		LastMessage:  toProtoMessage(&msg),
	}, nil
}

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
		return nil, gerror.Wrap(err, "query conversation failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

func (s *sChat) getConversationByNo(ctx context.Context, conversationNo string) (*entity.Conversation, error) {
	var row entity.Conversation
	if err := dao.Conversation.Ctx(ctx).Where(dao.Conversation.Columns().ConversationNo, conversationNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query conversation by no failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "conversation not found")
	}
	return &row, nil
}

func (s *sChat) sendMessage(ctx context.Context, conversationNo string, senderType int, senderUserID uint64, clientMessageNo string, messageType v1.MessageType, contentText string, mediaAssetID uint64, extJSON string) (*entity.ChatMessage, *entity.Conversation, bool, error) {
	conv, err := s.getConversationByNo(ctx, conversationNo)
	if err != nil {
		return nil, nil, false, err
	}
	if senderType == senderTypeBuyer && senderUserID > 0 && conv.BuyerId != senderUserID {
		return nil, nil, false, gerror.NewCode(gcode.CodeNotAuthorized, "conversation not belongs to buyer")
	}
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
			return nil, nil, false, gerror.Wrap(err, "query message idempotency failed")
		}
	}
	if messageType == v1.MessageType_MESSAGE_TYPE_UNSPECIFIED {
		messageType = v1.MessageType_MESSAGE_TYPE_TEXT
	}
	msgNo := generateBizNo("MSG")
	now := gtime.Now()
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
			ExtJson:         strings.TrimSpace(extJSON),
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

func (s *sChat) markRead(ctx context.Context, conversationNo string, readerType int, readerUserID uint64, readToMessageNo string) (string, uint64, error) {
	var readMsg entity.ChatMessage
	if strings.TrimSpace(readToMessageNo) != "" {
		if err := dao.ChatMessage.Ctx(ctx).
			Where(dao.ChatMessage.Columns().ConversationNo, conversationNo).
			Where(dao.ChatMessage.Columns().MessageNo, strings.TrimSpace(readToMessageNo)).
			Scan(&readMsg); err != nil {
			return "", 0, gerror.Wrap(err, "query read message failed")
		}
	}
	if readMsg.Id == 0 {
		if err := dao.ChatMessage.Ctx(ctx).
			Where(dao.ChatMessage.Columns().ConversationNo, conversationNo).
			OrderDesc(dao.ChatMessage.Columns().Id).
			Scan(&readMsg); err != nil {
			return "", 0, gerror.Wrap(err, "query latest message failed")
		}
	}
	if readMsg.Id == 0 {
		return "", 0, nil
	}
	var offset entity.ChatReadOffset
	if err := dao.ChatReadOffset.Ctx(ctx).
		Where(dao.ChatReadOffset.Columns().ConversationNo, conversationNo).
		Where(dao.ChatReadOffset.Columns().ReaderType, readerType).
		Where(dao.ChatReadOffset.Columns().ReaderUserId, readerUserID).
		Scan(&offset); err != nil {
		return "", 0, gerror.Wrap(err, "query read offset failed")
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

func (s *sChat) countUnread(ctx context.Context, conversationNo string, readerType int, readerUserID uint64) (uint32, error) {
	var offset entity.ChatReadOffset
	_ = dao.ChatReadOffset.Ctx(ctx).
		Where(dao.ChatReadOffset.Columns().ConversationNo, conversationNo).
		Where(dao.ChatReadOffset.Columns().ReaderType, readerType).
		Where(dao.ChatReadOffset.Columns().ReaderUserId, readerUserID).
		Scan(&offset)
	model := dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().ConversationNo, conversationNo)
	if readerType == readerTypeBuyer {
		model = model.WhereIn(dao.ChatMessage.Columns().SenderType, []int{senderTypeSeller, senderTypeSystem})
	} else {
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

func toConversationWithUnread(row *entity.Conversation, unread uint32) *entity.Conversation {
	if row == nil {
		return nil
	}
	copied := *row
	_ = unread
	return &copied
}

func toProtoConversation(row *entity.Conversation, unread uint32) *v1.Conversation {
	if row == nil {
		return nil
	}
	return &v1.Conversation{
		ConversationNo:     row.ConversationNo,
		BuyerId:            row.BuyerId,
		ShopNo:             row.ShopNo,
		SceneCode:          row.SceneCode,
		OrderNo:            row.OrderNo,
		SubOrderNo:         row.SubOrderNo,
		AnchorSpuNo:        row.AnchorSpuNo,
		AnchorSkuNo:        row.AnchorSkuNo,
		LastMessagePreview: row.LastMessagePreview,
		LastMessageAt:      toProtoTs(row.LastMessageAt),
		ConversationStatus: v1.ConversationStatus(row.ConversationStatus),
		UnreadCount:        unread,
		CreatedAt:          toProtoTs(row.CreatedAt),
		UpdatedAt:          toProtoTs(row.UpdatedAt),
	}
}

func toProtoMessage(row *entity.ChatMessage) *v1.ChatMessage {
	if row == nil || row.Id == 0 {
		return nil
	}
	return &v1.ChatMessage{
		MessageNo:       row.MessageNo,
		ConversationNo:  row.ConversationNo,
		SenderType:      v1.SenderType(row.SenderType),
		SenderUserId:    row.SenderUserId,
		MessageType:     v1.MessageType(row.MessageType),
		ClientMessageNo: row.ClientMessageNo,
		ContentText:     row.ContentText,
		MediaAssetId:    row.MediaAssetId,
		ExtJson:         row.ExtJson,
		SentAt:          toProtoTs(row.SentAt),
	}
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

func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
