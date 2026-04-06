package agent

import (
	"context"
	"strings"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

func (s *sAgent) CreateKnowledgeDoc(ctx context.Context, req *agentv1.CreateKnowledgeDocReq) (*agentv1.CreateKnowledgeDocRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || req.Doc == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "doc is required")
	}
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	docNo := generateBizNo("AKD")
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err := dao.AgentKnowledgeDoc.Ctx(ctx).Data(do.AgentKnowledgeDoc{
		KnowledgeDocNo:     docNo,
		Title:              strings.TrimSpace(req.Doc.GetTitle()),
		KnowledgeScopeCode: strings.ToUpper(strings.TrimSpace(req.Doc.GetKnowledgeScopeCode())),
		SourceTypeCode:     strings.ToUpper(strings.TrimSpace(req.Doc.GetSourceTypeCode())),
		SourceId:           strings.TrimSpace(req.Doc.GetSourceId()),
		SourceVersion:      1,
		ShopNo:             strings.TrimSpace(req.Doc.GetShopNo()),
		ContentTypeCode:    strings.ToUpper(strings.TrimSpace(req.Doc.GetContentTypeCode())),
		ContentText:        strings.TrimSpace(req.Doc.GetContentText()),
		AssetIdsJson:       mustJSON(req.Doc.GetAssetIds()),
		TagsJson:           mustJSON(req.Doc.GetTags()),
		StatusCode:         "DRAFT",
		Version:            1,
		CreatedAt:          now,
		UpdatedAt:          now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create knowledge doc failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	doc, err := s.getKnowledgeDocByNo(ctx, docNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.CreateKnowledgeDocRes{Doc: toProtoKnowledgeDoc(doc)}, nil
}

func (s *sAgent) UpdateKnowledgeDoc(ctx context.Context, req *agentv1.UpdateKnowledgeDocReq) (*agentv1.UpdateKnowledgeDocRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetKnowledgeDocNo()) == "" || req.Doc == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "knowledge_doc_no and doc are required")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	doc, err := s.getKnowledgeDocByNo(ctx, req.GetKnowledgeDocNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if req.GetExpectedVersion() > 0 && doc.Version != req.GetExpectedVersion() {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "knowledge doc version conflict")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	data := do.AgentKnowledgeDoc{Version: doc.Version + 1}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	paths := fieldMaskPaths(req.GetUpdateMask().GetPaths())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if len(paths) == 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		paths = []string{"title", "knowledge_scope_code", "source_type_code", "source_id", "shop_no", "content_type_code", "content_text", "asset_ids", "tags"}
	}
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, path := range paths {
		// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
		switch path {
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "title":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.Title = strings.TrimSpace(req.Doc.GetTitle())
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "knowledge_scope_code":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.KnowledgeScopeCode = strings.ToUpper(strings.TrimSpace(req.Doc.GetKnowledgeScopeCode()))
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "source_type_code":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.SourceTypeCode = strings.ToUpper(strings.TrimSpace(req.Doc.GetSourceTypeCode()))
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "source_id":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.SourceId = strings.TrimSpace(req.Doc.GetSourceId())
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "shop_no":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.ShopNo = strings.TrimSpace(req.Doc.GetShopNo())
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "content_type_code":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.ContentTypeCode = strings.ToUpper(strings.TrimSpace(req.Doc.GetContentTypeCode()))
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "content_text":
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			data.ContentText = strings.TrimSpace(req.Doc.GetContentText())
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "asset_ids":
			// 把复杂结构序列化成 JSON 字符串，便于在当前表结构下稳定存储和回显。
			data.AssetIdsJson = mustJSON(req.Doc.GetAssetIds())
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "tags":
			// 把复杂结构序列化成 JSON 字符串，便于在当前表结构下稳定存储和回显。
			data.TagsJson = mustJSON(req.Doc.GetTags())
		}
	}
	// 将本次状态变化写回数据库，保证状态机推进结果能被后续查询立即观察到。
	_, err = dao.AgentKnowledgeDoc.Ctx(ctx).Where(dao.AgentKnowledgeDoc.Columns().KnowledgeDocNo, doc.KnowledgeDocNo).Data(data).Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "update knowledge doc failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	updated, err := s.getKnowledgeDocByNo(ctx, doc.KnowledgeDocNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.UpdateKnowledgeDocRes{Doc: toProtoKnowledgeDoc(updated)}, nil
}

