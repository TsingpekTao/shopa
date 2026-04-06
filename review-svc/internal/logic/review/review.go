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
	// defaultPageSize 表示未显式传入分页大小时使用的默认值。
	defaultPageSize = 20
	// maxPageSize 表示接口允许的最大分页大小，避免单次查询拖垮数据库。
	maxPageSize = 100
)

// sReview 承载评价领域的逻辑实现。
// 该结构当前无状态，因此使用空结构体即可。
type sReview struct{}

// New 创建评价领域逻辑实例。
func New() *sReview {
	// 当前逻辑对象不持有运行时状态，直接返回空结构体实例。
	return &sReview{}
}

// init 在包初始化阶段把评价逻辑注册到 service 门面。
func init() {
	// controller 最终会通过 service.Review() 获取这里注册的实现。
	service.RegisterReview(New())
}

// CreateReview 创建一条新的商品评价记录。
func (s *sReview) CreateReview(ctx context.Context, req *v1.CreateReviewReq) (*v1.CreateReviewRes, error) {
	// 创建评价时，订单号、子单号和商品项编号都是必填字段。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetItemNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/item_no are required")
	}

	// 从上下文中解析当前登录用户，保证评价主体真实可追溯。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 先把媒体列表编码成 JSON，便于直接落到单列中保存。
	mediasJSON, _ := json.Marshal(req.GetMedias())
	// 为本次评价生成唯一评价号。
	reviewNo := generateBizNo("RV")

	// 将评价主记录写入数据库。
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

	// 回查刚刚写入的评价记录，用于组装完整响应结构。
	row, err := s.getReviewByNo(ctx, reviewNo)
	if err != nil {
		return nil, err
	}

	// 评价创建后异步友好地重建一次 SPU 评分汇总，失败不阻断主链路。
	_, _ = s.rebuildSummary(ctx, row.SpuNo)

	// 返回协议层评价对象。
	return &v1.CreateReviewRes{Review: toProtoReview(row)}, nil
}

// AppendReview 为已有评价追加追评内容。
func (s *sReview) AppendReview(ctx context.Context, req *v1.AppendReviewReq) (*v1.AppendReviewRes, error) {
	// 追评必须明确指定目标评价编号。
	if req == nil || strings.TrimSpace(req.GetReviewNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "review_no is required")
	}

	// 解析当前用户身份，用于校验评价归属。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 读取原始评价记录，后续要校验归属和版本。
	row, err := s.getReviewByNo(ctx, req.GetReviewNo())
	if err != nil {
		return nil, err
	}
	if row.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "review does not belong to user")
	}

	// 如果客户端传入了乐观锁版本，就强校验版本一致。
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}

	// 将追评媒体序列化后一起更新。
	appendMediasJSON, _ := json.Marshal(req.GetAppendMedias())

	// 只更新追评相关字段，并同步把版本号加一。
	_, err = dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, row.ReviewNo).Data(do.ReviewRecord{
		AppendContent:    strings.TrimSpace(req.GetAppendContent()),
		AppendMediasJson: string(appendMediasJSON),
		AppendAt:         gtime.Now(),
		Version:          row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "append review failed")
	}

	// 重新查询最新记录，保证返回值包含最新版本号和时间。
	row, _ = s.getReviewByNo(ctx, row.ReviewNo)
	return &v1.AppendReviewRes{Review: toProtoReview(row)}, nil
}

