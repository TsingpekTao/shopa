package review

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/review-svc/api/v1"
	"github.com/TsingpekTao/shopa/review-svc/internal/dao"
	"github.com/TsingpekTao/shopa/review-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/review-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/review-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
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
)

type sReview struct{}

func New() *sReview {
	return &sReview{}
}

func init() {
	service.RegisterReview(New())
}

func (s *sReview) CreateReview(ctx context.Context, req *v1.CreateReviewReq) (*v1.CreateReviewRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetItemNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/item_no are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	mediasJSON, _ := json.Marshal(req.GetMedias())
	reviewNo := generateBizNo("RV")
	_, err = dao.ReviewRecord.Ctx(ctx).Data(do.ReviewRecord{
		ReviewNo:     reviewNo,
		OrderNo:      strings.TrimSpace(req.GetOrderNo()),
		SubOrderNo:   strings.TrimSpace(req.GetSubOrderNo()),
		ItemNo:       strings.TrimSpace(req.GetItemNo()),
		UserId:       userID,
		ShopNo:       strings.TrimSpace(req.GetShopNo()),
		SpuNo:        strings.TrimSpace(req.GetSpuNo()),
		SkuNo:        strings.TrimSpace(req.GetSkuNo()),
		Score:        req.GetScore(),
		Content:      strings.TrimSpace(req.GetContent()),
		MediasJson:   string(mediasJSON),
		Anonymous:    boolToTinyInt(req.GetAnonymous()),
		ReviewStatus: int(v1.ReviewStatus_REVIEW_STATUS_PUBLISHED),
		Version:      1,
	}).Insert()
	if err != nil {
		return nil, gerror.Wrap(err, "create review failed")
	}
	row, err := s.getReviewByNo(ctx, reviewNo)
	if err != nil {
		return nil, err
	}
	_, _ = s.rebuildSummary(ctx, row.SpuNo)
	return &v1.CreateReviewRes{Review: toProtoReview(row)}, nil
}

func (s *sReview) AppendReview(ctx context.Context, req *v1.AppendReviewReq) (*v1.AppendReviewRes, error) {
	if req == nil || strings.TrimSpace(req.GetReviewNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "review_no is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	row, err := s.getReviewByNo(ctx, req.GetReviewNo())
	if err != nil {
		return nil, err
	}
	if row.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "review does not belong to user")
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	appendMediasJSON, _ := json.Marshal(req.GetAppendMedias())
	_, err = dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, row.ReviewNo).Data(do.ReviewRecord{
		AppendContent:    strings.TrimSpace(req.GetAppendContent()),
		AppendMediasJson: string(appendMediasJSON),
		AppendAt:         gtime.Now(),
		Version:          row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "append review failed")
	}
	row, _ = s.getReviewByNo(ctx, row.ReviewNo)
	return &v1.AppendReviewRes{Review: toProtoReview(row)}, nil
}

func (s *sReview) ListMyReviews(ctx context.Context, req *v1.ListMyReviewsReq) (*v1.ListMyReviewsRes, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().UserId, userID).OrderDesc(dao.ReviewRecord.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.ReviewRecord.Columns().Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			statuses = append(statuses, int(status))
		}
		model = model.WhereIn(dao.ReviewRecord.Columns().ReviewStatus, statuses)
	}
	var rows []entity.ReviewRecord
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list my reviews failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	list := make([]*v1.Review, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoReview(&row))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMyReviewsRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sReview) ReplyReview(ctx context.Context, req *v1.ReplyReviewReq) (*v1.ReplyReviewRes, error) {
	if req == nil || strings.TrimSpace(req.GetReviewNo()) == "" || strings.TrimSpace(req.GetSellerReply()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "review_no/seller_reply are required")
	}
	row, err := s.getReviewByNo(ctx, req.GetReviewNo())
	if err != nil {
		return nil, err
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	_, err = dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, row.ReviewNo).Data(do.ReviewRecord{
		SellerReply:   strings.TrimSpace(req.GetSellerReply()),
		SellerReplyAt: gtime.Now(),
		Version:       row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "reply review failed")
	}
	row, _ = s.getReviewByNo(ctx, row.ReviewNo)
	return &v1.ReplyReviewRes{Review: toProtoReview(row)}, nil
}