func (s *sAgent) GetKnowledgeDocDetail(ctx context.Context, req *agentv1.GetKnowledgeDocDetailReq) (*agentv1.GetKnowledgeDocDetailRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetKnowledgeDocNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "knowledge_doc_no is required")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	doc, err := s.getKnowledgeDocByNo(ctx, req.GetKnowledgeDocNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.GetKnowledgeDocDetailRes{Doc: toProtoKnowledgeDoc(doc)}, nil
}

func (s *sAgent) ListKnowledgeDocs(ctx context.Context, req *agentv1.ListKnowledgeDocsReq) (*agentv1.ListKnowledgeDocsRes, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	page, pageSize := normalizeOffsetPage(req.GetPage(), req.GetPageSize())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	model := dao.AgentKnowledgeDoc.Ctx(ctx).OrderDesc(dao.AgentKnowledgeDoc.Columns().Id)
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetKnowledgeScopeCode()); v != "" {
		// 统一把编码转换成约定的大写形式，避免状态码和场景码因大小写不一致而失配。
		model = model.Where(dao.AgentKnowledgeDoc.Columns().KnowledgeScopeCode, strings.ToUpper(v))
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetSourceTypeCode()); v != "" {
		// 统一把编码转换成约定的大写形式，避免状态码和场景码因大小写不一致而失配。
		model = model.Where(dao.AgentKnowledgeDoc.Columns().SourceTypeCode, strings.ToUpper(v))
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetShopNo()); v != "" {
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		model = model.Where(dao.AgentKnowledgeDoc.Columns().ShopNo, v)
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetStatusCode()); v != "" {
		// 统一把编码转换成约定的大写形式，避免状态码和场景码因大小写不一致而失配。
		model = model.Where(dao.AgentKnowledgeDoc.Columns().StatusCode, strings.ToUpper(v))
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetKeyword()); v != "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		model = model.WhereLike(dao.AgentKnowledgeDoc.Columns().Title, "%"+v+"%")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	total, err := model.Clone().Count()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "count knowledge docs failed")
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var rows []entity.AgentKnowledgeDoc
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if err = model.Page(page, pageSize).Scan(&rows); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query knowledge docs failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	list := make([]*agentv1.KnowledgeDoc, 0, len(rows))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := range rows {
		// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
		list = append(list, toProtoKnowledgeDoc(&rows[i]))
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	totalPages := int32(0)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if total > 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		totalPages = int32((total + pageSize - 1) / pageSize)
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.ListKnowledgeDocsRes{List: list, Total: uint64(total), HasMore: page*pageSize < total, TotalPages: totalPages}, nil
}

func (s *sAgent) SubmitKnowledgeDoc(ctx context.Context, req *agentv1.SubmitKnowledgeDocReq) (*agentv1.SubmitKnowledgeDocRes, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	res, err := s.transitionKnowledgeDoc(ctx, req.GetKnowledgeDocNo(), req.GetExpectedVersion(), "REVIEWING", "", "")
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.SubmitKnowledgeDocRes{Doc: res.Doc}, nil
}

func (s *sAgent) ApproveKnowledgeDoc(ctx context.Context, req *agentv1.ApproveKnowledgeDocReq) (*agentv1.ApproveKnowledgeDocRes, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	res, err := s.transitionKnowledgeDoc(ctx, req.GetKnowledgeDocNo(), req.GetExpectedVersion(), "APPROVED", "", "")
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.ApproveKnowledgeDocRes{Doc: res.Doc}, nil
}

func (s *sAgent) RejectKnowledgeDoc(ctx context.Context, req *agentv1.RejectKnowledgeDocReq) (*agentv1.RejectKnowledgeDocRes, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	res, err := s.transitionKnowledgeDoc(ctx, req.GetKnowledgeDocNo(), req.GetExpectedVersion(), "REJECTED", req.GetRejectReasonCode(), req.GetRejectComment())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.RejectKnowledgeDocRes{Doc: res.Doc}, nil
}

func (s *sAgent) PublishKnowledgeDoc(ctx context.Context, req *agentv1.PublishKnowledgeDocReq) (*agentv1.PublishKnowledgeDocRes, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	res, err := s.transitionKnowledgeDoc(ctx, req.GetKnowledgeDocNo(), req.GetExpectedVersion(), "PUBLISHED", "", "")
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.PublishKnowledgeDocRes{Doc: res.Doc}, nil
}

