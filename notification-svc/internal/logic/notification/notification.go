package notification

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/notification-svc/api/v1"
	"github.com/TsingpekTao/shopa/notification-svc/internal/dao"
	"github.com/TsingpekTao/shopa/notification-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/notification-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/notification-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sNotification 是通知领域服务实现。
type sNotification struct{}

// init 在包加载时注册通知服务实现。
func init() {
	// 把逻辑实现注册到 service 层，供 controller 统一调用。
	service.RegisterNotification(New())
}

// New 创建通知服务实例。
func New() *sNotification {
	// 返回无状态实例，便于并发复用。
	return &sNotification{}
}

// SendTemplateMessage 发送单条模板消息。
func (s *sNotification) SendTemplateMessage(ctx context.Context, req *v1.SendTemplateMessageReq) (*v1.SendTemplateMessageRes, error) {
	// 读取并去除模板编码两端空白，避免空字符串写入任务表。
	templateCode := strings.TrimSpace(req.GetTemplateCode())
	// 校验模板编码必填。
	if templateCode == "" {
		// 返回参数错误，提示调用方补齐模板编码。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "template_code is required")
	}

	// 读取并去除目标地址两端空白，支持未登录态直接发目标地址。
	targetAddress := strings.TrimSpace(req.GetTargetAddress())
	// 校验 user_id 与 target_address 至少一个可用，保证有明确接收者。
	if req.GetUserId() == 0 && targetAddress == "" {
		// 返回参数错误，提示必须提供用户或目标地址。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id or target_address is required")
	}

	// 初始化是否被用户偏好拦截标记。
	blockedByPreference := false
	// 初始化是否被频控拦截标记。
	blockedByFrequency := false
	// 初始化消息状态为待投递。
	status := v1.DeliveryStatus_DELIVERY_STATUS_PENDING

	// 仅在营销消息且 user_id 有效时检查用户偏好。
	if req.GetBizType() == v1.NotificationBizType_NOTIFICATION_BIZ_TYPE_MARKETING && req.GetUserId() > 0 {
		// 查询当前用户在指定渠道下是否允许营销触达。
		allowed, err := s.allowMarketingByPreference(ctx, req.GetUserId(), req.GetChannel())
		// 偏好查询失败直接返回错误，避免误发营销消息。
		if err != nil {
			// 向上返回查询错误。
			return nil, err
		}
		// 如果不允许营销触达则标记偏好拦截。
		if !allowed {
			// 记录拦截标记用于响应与审计。
			blockedByPreference = true
		}
	}

	// 仅在未被偏好拦截时才进行频控检查。
	if !blockedByPreference {
		// 执行默认频控：同模板 1 分钟最多 1 条、1 天最多 5 条。
		hit, err := s.hitFrequencyLimit(ctx, req.GetUserId(), templateCode, targetAddress)
		// 频控查询失败直接返回错误，避免越限发送。
		if err != nil {
			// 向上返回频控查询异常。
			return nil, err
		}
		// 命中频控时设置频控拦截标记。
		if hit {
			// 记录拦截标记用于响应与审计。
			blockedByFrequency = true
		}
	}

	// 偏好或频控命中时将状态置为 dropped。
	if blockedByPreference || blockedByFrequency {
		// 标记消息为被丢弃，表示逻辑拦截未下发渠道。
		status = v1.DeliveryStatus_DELIVERY_STATUS_DROPPED
	} else {
		// 当前版本使用 MOCK 渠道直返 sent 状态。
		status = v1.DeliveryStatus_DELIVERY_STATUS_SENT
	}

	// 生成业务通知号作为任务唯一标识。
	notificationNo := genNo("NTF")

	// 写入通知任务主表，沉淀全量发送审计记录。
	_, err := dao.NotificationTask.Ctx(ctx).Data(do.NotificationTask{
		NotificationNo:      notificationNo,                     // 记录通知号用于后续状态查询。
		UserId:              req.GetUserId(),                    // 记录用户 ID，未登录态时为 0。
		TargetAddress:       targetAddress,                      // 记录目标地址，支持未登录态短信/邮箱。
		TemplateCode:        templateCode,                       // 记录模板编码用于追踪与频控维度。
		Channel:             uint(req.GetChannel()),             // 记录消息渠道枚举值。
		BizType:             uint(req.GetBizType()),             // 记录业务类型枚举值。
		TemplateParamsJson:  req.GetTemplateParamsJson(),        // 记录模板参数快照。
		Status:              uint(status),                       // 记录当前投递状态。
		BlockedByPreference: boolToTinyint(blockedByPreference), // 记录是否被偏好拦截。
		BlockedByFrequency:  boolToTinyint(blockedByFrequency),  // 记录是否被频控拦截。
		RetryCount:          0,                                  // 初始化重试次数为 0。
	}).Insert()
	// 任务写入失败时直接返回错误。
	if err != nil {
		// 包装错误栈并向上抛出。
		return nil, gerror.Wrap(err, "insert notification_task failed")
	}

	// 非 dropped 状态时补写渠道发送记录。
	if status != v1.DeliveryStatus_DELIVERY_STATUS_DROPPED {
		// 生成模拟渠道消息 ID。
		providerMessageID := genNo("MSG")
		// 写入渠道记录表用于回溯第三方交互结果。
		_, _ = dao.NotificationChannelRecord.Ctx(ctx).Data(do.NotificationChannelRecord{
			NotificationNo:    notificationNo,              // 关联通知号。
			ProviderCode:      "MOCK",                      // 标识当前使用模拟渠道。
			ProviderMessageId: providerMessageID,           // 记录渠道消息 ID。
			RequestPayload:    req.GetTemplateParamsJson(), // 记录发送请求参数。
			ResponsePayload:   `{"ok":true}`,               // 记录模拟发送响应。
			Status:            2,                           // 标记渠道记录状态为 sent。
			ErrorCode:         "",                          // 正常发送无错误码。
			ErrorMessage:      "",                          // 正常发送无错误信息。
		}).Insert()
	}

	// 组装并返回发送结果。
	return &v1.SendTemplateMessageRes{
		NotificationNo:      notificationNo,      // 返回通知号供后续查询。
		Status:              status,              // 返回本次最终状态。
		BlockedByPreference: blockedByPreference, // 返回偏好拦截标记。
		BlockedByFrequency:  blockedByFrequency,  // 返回频控拦截标记。
	}, nil
}