func (s *sReview) ListSpuReviews(ctx context.Context, req *v1.ListSpuReviewsReq) (*v1.ListSpuReviewsRes, error) {
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	model := dao.ReviewRecord.Ctx(ctx).
		Where(dao.ReviewRecord.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.ReviewRecord.Columns().ReviewStatus, int(v1.ReviewStatus_REVIEW_STATUS_PUBLISHED)).
		Limit(pageSize + 1)
	if cursorID > 0 {
		model = model.WhereLT(dao.ReviewRecord.Columns().Id, cursorID)
	}
	if req.GetWithMediaOnly() {
		model = model.WhereNotIn(dao.ReviewRecord.Columns().MediasJson, []string{"", "[]", "null"})
	}
	switch req.GetSortCode() {
	case v1.ReviewSortCode_REVIEW_SORT_CODE_SCORE_DESC:
		model = model.OrderDesc(dao.ReviewRecord.Columns().Score).OrderDesc(dao.ReviewRecord.Columns().Id)
	case v1.ReviewSortCode_REVIEW_SORT_CODE_SCORE_ASC:
		model = model.OrderAsc(dao.ReviewRecord.Columns().Score).OrderDesc(dao.ReviewRecord.Columns().Id)
	case v1.ReviewSortCode_REVIEW_SORT_CODE_WITH_MEDIA_FIRST:
		model = model.Order("CASE WHEN medias_json IN ('', '[]', 'null') THEN 0 ELSE 1 END DESC").OrderDesc(dao.ReviewRecord.Columns().Id)
	default:
		model = model.OrderDesc(dao.ReviewRecord.Columns().Id)
	}
	var rows []entity.ReviewRecord
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list spu reviews failed")
	}
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}
	list := make([]*v1.Review, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoReview(&row))
	}
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListSpuReviewsRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

func (s *sReview) GetSpuRatingSummary(ctx context.Context, req *v1.GetSpuRatingSummaryReq) (*v1.GetSpuRatingSummaryRes, error) {
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	var row entity.ReviewSpuSummary
	err := dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, req.GetSpuNo()).Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "query review_spu_summary failed")
	}
	if row.Id == 0 {
		row, err = s.rebuildSummary(ctx, req.GetSpuNo())
		if err != nil {
			return nil, err
		}
	}
	return &v1.GetSpuRatingSummaryRes{Summary: toProtoSummary(&row)}, nil
}

func (s *sReview) ModerateReview(ctx context.Context, req *v1.ModerateReviewReq) (*v1.ModerateReviewRes, error) {
	if req == nil || strings.TrimSpace(req.GetReviewNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "review_no is required")
	}
	row, err := s.getReviewByNo(ctx, req.GetReviewNo())
	if err != nil {
		return nil, err
	}
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}
	_, err = dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, row.ReviewNo).Data(do.ReviewRecord{
		ReviewStatus: int(req.GetTargetStatus()),
		Version:      row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "moderate review failed")
	}
	row, _ = s.getReviewByNo(ctx, row.ReviewNo)
	_, _ = s.rebuildSummary(ctx, row.SpuNo)
	return &v1.ModerateReviewRes{Review: toProtoReview(row)}, nil
}

func (s *sReview) RebuildSpuRatingSummary(ctx context.Context, req *v1.RebuildSpuRatingSummaryReq) (*v1.RebuildSpuRatingSummaryRes, error) {
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	row, err := s.rebuildSummary(ctx, req.GetSpuNo())
	if err != nil {
		return nil, err
	}
	return &v1.RebuildSpuRatingSummaryRes{Summary: toProtoSummary(&row)}, nil
}