func (s *sAgent) OfflineKnowledgeDoc(ctx context.Context, req *agentv1.OfflineKnowledgeDocReq) (*agentv1.OfflineKnowledgeDocRes, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	res, err := s.transitionKnowledgeDoc(ctx, req.GetKnowledgeDocNo(), req.GetExpectedVersion(), "OFFLINE", req.GetReasonCode(), "")
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.OfflineKnowledgeDocRes{Doc: res.Doc}, nil
}

func (s *sAgent) ReindexKnowledgeDoc(ctx context.Context, req *agentv1.ReindexKnowledgeDocReq) (*agentv1.ReindexKnowledgeDocRes, error) {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(req.GetKnowledgeDocNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "knowledge_doc_no is required")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.ReindexKnowledgeDocRes{KnowledgeDocNo: req.GetKnowledgeDocNo(), StatusCode: "ACCEPTED"}, nil
}

func (s *sAgent) SearchKnowledgeChunks(ctx context.Context, req *agentv1.SearchKnowledgeChunksReq) (*agentv1.SearchKnowledgeChunksRes, error) {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	query := strings.TrimSpace(req.GetQuery())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if query == "" {
		// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
		return &agentv1.SearchKnowledgeChunksRes{}, nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	limit := int(req.GetLimit())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if limit <= 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		limit = defaultKnowledgeSize
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if limit > 50 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		limit = 50
	}
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	model := dao.AgentKnowledgeChunk.Ctx(ctx).Where(dao.AgentKnowledgeChunk.Columns().IsDeleted, 0).OrderAsc(dao.AgentKnowledgeChunk.Columns().ChunkIndex).Limit(limit).WhereLike(dao.AgentKnowledgeChunk.Columns().ChunkText, "%"+query+"%")
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetSourceTypeCode()); v != "" {
		// 统一把编码转换成约定的大写形式，避免状态码和场景码因大小写不一致而失配。
		model = model.Where(dao.AgentKnowledgeChunk.Columns().SourceTypeCode, strings.ToUpper(v))
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if v := strings.TrimSpace(req.GetSourceId()); v != "" {
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		model = model.Where(dao.AgentKnowledgeChunk.Columns().SourceId, v)
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var rows []entity.AgentKnowledgeChunk
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if err := model.Scan(&rows); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "search knowledge chunks failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	list := make([]*agentv1.SearchKnowledgeChunkHit, 0, len(rows))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := range rows {
		// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
		list = append(list, &agentv1.SearchKnowledgeChunkHit{Chunk: toProtoKnowledgeChunk(&rows[i]), Score: 1})
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.SearchKnowledgeChunksRes{List: list}, nil
}

func (s *sAgent) transitionKnowledgeDoc(ctx context.Context, knowledgeDocNo string, expectedVersion uint64, statusCode, reasonCode, comment string) (*agentv1.UpdateKnowledgeDocRes, error) {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(knowledgeDocNo) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "knowledge_doc_no is required")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	doc, err := s.getKnowledgeDocByNo(ctx, knowledgeDocNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if expectedVersion > 0 && doc.Version != expectedVersion {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "knowledge doc version conflict")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	data := do.AgentKnowledgeDoc{StatusCode: statusCode, Version: doc.Version + 1}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if reasonCode != "" {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		data.RejectReasonCode = strings.ToUpper(strings.TrimSpace(reasonCode))
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if comment != "" {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		data.RejectComment = strings.TrimSpace(comment)
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if statusCode == "PUBLISHED" {
		// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
		data.PublishedAt = gtime.Now()
	}
	// 将本次状态变化写回数据库，保证状态机推进结果能被后续查询立即观察到。
	_, err = dao.AgentKnowledgeDoc.Ctx(ctx).Where(dao.AgentKnowledgeDoc.Columns().KnowledgeDocNo, doc.KnowledgeDocNo).Data(data).Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "update knowledge doc status failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	updated, err := s.getKnowledgeDocByNo(ctx, doc.KnowledgeDocNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.UpdateKnowledgeDocRes{Doc: toProtoKnowledgeDoc(updated)}, nil
}