// ListMyReviews 分页列出当前用户自己的评价记录。
func (s *sReview) ListMyReviews(ctx context.Context, req *v1.ListMyReviewsReq) (*v1.ListMyReviewsRes, error) {
	// 先从上下文中识别调用者用户。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 统一规范分页大小，避免过大页压垮查询。
	pageSize := normalizePageSize(req.GetPageSize())
	// 解析游标，恢复翻页起点。
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}

	// 以当前用户为维度构造基础查询，并多查一条用于判断 has_more。
	model := dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().UserId, userID).OrderDesc(dao.ReviewRecord.Columns().Id).Limit(pageSize + 1)
	if cursorID > 0 {
		// 只有传了游标时，才按主键继续向后翻页。
		model = model.WhereLT(dao.ReviewRecord.Columns().Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		// 将协议枚举转成数据库状态值，用于筛选指定状态评价。
		statuses := make([]int, 0, len(req.GetStatuses()))
		for _, status := range req.GetStatuses() {
			statuses = append(statuses, int(status))
		}
		model = model.WhereIn(dao.ReviewRecord.Columns().ReviewStatus, statuses)
	}

	// 执行查询并加载评价列表。
	var rows []entity.ReviewRecord
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list my reviews failed")
	}

	// 根据是否多查出一条记录来判断是否仍有下一页。
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}

	// 将数据库实体逐条映射成协议对象。
	list := make([]*v1.Review, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoReview(&row))
	}

	// 用最后一条记录的主键作为下一页游标。
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMyReviewsRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

// ReplyReview 允许商家对买家评价进行回复。
func (s *sReview) ReplyReview(ctx context.Context, req *v1.ReplyReviewReq) (*v1.ReplyReviewRes, error) {
	// 回复时必须提供评价编号和回复内容。
	if req == nil || strings.TrimSpace(req.GetReviewNo()) == "" || strings.TrimSpace(req.GetSellerReply()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "review_no/seller_reply are required")
	}

	// 先读取目标评价，后续要在其基础上做版本校验和更新。
	row, err := s.getReviewByNo(ctx, req.GetReviewNo())
	if err != nil {
		return nil, err
	}

	// 如果调用方带了期望版本，则按版本号做并发保护。
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}

	// 更新商家回复内容、回复时间和版本号。
	_, err = dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, row.ReviewNo).Data(do.ReviewRecord{
		SellerReply:   strings.TrimSpace(req.GetSellerReply()),
		SellerReplyAt: gtime.Now(),
		Version:       row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "reply review failed")
	}

	// 回查最新评价记录后返回。
	row, _ = s.getReviewByNo(ctx, row.ReviewNo)
	return &v1.ReplyReviewRes{Review: toProtoReview(row)}, nil
}

// ListSpuReviews 分页列出某个 SPU 下的公开评价。
func (s *sReview) ListSpuReviews(ctx context.Context, req *v1.ListSpuReviewsReq) (*v1.ListSpuReviewsRes, error) {
	// 商品维度查询时，SPU 编号是必填条件。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}

	// 归一化分页参数并解析游标。
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursor(req.GetNextCursor())
	if err != nil {
		return nil, err
	}

	// 基础查询只读取指定 SPU 且已发布的评价。
	model := dao.ReviewRecord.Ctx(ctx).
		Where(dao.ReviewRecord.Columns().SpuNo, req.GetSpuNo()).
		Where(dao.ReviewRecord.Columns().ReviewStatus, int(v1.ReviewStatus_REVIEW_STATUS_PUBLISHED)).
		Limit(pageSize + 1)
	if cursorID > 0 {
		// 如果有游标，则从上一页最后一条之后继续翻页。
		model = model.WhereLT(dao.ReviewRecord.Columns().Id, cursorID)
	}
	if req.GetWithMediaOnly() {
		// 仅带图/视频评价时，排除空媒体字段。
		model = model.WhereNotIn(dao.ReviewRecord.Columns().MediasJson, []string{"", "[]", "null"})
	}

	// 根据请求中的排序码切换不同排序策略。
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

	// 查询评价列表。
	var rows []entity.ReviewRecord
	if err = model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list spu reviews failed")
	}

	// 使用超一条查询法计算是否存在下一页。
	hasMore := false
	if len(rows) > pageSize {
		hasMore = true
		rows = rows[:pageSize]
	}

	// 转换成协议对象列表。
	list := make([]*v1.Review, 0, len(rows))
	for _, row := range rows {
		list = append(list, toProtoReview(&row))
	}

	// 生成下一页游标。
	next := ""
	if hasMore && len(rows) > 0 {
		next = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListSpuReviewsRes{List: list, NextCursor: next, HasMore: hasMore}, nil
}