// BatchSendTemplateMessage 批量发送模板消息。
func (s *sNotification) BatchSendTemplateMessage(ctx context.Context, req *v1.BatchSendTemplateMessageReq) (*v1.BatchSendTemplateMessageRes, error) {
	// 初始化批量响应容器。
	resp := &v1.BatchSendTemplateMessageRes{
		Items: make([]*v1.SendTemplateMessageRes, 0, len(req.GetItems())), // 预分配容量减少扩容开销。
	}
	// 逐条复用单发逻辑，确保规则一致。
	for _, item := range req.GetItems() {
		// 执行单条发送。
		one, err := s.SendTemplateMessage(ctx, item)
		// 任一条失败直接中断并返回错误。
		if err != nil {
			// 向上返回第一条失败错误。
			return nil, err
		}
		// 将单条结果追加到批量结果。
		resp.Items = append(resp.Items, one)
	}
	// 返回批量发送结果。
	return resp, nil
}

// GetDeliveryStatus 查询单条通知投递状态。
func (s *sNotification) GetDeliveryStatus(ctx context.Context, req *v1.GetDeliveryStatusReq) (*v1.GetDeliveryStatusRes, error) {
	// 提取并去除通知号两端空白。
	notificationNo := strings.TrimSpace(req.GetNotificationNo())
	// 校验通知号必填。
	if notificationNo == "" {
		// 返回参数错误，提示必须传 notification_no。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "notification_no is required")
	}

	// 声明任务实体用于接收查询结果。
	var task entity.NotificationTask
	// 根据通知号查询任务主记录。
	err := dao.NotificationTask.Ctx(ctx).Where(do.NotificationTask{
		NotificationNo: notificationNo, // 精确命中该通知任务。
	}).Scan(&task)
	// 主记录查询失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return nil, gerror.Wrap(err, "query notification_task failed")
	}
	// 未命中任务时返回 not found。
	if task.NotificationNo == "" {
		// 返回不存在错误供调用方处理。
		return nil, gerror.NewCode(gcode.CodeNotFound, "notification not found")
	}

	// 初始化渠道消息 ID 返回值为空。
	providerMessageID := ""
	// 初始化错误码返回值为空。
	errorCode := ""
	// 初始化错误信息返回值为空。
	errorMessage := ""

	// 声明渠道记录实体用于读取最近一次渠道结果。
	var record entity.NotificationChannelRecord
	// 按通知号查询最新一条渠道记录。
	_ = dao.NotificationChannelRecord.Ctx(ctx).Where(do.NotificationChannelRecord{
		NotificationNo: notificationNo, // 精确命中同一通知的渠道记录。
	}).OrderDesc(dao.NotificationChannelRecord.Columns().Id).Scan(&record)
	// 若渠道记录存在则补齐渠道字段。
	if record.NotificationNo != "" {
		// 填充渠道消息 ID。
		providerMessageID = record.ProviderMessageId
		// 填充渠道错误码。
		errorCode = record.ErrorCode
		// 填充渠道错误信息。
		errorMessage = record.ErrorMessage
	}

	// 返回投递状态响应。
	return &v1.GetDeliveryStatusRes{
		NotificationNo:    task.NotificationNo,            // 返回通知号。
		Status:            v1.DeliveryStatus(task.Status), // 返回任务状态枚举。
		ProviderMessageId: providerMessageID,              // 返回渠道消息 ID。
		ErrorCode:         errorCode,                      // 返回错误码。
		ErrorMessage:      errorMessage,                   // 返回错误信息。
	}, nil
}