func (s *sReview) rebuildSummary(ctx context.Context, spuNo string) (entity.ReviewSpuSummary, error) {
	var rows []entity.ReviewRecord
	if err := dao.ReviewRecord.Ctx(ctx).
		Where(dao.ReviewRecord.Columns().SpuNo, spuNo).
		Where(dao.ReviewRecord.Columns().ReviewStatus, int(v1.ReviewStatus_REVIEW_STATUS_PUBLISHED)).
		Scan(&rows); err != nil {
		return entity.ReviewSpuSummary{}, gerror.Wrap(err, "query reviews for summary failed")
	}
	summary := entity.ReviewSpuSummary{SpuNo: spuNo}
	var scoreSum uint64
	for _, row := range rows {
		summary.TotalReviews++
		switch row.Score {
		case 1:
			summary.Score1Count++
		case 2:
			summary.Score2Count++
		case 3:
			summary.Score3Count++
		case 4:
			summary.Score4Count++
		case 5:
			summary.Score5Count++
		}
		scoreSum += uint64(row.Score)
	}
	if summary.TotalReviews > 0 {
		summary.AvgScoreX100 = scoreSum * 100 / summary.TotalReviews
		positive := summary.Score4Count + summary.Score5Count
		summary.PositiveRateX100 = positive * 10000 / summary.TotalReviews
	}

	var existed entity.ReviewSpuSummary
	if err := dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, spuNo).Scan(&existed); err != nil {
		return entity.ReviewSpuSummary{}, gerror.Wrap(err, "query existed summary failed")
	}
	if existed.Id == 0 {
		_, err := dao.ReviewSpuSummary.Ctx(ctx).Data(do.ReviewSpuSummary{
			SpuNo:            spuNo,
			TotalReviews:     summary.TotalReviews,
			Score1Count:      summary.Score1Count,
			Score2Count:      summary.Score2Count,
			Score3Count:      summary.Score3Count,
			Score4Count:      summary.Score4Count,
			Score5Count:      summary.Score5Count,
			AvgScoreX100:     summary.AvgScoreX100,
			PositiveRateX100: summary.PositiveRateX100,
		}).Insert()
		if err != nil {
			return entity.ReviewSpuSummary{}, gerror.Wrap(err, "insert summary failed")
		}
	} else {
		_, err := dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, spuNo).Data(do.ReviewSpuSummary{
			TotalReviews:     summary.TotalReviews,
			Score1Count:      summary.Score1Count,
			Score2Count:      summary.Score2Count,
			Score3Count:      summary.Score3Count,
			Score4Count:      summary.Score4Count,
			Score5Count:      summary.Score5Count,
			AvgScoreX100:     summary.AvgScoreX100,
			PositiveRateX100: summary.PositiveRateX100,
		}).Update()
		if err != nil {
			return entity.ReviewSpuSummary{}, gerror.Wrap(err, "update summary failed")
		}
	}
	_ = dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, spuNo).Scan(&summary)
	return summary, nil
}

func (s *sReview) getReviewByNo(ctx context.Context, reviewNo string) (*entity.ReviewRecord, error) {
	var row entity.ReviewRecord
	if err := dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, reviewNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query review by no failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "review not found")
	}
	return &row, nil
}

func toProtoReview(row *entity.ReviewRecord) *v1.Review {
	if row == nil {
		return nil
	}
	var medias []*v1.ReviewMedia
	if strings.TrimSpace(row.MediasJson) != "" {
		_ = json.Unmarshal([]byte(row.MediasJson), &medias)
	}
	var appendMedias []*v1.ReviewMedia
	if strings.TrimSpace(row.AppendMediasJson) != "" {
		_ = json.Unmarshal([]byte(row.AppendMediasJson), &appendMedias)
	}
	return &v1.Review{
		ReviewNo:      row.ReviewNo,
		OrderNo:       row.OrderNo,
		SubOrderNo:    row.SubOrderNo,
		ItemNo:        row.ItemNo,
		UserId:        row.UserId,
		ShopNo:        row.ShopNo,
		SpuNo:         row.SpuNo,
		SkuNo:         row.SkuNo,
		Score:         uint32(row.Score),
		Content:       row.Content,
		Medias:        medias,
		Anonymous:     row.Anonymous == 1,
		AppendContent: row.AppendContent,
		AppendMedias:  appendMedias,
		AppendAt:      toProtoTs(row.AppendAt),
		SellerReply:   row.SellerReply,
		SellerReplyAt: toProtoTs(row.SellerReplyAt),
		ReviewStatus:  v1.ReviewStatus(row.ReviewStatus),
		LikeCount:     row.LikeCount,
		Version:       row.Version,
		CreatedAt:     toProtoTs(row.CreatedAt),
		UpdatedAt:     toProtoTs(row.UpdatedAt),
	}
}

func toProtoSummary(row *entity.ReviewSpuSummary) *v1.SpuRatingSummary {
	if row == nil {
		return nil
	}
	return &v1.SpuRatingSummary{
		SpuNo:            row.SpuNo,
		TotalReviews:     row.TotalReviews,
		Score_1Count:     row.Score1Count,
		Score_2Count:     row.Score2Count,
		Score_3Count:     row.Score3Count,
		Score_4Count:     row.Score4Count,
		Score_5Count:     row.Score5Count,
		AvgScoreX100:     row.AvgScoreX100,
		PositiveRateX100: row.PositiveRateX100,
		UpdatedAt:        toProtoTs(row.UpdatedAt),
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

func boolToTinyInt(v bool) int {
	if v {
		return 1
	}
	return 0
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