// GetSpuRatingSummary 获取指定 SPU 的评分汇总信息。
func (s *sReview) GetSpuRatingSummary(ctx context.Context, req *v1.GetSpuRatingSummaryReq) (*v1.GetSpuRatingSummaryRes, error) {
	// SPU 编号为空时无法定位汇总对象。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}

	// 先查汇总表，正常情况下直接命中缓存好的聚合结果。
	var row entity.ReviewSpuSummary
	err := dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, req.GetSpuNo()).Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "query review_spu_summary failed")
	}
	if row.Id == 0 {
		// 如果汇总不存在，则现场重建一次。
		row, err = s.rebuildSummary(ctx, req.GetSpuNo())
		if err != nil {
			return nil, err
		}
	}
	return &v1.GetSpuRatingSummaryRes{Summary: toProtoSummary(&row)}, nil
}

// ModerateReview 供管理端修改评价审核状态。
func (s *sReview) ModerateReview(ctx context.Context, req *v1.ModerateReviewReq) (*v1.ModerateReviewRes, error) {
	// 审核时必须指明评价编号。
	if req == nil || strings.TrimSpace(req.GetReviewNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "review_no is required")
	}

	// 先加载原始评价，供版本校验和更新使用。
	row, err := s.getReviewByNo(ctx, req.GetReviewNo())
	if err != nil {
		return nil, err
	}

	// 使用版本号避免并发审核覆盖。
	if req.GetExpectedVersion() > 0 && row.Version != req.GetExpectedVersion() {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_version mismatch")
	}

	// 更新评价审核状态并提升版本号。
	_, err = dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, row.ReviewNo).Data(do.ReviewRecord{
		ReviewStatus: int(req.GetTargetStatus()),
		Version:      row.Version + 1,
	}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "moderate review failed")
	}

	// 审核状态变化后重新加载评价，并重建汇总数据。
	row, _ = s.getReviewByNo(ctx, row.ReviewNo)
	_, _ = s.rebuildSummary(ctx, row.SpuNo)
	return &v1.ModerateReviewRes{Review: toProtoReview(row)}, nil
}

// RebuildSpuRatingSummary 手动触发某个 SPU 的评分汇总重建。
func (s *sReview) RebuildSpuRatingSummary(ctx context.Context, req *v1.RebuildSpuRatingSummaryReq) (*v1.RebuildSpuRatingSummaryRes, error) {
	// 重建汇总时必须指明商品 SPU。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}

	// 调用内部聚合逻辑重算并返回最新结果。
	row, err := s.rebuildSummary(ctx, req.GetSpuNo())
	if err != nil {
		return nil, err
	}
	return &v1.RebuildSpuRatingSummaryRes{Summary: toProtoSummary(&row)}, nil
}

// rebuildSummary 重算某个 SPU 的评分统计并回写汇总表。
func (s *sReview) rebuildSummary(ctx context.Context, spuNo string) (entity.ReviewSpuSummary, error) {
	// 先拉取该 SPU 下所有已发布评价，作为汇总计算源数据。
	var rows []entity.ReviewRecord
	if err := dao.ReviewRecord.Ctx(ctx).
		Where(dao.ReviewRecord.Columns().SpuNo, spuNo).
		Where(dao.ReviewRecord.Columns().ReviewStatus, int(v1.ReviewStatus_REVIEW_STATUS_PUBLISHED)).
		Scan(&rows); err != nil {
		return entity.ReviewSpuSummary{}, gerror.Wrap(err, "query reviews for summary failed")
	}

	// 初始化汇总对象，并累计评分分布和平均分。
	summary := entity.ReviewSpuSummary{SpuNo: spuNo}
	var scoreSum uint64
	for _, row := range rows {
		// 每遍历一条已发布评价，总评价数加一。
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
		// 累加总分，后面用于计算平均分。
		scoreSum += uint64(row.Score)
	}
	if summary.TotalReviews > 0 {
		// 将平均分保留两位小数的 x100 形式。
		summary.AvgScoreX100 = scoreSum * 100 / summary.TotalReviews
		// 4 星和 5 星评价视为正向评价。
		positive := summary.Score4Count + summary.Score5Count
		summary.PositiveRateX100 = positive * 10000 / summary.TotalReviews
	}

	// 查看汇总表中是否已经存在当前 SPU 的聚合记录。
	var existed entity.ReviewSpuSummary
	if err := dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, spuNo).Scan(&existed); err != nil {
		return entity.ReviewSpuSummary{}, gerror.Wrap(err, "query existed summary failed")
	}
	if existed.Id == 0 {
		// 首次生成汇总时直接插入。
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
		// 已存在汇总时改为更新，保持一个 SPU 只有一条汇总记录。
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

	// 再查一次数据库，返回最终落库后的最新汇总值。
	_ = dao.ReviewSpuSummary.Ctx(ctx).Where(dao.ReviewSpuSummary.Columns().SpuNo, spuNo).Scan(&summary)
	return summary, nil
}