// GetMyNotificationPreference 查询当前用户通知偏好。
func (s *sNotification) GetMyNotificationPreference(ctx context.Context, req *v1.GetMyNotificationPreferenceReq) (*v1.GetMyNotificationPreferenceRes, error) {
	// 从上下文提取用户 ID。
	userID, err := userIDFromContext(ctx)
	// 提取失败说明未鉴权，直接返回错误。
	if err != nil {
		// 向上返回未授权错误。
		return nil, err
	}
	// 查询或初始化用户偏好记录。
	preference, err := s.getOrInitPreference(ctx, userID)
	// 查询失败时返回错误。
	if err != nil {
		// 向上返回数据库错误。
		return nil, err
	}
	// 返回用户偏好响应。
	return &v1.GetMyNotificationPreferenceRes{
		Preference: toProtoPreference(preference), // 返回协议层偏好对象。
	}, nil
}

// UpdateMyNotificationPreference 更新当前用户通知偏好。
func (s *sNotification) UpdateMyNotificationPreference(ctx context.Context, req *v1.UpdateMyNotificationPreferenceReq) (*v1.UpdateMyNotificationPreferenceRes, error) {
	// 从上下文提取用户 ID。
	userID, err := userIDFromContext(ctx)
	// 提取失败说明未鉴权，直接返回错误。
	if err != nil {
		// 向上返回未授权错误。
		return nil, err
	}

	// 先确保偏好记录存在，避免更新空记录。
	preference, err := s.getOrInitPreference(ctx, userID)
	// 初始化失败时返回错误。
	if err != nil {
		// 向上返回数据库错误。
		return nil, err
	}

	// 执行偏好更新，仅开放营销和 push 开关给用户修改。
	_, err = dao.UserNotificationPreference.Ctx(ctx).Where(do.UserNotificationPreference{
		UserId: userID, // 定位当前用户偏好行。
	}).Data(do.UserNotificationPreference{
		AllowTransactionalSms:   preference.AllowTransactionalSms,            // 事务类短信保持系统默认策略。
		AllowMarketingSms:       boolToTinyint(req.GetAllowMarketingSms()),   // 更新营销短信开关。
		AllowTransactionalEmail: preference.AllowTransactionalEmail,          // 事务类邮件保持系统默认策略。
		AllowMarketingEmail:     boolToTinyint(req.GetAllowMarketingEmail()), // 更新营销邮件开关。
		AllowPush:               boolToTinyint(req.GetAllowPush()),           // 更新推送开关。
	}).Update()
	// 更新失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return nil, gerror.Wrap(err, "update user_notification_preference failed")
	}

	// 重新查询最新偏好用于响应回显。
	updated, err := s.getOrInitPreference(ctx, userID)
	// 回查失败时返回错误。
	if err != nil {
		// 向上返回数据库错误。
		return nil, err
	}

	// 返回更新后的偏好响应。
	return &v1.UpdateMyNotificationPreferenceRes{
		Preference: toProtoPreference(updated), // 返回最新偏好对象。
	}, nil
}

