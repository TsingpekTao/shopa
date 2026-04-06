package chat

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/TsingpekTao/shopa/chat-svc/internal/dao"
	"github.com/TsingpekTao/shopa/chat-svc/internal/model/entity"
	sellershopv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	userprofilev1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type profileSummary struct {
	displayName string
	avatarURL   string
}

type shopSummary struct {
	name      string
	avatarURL string
}

type conversationPresentation struct {
	unread           uint32
	viewerType       int
	buyerProfile     profileSummary
	shop             shopSummary
	buyerReadOffset  *entity.ChatReadOffset
	sellerReadOffset *entity.ChatReadOffset
	lastMessage      *entity.ChatMessage
}

type messagePresentation struct {
	buyerProfile     profileSummary
	shop             shopSummary
	buyerReadOffset  *entity.ChatReadOffset
	sellerReadOffset *entity.ChatReadOffset
}

var (
	upstreamOnce      sync.Once
	upstreamInitErr   error
	userProfileConn   *grpc.ClientConn
	sellerShopConn    *grpc.ClientConn
	userProfileClient userprofilev1.UserProfileServiceClient
	sellerShopClient  sellershopv1.InternalShopServiceClient
)

func ensureUpstreamClients(ctx context.Context) error {
	upstreamOnce.Do(func() {
		var err error
		if userProfileConn, err = dialUpstreamGRPC(ctx, "upstream.userProfileGrpc", "127.0.0.1:8002"); err != nil {
			upstreamInitErr = err
			return
		}
		if sellerShopConn, err = dialUpstreamGRPC(ctx, "upstream.sellerShopGrpc", "127.0.0.1:50051"); err != nil {
			upstreamInitErr = err
			return
		}
		userProfileClient = userprofilev1.NewUserProfileServiceClient(userProfileConn)
		sellerShopClient = sellershopv1.NewInternalShopServiceClient(sellerShopConn)
	})
	return upstreamInitErr
}

func dialUpstreamGRPC(ctx context.Context, cfgKey, defaultAddr string) (*grpc.ClientConn, error) {
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, cfgKey, defaultAddr).String())
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func loadBuyerProfiles(ctx context.Context, buyerIDs []uint64) map[uint64]profileSummary {
	result := make(map[uint64]profileSummary)
	if len(buyerIDs) == 0 {
		return result
	}
	if err := ensureUpstreamClients(ctx); err != nil {
		return result
	}
	res, err := userProfileClient.BatchGetProfileSummary(ctx, &userprofilev1.BatchGetProfileSummaryReq{UserIds: buyerIDs})
	if err != nil {
		return result
	}
	for _, user := range res.GetUsers() {
		if user == nil || user.GetUserId() == 0 {
			continue
		}
		result[user.GetUserId()] = profileSummary{
			displayName: strings.TrimSpace(user.GetDisplayName()),
			avatarURL:   strings.TrimSpace(user.GetAvatar().GetUrl()),
		}
	}
	return result
}

func loadShopSummaries(ctx context.Context, shopNos []string) map[string]shopSummary {
	result := make(map[string]shopSummary)
	if len(shopNos) == 0 {
		return result
	}
	if err := ensureUpstreamClients(ctx); err != nil {
		return result
	}
	for _, rawShopNo := range shopNos {
		shopNo := strings.TrimSpace(rawShopNo)
		if shopNo == "" {
			continue
		}
		if _, ok := result[shopNo]; ok {
			continue
		}
		res, err := sellerShopClient.GetShopByNo(ctx, &sellershopv1.GetShopByNoReq{ShopNo: shopNo})
		if err != nil || res.GetShop() == nil {
			continue
		}
		shop := res.GetShop()
		name := strings.TrimSpace(shop.GetShopDisplayName())
		if name == "" {
			name = strings.TrimSpace(shop.GetShopName())
		}
		result[shopNo] = shopSummary{name: name}
	}
	return result
}

func loadReadOffsets(ctx context.Context, conversationNo string) (*entity.ChatReadOffset, *entity.ChatReadOffset) {
	var offsets []entity.ChatReadOffset
	if err := dao.ChatReadOffset.Ctx(ctx).
		Where(dao.ChatReadOffset.Columns().ConversationNo, conversationNo).
		WhereIn(dao.ChatReadOffset.Columns().ReaderType, []int{readerTypeBuyer, readerTypeSeller}).
		Scan(&offsets); err != nil {
		return nil, nil
	}
	var (
		buyerOffset  *entity.ChatReadOffset
		sellerOffset *entity.ChatReadOffset
	)
	for i := range offsets {
		item := offsets[i]
		switch item.ReaderType {
		case readerTypeBuyer:
			buyerOffset = &item
		case readerTypeSeller:
			sellerOffset = &item
		}
	}
	return buyerOffset, sellerOffset
}

func loadLastMessage(ctx context.Context, messageNo string) *entity.ChatMessage {
	messageNo = strings.TrimSpace(messageNo)
	if messageNo == "" {
		return nil
	}
	var row entity.ChatMessage
	if err := dao.ChatMessage.Ctx(ctx).Where(dao.ChatMessage.Columns().MessageNo, messageNo).Scan(&row); err != nil {
		return nil
	}
	if row.Id == 0 {
		return nil
	}
	return &row
}

func peerReadForMessage(row *entity.ChatMessage, buyerOffset, sellerOffset *entity.ChatReadOffset) (bool, *entity.ChatReadOffset) {
	if row == nil {
		return false, nil
	}
	switch row.SenderType {
	case senderTypeBuyer:
		if sellerOffset != nil && sellerOffset.ReadToMessageId >= row.Id {
			return true, sellerOffset
		}
	case senderTypeSeller:
		if buyerOffset != nil && buyerOffset.ReadToMessageId >= row.Id {
			return true, buyerOffset
		}
	}
	return false, nil
}

func defaultBuyerName(profile profileSummary) string {
	if strings.TrimSpace(profile.displayName) != "" {
		return strings.TrimSpace(profile.displayName)
	}
	return "买家"
}

func defaultShopName(shop shopSummary) string {
	if strings.TrimSpace(shop.name) != "" {
		return strings.TrimSpace(shop.name)
	}
	return "店铺客服"
}