// getReviewByNo 按评价编号查询评价主记录。
func (s *sReview) getReviewByNo(ctx context.Context, reviewNo string) (*entity.ReviewRecord, error) {
	// 通过 review_no 精确查询单条评价。
	var row entity.ReviewRecord
	if err := dao.ReviewRecord.Ctx(ctx).Where(dao.ReviewRecord.Columns().ReviewNo, reviewNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query review by no failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "review not found")
	}
	return &row, nil
}

// toProtoReview 将评价实体转换成协议层对象。
func toProtoReview(row *entity.ReviewRecord) *v1.Review {
	// 空实体直接返回空协议对象。
	if row == nil {
		return nil
	}

	// 解析主评价媒体列表。
	var medias []*v1.ReviewMedia
	if strings.TrimSpace(row.MediasJson) != "" {
		_ = json.Unmarshal([]byte(row.MediasJson), &medias)
	}

	// 解析追评媒体列表。
	var appendMedias []*v1.ReviewMedia
	if strings.TrimSpace(row.AppendMediasJson) != "" {
		_ = json.Unmarshal([]byte(row.AppendMediasJson), &appendMedias)
	}

	// 按字段一一映射成协议对象。
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

// toProtoSummary 将 SPU 评分汇总实体转换成协议对象。
func toProtoSummary(row *entity.ReviewSpuSummary) *v1.SpuRatingSummary {
	// 空实体直接返回 nil。
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

// userIDFromContext 从 HTTP Header 或 gRPC 元数据中解析用户 ID。
func userIDFromContext(ctx context.Context) (uint64, error) {
	// HTTP 请求优先从 Header 中解析用户身份。
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				// 能成功解析成正整数时直接返回。
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil && uid > 0 {
					return uid, nil
				}
			}
		}
	}

	// 如果不是 HTTP 场景，则继续从 gRPC 元数据里找用户信息。
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		if val := md.Get(key); val != nil {
			// gconv 可以兼容不同底层类型的元数据值。
			uid := gconv.Uint64(val)
			if uid > 0 {
				return uid, nil
			}
		}
	}

	// 两条路径都没有拿到用户身份时，直接返回未授权错误。
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

// normalizePageSize 将分页大小约束在服务可接受范围内。
func normalizePageSize(reqSize int32) int {
	// 先转成 int，便于后续参与 ORM 查询。
	size := int(reqSize)
	if size <= 0 {
		// 未传或非法值时回退到默认页大小。
		size = defaultPageSize
	}
	if size > maxPageSize {
		// 超过服务上限时强制裁剪。
		size = maxPageSize
	}
	return size
}

// parseCursor 解析字符串游标到数据库主键值。
func parseCursor(cursor string) (uint64, error) {
	// 先去掉两端空白字符。
	cursor = strings.TrimSpace(cursor)
	if cursor == "" {
		// 空游标表示从第一页开始。
		return 0, nil
	}

	// 将游标解析为无符号整型主键。
	id, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	return id, nil
}

// boolToTinyInt 将布尔值转换为数据库中的 0/1 表示。
func boolToTinyInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// toProtoTs 把 GoFrame 时间对象转换成 protobuf 时间戳。
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

// generateBizNo 生成带前缀的轻量业务流水号。
func generateBizNo(prefix string) string {
	// 用当前时间和纳秒尾段拼装业务号，满足当前服务唯一性需求。
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