// RetryFailedMessage 重试失败消息。
func (s *sNotification) RetryFailedMessage(ctx context.Context, req *v1.RetryFailedMessageReq) (*v1.RetryFailedMessageRes, error) {
	// 提取并去除通知号两端空白。
	notificationNo := strings.TrimSpace(req.GetNotificationNo())
	// 校验通知号必填。
	if notificationNo == "" {
		// 返回参数错误，提示调用方补齐通知号。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "notification_no is required")
	}

	// 更新任务状态为 pending 并递增重试次数。
	result, err := dao.NotificationTask.Ctx(ctx).Where(do.NotificationTask{
		NotificationNo: notificationNo, // 精确命中通知任务。
	}).Data(do.NotificationTask{
		Status:      uint(v1.DeliveryStatus_DELIVERY_STATUS_PENDING), // 重置为待投递。
		RetryCount:  g.DB().Raw("retry_count + 1"),                   // 原子递增重试次数。
		NextRetryAt: gtime.Now(),                                     // 记录下一次重试触发时间。
	}).Update()
	// 更新失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return nil, gerror.Wrap(err, "retry notification failed")
	}
	// 读取受影响行数用于判断任务是否存在。
	affected, rowsErr := result.RowsAffected()
	// 读取行数失败时返回数据库错误。
	if rowsErr != nil {
		// 包装错误栈并返回。
		return nil, gerror.Wrap(rowsErr, "read rows affected failed")
	}
	// 未命中任务时返回 not found。
	if affected == 0 {
		// 返回不存在错误供调用方处理。
		return nil, gerror.NewCode(gcode.CodeNotFound, "notification not found")
	}

	// 返回重试受理结果。
	return &v1.RetryFailedMessageRes{
		NotificationNo: notificationNo,                            // 返回通知号。
		Status:         v1.DeliveryStatus_DELIVERY_STATUS_PENDING, // 返回重试后的目标状态。
	}, nil
}

// allowMarketingByPreference 判断用户是否允许营销类通知。
func (s *sNotification) allowMarketingByPreference(ctx context.Context, userID uint64, channel v1.ChannelType) (bool, error) {
	// 查询或初始化用户偏好记录。
	pref, err := s.getOrInitPreference(ctx, userID)
	// 偏好查询失败时返回错误。
	if err != nil {
		// 返回错误给上游调用方。
		return false, err
	}

	// 按渠道类型返回对应营销开关。
	switch channel {
	case v1.ChannelType_CHANNEL_TYPE_SMS:
		// 短信渠道读取营销短信开关。
		return pref.AllowMarketingSms == 1, nil
	case v1.ChannelType_CHANNEL_TYPE_EMAIL:
		// 邮件渠道读取营销邮件开关。
		return pref.AllowMarketingEmail == 1, nil
	default:
		// 其余渠道统一读取 push 开关。
		return pref.AllowPush == 1, nil
	}
}

// hitFrequencyLimit 执行默认频控（1 分钟 1 条、1 天 5 条）。
func (s *sNotification) hitFrequencyLimit(ctx context.Context, userID uint64, templateCode string, targetAddress string) (bool, error) {
	// 当 user_id 与 target_address 都为空时无法建立频控维度，直接放行。
	if userID == 0 && strings.TrimSpace(targetAddress) == "" {
		// 返回未命中频控。
		return false, nil
	}

	// 计算 1 分钟窗口起点。
	minuteAgo := gtime.Now().Add(-time.Minute)
	// 计算 1 天窗口起点。
	dayAgo := gtime.Now().Add(-24 * time.Hour)

	// 构建按模板编码过滤的基础查询。
	m := dao.NotificationTask.Ctx(ctx).Where(do.NotificationTask{
		TemplateCode: templateCode, // 固定模板维度，避免不同模板互相影响。
	})
	// user_id 存在时按 user_id 维度频控。
	if userID > 0 {
		// 追加 user_id 条件。
		m = m.Where(do.NotificationTask{UserId: userID})
	} else {
		// 未登录态按 target_address 维度频控。
		m = m.Where(do.NotificationTask{TargetAddress: strings.TrimSpace(targetAddress)})
	}

	// 统计 1 分钟内历史发送次数。
	minuteCount, err := m.Clone().WhereGTE(dao.NotificationTask.Columns().CreatedAt, minuteAgo).Count()
	// 分钟窗口统计失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return false, gerror.Wrap(err, "query minute frequency failed")
	}

	// 统计 1 天内历史发送次数。
	dayCount, err := m.Clone().WhereGTE(dao.NotificationTask.Columns().CreatedAt, dayAgo).Count()
	// 天窗口统计失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return false, gerror.Wrap(err, "query day frequency failed")
	}

	// 命中任一窗口阈值即拦截。
	return minuteCount >= 1 || dayCount >= 5, nil
}

// getOrInitPreference 查询或初始化用户偏好记录。
func (s *sNotification) getOrInitPreference(ctx context.Context, userID uint64) (*entity.UserNotificationPreference, error) {
	// 声明偏好实体承接查询结果。
	var pref entity.UserNotificationPreference
	// 按 user_id 查询偏好记录。
	err := dao.UserNotificationPreference.Ctx(ctx).Where(do.UserNotificationPreference{
		UserId: userID, // user_id 唯一约束保证最多一行。
	}).Scan(&pref)
	// 查询失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return nil, gerror.Wrap(err, "query user_notification_preference failed")
	}
	// 记录存在时直接返回。
	if pref.UserId > 0 {
		// 返回已存在偏好对象。
		return &pref, nil
	}

	// 不存在时插入默认偏好配置。
	_, err = dao.UserNotificationPreference.Ctx(ctx).Data(do.UserNotificationPreference{
		UserId:                  userID, // 绑定当前用户。
		AllowTransactionalSms:   1,      // 默认允许事务短信。
		AllowMarketingSms:       0,      // 默认关闭营销短信。
		AllowTransactionalEmail: 1,      // 默认允许事务邮件。
		AllowMarketingEmail:     0,      // 默认关闭营销邮件。
		AllowPush:               1,      // 默认允许推送。
	}).InsertIgnore()
	// 初始化失败时返回错误。
	if err != nil {
		// 包装错误栈并返回。
		return nil, gerror.Wrap(err, "insert user_notification_preference failed")
	}

	// 再次查询并返回最新偏好记录。
	return s.getOrInitPreference(ctx, userID)
}

// toProtoPreference 把实体偏好转换为 proto 偏好对象。
func toProtoPreference(pref *entity.UserNotificationPreference) *v1.UserNotificationPreference {
	// 输入为空时返回空对象避免空指针。
	if pref == nil {
		// 返回空 proto 对象。
		return &v1.UserNotificationPreference{}
	}

	// 返回字段映射后的 proto 偏好对象。
	return &v1.UserNotificationPreference{
		UserId:                  pref.UserId,                       // 映射用户 ID。
		AllowTransactionalSms:   pref.AllowTransactionalSms == 1,   // tinyint 转 bool。
		AllowMarketingSms:       pref.AllowMarketingSms == 1,       // tinyint 转 bool。
		AllowTransactionalEmail: pref.AllowTransactionalEmail == 1, // tinyint 转 bool。
		AllowMarketingEmail:     pref.AllowMarketingEmail == 1,     // tinyint 转 bool。
		AllowPush:               pref.AllowPush == 1,               // tinyint 转 bool。
		UpdatedAt:               toProtoTs(pref.UpdatedAt),         // 时间类型转换。
	}
}

// userIDFromContext 从上下文解析用户 ID。
func userIDFromContext(ctx context.Context) (uint64, error) {
	// 优先读取 HTTP 上下文中的 x-user-id。
	if r := g.RequestFromCtx(ctx); r != nil {
		// 读取并去除请求头空白。
		raw := strings.TrimSpace(r.Header.Get("x-user-id"))
		// 请求头存在时解析为 uint64。
		if raw != "" {
			// 把字符串转为无符号整数。
			uid, err := strconv.ParseUint(raw, 10, 64)
			// 解析成功且 uid 合法时直接返回。
			if err == nil && uid > 0 {
				// 返回解析后的用户 ID。
				return uid, nil
			}
		}
	}

	// 再尝试读取 gRPC metadata 中的 x-user-id。
	md, ok := metadata.FromIncomingContext(ctx)
	// metadata 存在时尝试解析。
	if ok {
		// 读取 metadata 对应键值数组。
		values := md.Get("x-user-id")
		// 至少存在一个值时进行解析。
		if len(values) > 0 {
			// 读取第一个用户 ID 值并去空白。
			raw := strings.TrimSpace(values[0])
			// 值非空时执行数值解析。
			if raw != "" {
				// 把字符串转为无符号整数。
				uid, err := strconv.ParseUint(raw, 10, 64)
				// 解析成功且 uid 合法时直接返回。
				if err == nil && uid > 0 {
					// 返回解析后的用户 ID。
					return uid, nil
				}
			}
		}
	}

	// 两种上下文都未取到合法用户 ID 时返回未授权错误。
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

// boolToTinyint 把 bool 转成 tinyint 语义值。
func boolToTinyint(v bool) int {
	// true 映射为 1。
	if v {
		// 返回 1。
		return 1
	}
	// false 映射为 0。
	return 0
}

// genNo 生成业务编号。
func genNo(prefix string) string {
	// 使用前缀 + 纳秒时间戳生成近似唯一编号。
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// toProtoTs 把 gtime 转成 protobuf timestamp。
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	// 空时间直接返回 nil。
	if t == nil {
		// 返回 nil 代表字段未设置。
		return nil
	}
	// 转换并返回 protobuf timestamp。
	return timestamppb.New(t.Time)
}
