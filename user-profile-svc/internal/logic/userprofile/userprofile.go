package userprofile

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/consts"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/ctxkey"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/dao"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sUserProfile 实现 user-profile 服务的业务逻辑。
// 核心约束如下：
// 1) 资料字段变更推进 profile_version，供用户资料乐观并发使用；
// 2) 默认地址投影变更推进 address_book_version，保障地址簿版本一致；
// 3) 地址替换采用追加新行 + 旧行退役方式，避免污染旧的快照。
type sUserProfile struct{}

// New 创建 user-profile 逻辑实例。
func New() *sUserProfile {
	return &sUserProfile{}
}

func init() {
	service.RegisterUserProfile(New())
}

// GetMyProfile 返回当前用户资料，可按需带地址列表。
func (s *sUserProfile) GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.ensureProfile(ctx, userID); err != nil {
		return nil, err
	}
	profile, err := s.getProfileEntity(ctx, userID)
	if err != nil {
		return nil, err
	}

	res = &v1.GetMyProfileRes{
		Profile: toProtoProfile(profile),
	}
	if req.GetIncludeAddresses() {
		addresses, listErr := s.listAddressEntities(ctx, userID, false, 1, 500)
		if listErr != nil {
			return nil, listErr
		}
		res.Addresses = toProtoAddressList(addresses)
	}
	return res, nil
}

// UpdateMyProfile 按 FieldMask 更新用户资料。
func (s *sUserProfile) UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error) {
	if req.GetProfile() == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "profile is required")
	}
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "update_mask.paths is required")
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.ensureProfile(ctx, userID); err != nil {
		return nil, err
	}

	var (
		// p 是请求中的变更值，按 update_mask 选择落库字段。
		p            = req.GetProfile()
		data         do.UserProfile
		updatedCount = 0
	)

	for _, rawPath := range req.GetUpdateMask().GetPaths() {
		path, normalizeErr := normalizeProfileMaskPath(rawPath)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		switch path {
		case "display_name":
			data.DisplayName = p.GetDisplayName()
			data.DisplayNameSource = consts.DisplayNameSourceUserSet
			updatedCount++

		case "avatar":
			if p.GetAvatar() != nil {
				data.AvatarAssetId = p.GetAvatar().GetAssetId()
				data.AvatarUrl = p.GetAvatar().GetUrl()
			} else {
				data.AvatarAssetId = uint64(0)
				data.AvatarUrl = ""
			}
			updatedCount++

		case "gender":
			if p.GetGender() < v1.Gender_GENDER_UNSPECIFIED || p.GetGender() > v1.Gender_GENDER_OTHER {
				return nil, gerror.NewCode(gcode.CodeInvalidParameter, "gender is invalid")
			}
			data.Gender = uint(p.GetGender())
			updatedCount++

		case "birthday":
			if strings.TrimSpace(p.GetBirthday()) == "" {
				return nil, gerror.NewCode(gcode.CodeInvalidParameter, "birthday cannot be empty in v1")
			}
			t, parseErr := time.Parse("2006-01-02", p.GetBirthday())
			if parseErr != nil {
				return nil, gerror.WrapCode(gcode.CodeInvalidParameter, parseErr, "birthday format must be YYYY-MM-DD")
			}
			data.Birthday = gtime.NewFromTime(t)
			updatedCount++

		case "locale":
			data.Locale = p.GetLocale()
			updatedCount++

		case "timezone":
			data.Timezone = p.GetTimezone()
			updatedCount++

		case "marketing_opt_in":
			data.MarketingOptIn = gconv.Int(p.GetMarketingOptIn())
			updatedCount++

		case "ext":
			if len(p.GetExt()) == 0 {
				data.Ext = "{}"
			} else {
				extJSON, marshalErr := json.Marshal(p.GetExt())
				if marshalErr != nil {
					return nil, gerror.WrapCode(gcode.CodeInvalidParameter, marshalErr, "ext marshal failed")
				}
				data.Ext = string(extJSON)
			}
			updatedCount++
		}
	}

	if updatedCount == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "no valid update fields")
	}

	cols := dao.UserProfile.Columns()
	// 递增 profile_version 以支持乐观并发控制。
	data.ProfileVersion = gdb.Raw(cols.ProfileVersion + " + 1")
	_, err = dao.UserProfile.Ctx(ctx).
		Where(cols.UserId, userID).
		Data(data).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "update profile failed")
	}

	profile, err := s.getProfileEntity(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMyProfileRes{
		Profile: toProtoProfile(profile),
	}, nil
}

// ListMyAddresses 返回当前用户地址列表。
func (s *sUserProfile) ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 新 proto 去掉分页请求字段，仍保持兼容的响应字段值。
	page := consts.DefaultPage
	pageSize := consts.MaxPageSize

	cols := dao.UserAddress.Columns()
	model := dao.UserAddress.Ctx(ctx).Where(cols.UserId, userID)
	if !req.GetIncludeDeleted() {
		model = model.Where(cols.Status, consts.AddressStatusActive)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count addresses failed")
	}

	addresses, err := s.listAddressEntities(ctx, userID, req.GetIncludeDeleted(), page, pageSize)
	if err != nil {
		return nil, err
	}

	return &v1.ListMyAddressesRes{
		Addresses: toProtoAddressList(addresses),
		Page:      uint32(page),
		PageSize:  uint32(pageSize),
		Total:     uint32(total),
	}, nil
}

// CreateMyAddress 创建地址。
func (s *sUserProfile) CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error) {
	if req.GetAddress() == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.ensureProfile(ctx, userID); err != nil {
		return nil, err
	}

	address := req.GetAddress()
	extStr, err := marshalExt(address.GetExt())
	if err != nil {
		return nil, err
	}

	var created *entity.UserAddress
	err = dao.UserProfile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 插入地址记录。
		r, insertErr := tx.Model(dao.UserAddress.Table()).
			Data(do.UserAddress{
				UserId:        userID,
				Status:        consts.AddressStatusActive,
				Label:         address.GetLabel(),
				ReceiverName:  address.GetReceiverName(),
				ReceiverPhone: address.GetReceiverPhone(),
				CountryCode:   address.GetCountryCode(),
				ProvinceCode:  address.GetProvinceCode(),
				ProvinceName:  address.GetProvinceName(),
				CityCode:      address.GetCityCode(),
				CityName:      address.GetCityName(),
				DistrictCode:  address.GetDistrictCode(),
				DistrictName:  address.GetDistrictName(),
				Street:        address.GetStreet(),
				Detail:        address.GetDetail(),
				PostalCode:    address.GetPostalCode(),
				IsDefault:     0,
				Latitude:      address.GetLatitude(),
				Longitude:     address.GetLongitude(),
				Ext:           extStr,
			}).
			Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "insert address failed")
		}
		lastID, idErr := r.LastInsertId()
		if idErr != nil {
			return gerror.Wrap(idErr, "read last insert id failed")
		}
		addressID := uint64(lastID)

		// 决定是否将其设为默认地址。
		if req.GetSetAsDefault() {
			if setErr := setDefaultAddressTx(ctx, tx, userID, addressID); setErr != nil {
				return setErr
			}
		} else {
			// 若画像尚无默认地址，则自动设置当前地址为默认。
			var profile entity.UserProfile
			cols := dao.UserProfile.Columns()
			if scanErr := tx.Model(dao.UserProfile.Table()).Where(cols.UserId, userID).Scan(&profile); scanErr != nil {
				return gerror.Wrap(scanErr, "query profile in tx failed")
			}
			if profile.DefaultAddressId == 0 {
				if setErr := setDefaultAddressTx(ctx, tx, userID, addressID); setErr != nil {
					return setErr
				}
			}
		}

		// 读取刚插入的地址行。
		var item entity.UserAddress
		addrCols := dao.UserAddress.Columns()
		if scanErr := tx.Model(dao.UserAddress.Table()).
			Where(addrCols.AddressId, addressID).
			Where(addrCols.UserId, userID).
			Scan(&item); scanErr != nil {
			return gerror.Wrap(scanErr, "query created address failed")
		}
		created = &item
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateMyAddressRes{
		Address: toProtoAddress(created),
	}, nil
}

// UpdateMyAddress 更新指定地址。
func (s *sUserProfile) UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error) {
	if req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address_id is required")
	}
	if req.GetAddress() == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address is required")
	}
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "update_mask.paths is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var (
		p            = req.GetAddress()
		data         do.UserAddress
		updatedCount = 0
	)
	for _, rawPath := range req.GetUpdateMask().GetPaths() {
		path, normalizeErr := normalizeAddressMaskPath(rawPath)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		switch path {
		case "label":
			data.Label = p.GetLabel()
			updatedCount++
		case "receiver_name":
			data.ReceiverName = p.GetReceiverName()
			updatedCount++
		case "receiver_phone":
			data.ReceiverPhone = p.GetReceiverPhone()
			updatedCount++
		case "country_code":
			data.CountryCode = p.GetCountryCode()
			updatedCount++
		case "province_code":
			data.ProvinceCode = p.GetProvinceCode()
			updatedCount++
		case "province_name":
			data.ProvinceName = p.GetProvinceName()
			updatedCount++
		case "city_code":
			data.CityCode = p.GetCityCode()
			updatedCount++
		case "city_name":
			data.CityName = p.GetCityName()
			updatedCount++
		case "district_code":
			data.DistrictCode = p.GetDistrictCode()
			updatedCount++
		case "district_name":
			data.DistrictName = p.GetDistrictName()
			updatedCount++
		case "street":
			data.Street = p.GetStreet()
			updatedCount++
		case "detail":
			data.Detail = p.GetDetail()
			updatedCount++
		case "postal_code":
			data.PostalCode = p.GetPostalCode()
			updatedCount++
		case "latitude":
			data.Latitude = p.GetLatitude()
			updatedCount++
		case "longitude":
			data.Longitude = p.GetLongitude()
			updatedCount++
		case "ext":
			extStr, extErr := marshalExt(p.GetExt())
			if extErr != nil {
				return nil, extErr
			}
			data.Ext = extStr
			updatedCount++
		}
	}
	if updatedCount == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "no valid update fields")
	}

	cols := dao.UserAddress.Columns()
	_, err = dao.UserAddress.Ctx(ctx).
		Where(cols.AddressId, req.GetAddressId()).
		Where(cols.UserId, userID).
		Where(cols.Status, consts.AddressStatusActive).
		Data(data).
		Update()
	if err != nil {
		return nil, gerror.Wrap(err, "update address failed")
	}

	item, err := s.getMyAddressEntity(ctx, userID, req.GetAddressId())
	if err != nil {
		return nil, err
	}
	return &v1.UpdateMyAddressRes{
		Address: toProtoAddress(item),
	}, nil
}

// DeleteMyAddress 软删除地址，并在命中默认地址时同步清空默认指针。
// 注意：这里只改状态，不做物理删除，便于后续审计与问题追踪。
func (s *sUserProfile) DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error) {
	if req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address_id is required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.ensureProfile(ctx, userID); err != nil {
		return nil, err
	}

	err = dao.UserProfile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		addrCols := dao.UserAddress.Columns()

		// 确保地址存在且属于当前用户。
		count, countErr := tx.Model(dao.UserAddress.Table()).
			Where(addrCols.AddressId, req.GetAddressId()).
			Where(addrCols.UserId, userID).
			Where(addrCols.Status, consts.AddressStatusActive).
			Count()
		if countErr != nil {
			return gerror.Wrap(countErr, "check address existence failed")
		}
		if count == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "address not found")
		}

		// 软删除地址并清除默认标记。
		_, delErr := tx.Model(dao.UserAddress.Table()).
			Where(addrCols.AddressId, req.GetAddressId()).
			Where(addrCols.UserId, userID).
			Data(do.UserAddress{
				Status:    consts.AddressStatusDeleted,
				IsDefault: 0,
			}).
			Update()
		if delErr != nil {
			return gerror.Wrap(delErr, "delete address failed")
		}

		// 若该地址为画像默认，则清理 profile.default_address_id。
		profileCols := dao.UserProfile.Columns()
		_, clearErr := tx.Model(dao.UserProfile.Table()).
			Where(profileCols.UserId, userID).
			Where(profileCols.DefaultAddressId, req.GetAddressId()).
			Data(do.UserProfile{
				DefaultAddressId: gdb.Raw("NULL"),
			}).
			Update()
		if clearErr != nil {
			return gerror.Wrap(clearErr, "clear default address id failed")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.DeleteMyAddressRes{}, nil
}

// SetMyDefaultAddress 设置或清空默认地址。
func (s *sUserProfile) SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.ensureProfile(ctx, userID); err != nil {
		return nil, err
	}
	// addressBookVersion 返回给调用方，配合 expected_address_book_version 做并发控制。
	var addressBookVersion uint64

	err = dao.UserProfile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// address_id 在 proto 中为可选，nil/0 表示清除默认地址。
		if req.AddressId == nil || req.GetAddressId() == 0 {
			if clearErr := clearDefaultAddressTx(ctx, tx, userID); clearErr != nil {
				return clearErr
			}
		} else {
			if setErr := setDefaultAddressTx(ctx, tx, userID, req.GetAddressId()); setErr != nil {
				return setErr
			}
		}
		addressBookVersion, err = getAddressBookVersionTx(ctx, tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &v1.SetMyDefaultAddressRes{
		DefaultAddressId:   req.GetAddressId(),
		AddressBookVersion: addressBookVersion,
	}, nil
}

// BatchGetProfileSummary 批量返回资料摘要。
func (s *sUserProfile) BatchGetProfileSummary(ctx context.Context, req *v1.BatchGetProfileSummaryReq) (res *v1.BatchGetProfileSummaryRes, err error) {
	ids := req.GetUserIds()
	if len(ids) == 0 {
		return &v1.BatchGetProfileSummaryRes{}, nil
	}
	if len(ids) > 500 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "too many user_ids")
	}

	cols := dao.UserProfile.Columns()
	var profiles []entity.UserProfile
	if err := dao.UserProfile.Ctx(ctx).
		WhereIn(cols.UserId, toInterfaceSliceUint64(ids)).
		Scan(&profiles); err != nil {
		return nil, gerror.Wrap(err, "batch query profiles failed")
	}

	m := make(map[uint64]*entity.UserProfile, len(profiles))
	for i := range profiles {
		p := profiles[i]
		m[p.UserId] = &p
	}

	out := make([]*v1.UserProfileSummary, 0, len(ids))
	for _, id := range ids {
		if p := m[id]; p != nil {
			out = append(out, &v1.UserProfileSummary{
				UserId:      p.UserId,
				DisplayName: p.DisplayName,
				Avatar:      &v1.ImageRef{AssetId: p.AvatarAssetId, Url: p.AvatarUrl},
			})
		} else {
			out = append(out, &v1.UserProfileSummary{
				UserId: id,
			})
		}
	}
	return &v1.BatchGetProfileSummaryRes{Users: out}, nil
}

// GetProfileByUserId 内部按 user_id 查询资料。
func (s *sUserProfile) GetProfileByUserId(ctx context.Context, req *v1.GetProfileByUserIdReq) (res *v1.GetProfileByUserIdRes, err error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	profile, err := s.getProfileEntity(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	out := &v1.GetProfileByUserIdRes{
		Profile: toProtoProfile(profile),
	}
	if req.GetIncludeAddresses() {
		addresses, err := s.listAddressEntities(ctx, req.GetUserId(), false, 1, 500)
		if err != nil {
			return nil, err
		}
		out.Addresses = toProtoAddressList(addresses)
	}
	return out, nil
}

// ListAddressesByUserId 内部按 user_id 查询地址列表。
func (s *sUserProfile) ListAddressesByUserId(ctx context.Context, req *v1.ListAddressesByUserIdReq) (res *v1.ListAddressesByUserIdRes, err error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	addresses, err := s.listAddressEntities(ctx, req.GetUserId(), req.GetIncludeDeleted(), 1, 500)
	if err != nil {
		return nil, err
	}
	return &v1.ListAddressesByUserIdRes{
		Addresses: toProtoAddressList(addresses),
	}, nil
}

// GetAddressById 内部按 address_id 查询地址。
func (s *sUserProfile) GetAddressById(ctx context.Context, req *v1.GetAddressByIdReq) (res *v1.GetAddressByIdRes, err error) {
	if req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address_id is required")
	}
	item, err := s.getAddressEntityByID(ctx, req.GetAddressId())
	if err != nil {
		return nil, err
	}
	return &v1.GetAddressByIdRes{
		Address: toProtoAddress(item),
	}, nil
}

// ReplaceMyAddress 使用“新建一条+旧地址置为 REPLACED”的方式修改地址。
// 这样可避免订单等下游引用旧 address_id 时被原地更新污染。
func (s *sUserProfile) ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error) {
	if req.GetSourceAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "source_address_id is required")
	}
	if req.GetAddress() == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address is required")
	}
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "update_mask.paths is required")
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err = s.ensureProfile(ctx, userID); err != nil {
		return nil, err
	}

	var (
		newAddressEntity entity.UserAddress
		// newAddressID 是替换后新地址主键。
		newAddressID uint64
		// addressBookVersion 用于告诉客户端地址簿最新版本。
		addressBookVersion uint64
	)

	err = dao.UserProfile.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		addrCols := dao.UserAddress.Columns()
		profileCols := dao.UserProfile.Columns()

		// 锁住源地址行以避免并发重复替换。
		var source entity.UserAddress
		if scanErr := tx.Model(dao.UserAddress.Table()).
			Where(addrCols.AddressId, req.GetSourceAddressId()).
			Where(addrCols.UserId, userID).
			Where(addrCols.Status, consts.AddressStatusActive).
			LockUpdate().
			Scan(&source); scanErr != nil {
			return gerror.Wrap(scanErr, "query source address failed")
		}
		if source.AddressId == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "source address not found")
		}
		if req.GetExpectedSourceAddressVersion() > 0 && source.AddressVersion != req.GetExpectedSourceAddressVersion() {
			return gerror.NewCode(gcode.CodeInvalidParameter, "source address version conflict")
		}

		// 基于原地址与请求字段构建替换行。
		data := do.UserAddress{
			UserId:                userID,
			Status:                consts.AddressStatusActive,
			AddressVersion:        1,
			ReplacedFromAddressId: source.AddressId,
			Label:                 source.Label,
			ReceiverName:          source.ReceiverName,
			ReceiverPhone:         source.ReceiverPhone,
			CountryCode:           source.CountryCode,
			ProvinceCode:          source.ProvinceCode,
			ProvinceName:          source.ProvinceName,
			CityCode:              source.CityCode,
			CityName:              source.CityName,
			DistrictCode:          source.DistrictCode,
			DistrictName:          source.DistrictName,
			Street:                source.Street,
			Detail:                source.Detail,
			PostalCode:            source.PostalCode,
			IsDefault:             0,
			Latitude:              source.Latitude,
			Longitude:             source.Longitude,
			Ext:                   source.Ext,
		}

		for _, rawPath := range req.GetUpdateMask().GetPaths() {
			path, normalizeErr := normalizeAddressMaskPath(rawPath)
			if normalizeErr != nil {
				return normalizeErr
			}
			switch path {
			case "label":
				data.Label = req.GetAddress().GetLabel()
			case "receiver_name":
				data.ReceiverName = req.GetAddress().GetReceiverName()
			case "receiver_phone":
				data.ReceiverPhone = req.GetAddress().GetReceiverPhone()
			case "country_code":
				data.CountryCode = req.GetAddress().GetCountryCode()
			case "province_code":
				data.ProvinceCode = req.GetAddress().GetProvinceCode()
			case "province_name":
				data.ProvinceName = req.GetAddress().GetProvinceName()
			case "city_code":
				data.CityCode = req.GetAddress().GetCityCode()
			case "city_name":
				data.CityName = req.GetAddress().GetCityName()
			case "district_code":
				data.DistrictCode = req.GetAddress().GetDistrictCode()
			case "district_name":
				data.DistrictName = req.GetAddress().GetDistrictName()
			case "street":
				data.Street = req.GetAddress().GetStreet()
			case "detail":
				data.Detail = req.GetAddress().GetDetail()
			case "postal_code":
				data.PostalCode = req.GetAddress().GetPostalCode()
			case "latitude":
				data.Latitude = req.GetAddress().GetLatitude()
			case "longitude":
				data.Longitude = req.GetAddress().GetLongitude()
			case "ext":
				extStr, extErr := marshalExt(req.GetAddress().GetExt())
				if extErr != nil {
					return extErr
				}
				data.Ext = extStr
			}
		}

		insertRes, insertErr := tx.Model(dao.UserAddress.Table()).Data(data).Insert()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "insert replacement address failed")
		}
		lastID, idErr := insertRes.LastInsertId()
		if idErr != nil {
			return gerror.Wrap(idErr, "read replacement address id failed")
		}
		newAddressID = uint64(lastID)

		// 将源地址标记为 REPLACED。
		_, updateErr := tx.Model(dao.UserAddress.Table()).
			Where(addrCols.AddressId, source.AddressId).
			Where(addrCols.UserId, userID).
			Data(do.UserAddress{
				Status:         consts.AddressStatusReplaced,
				IsDefault:      0,
				AddressVersion: gdb.Raw(addrCols.AddressVersion + " + 1"),
			}).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "retire source address failed")
		}

		shouldSetDefault := req.GetSetAsDefault() || source.IsDefault == 1
		if shouldSetDefault {
			if setErr := setDefaultAddressTx(ctx, tx, userID, newAddressID); setErr != nil {
				return setErr
			}
		} else {
			// 保持旧默认地址被替换但未重新分配的兼容性投影。
			if source.IsDefault == 1 {
				_, clearErr := tx.Model(dao.UserProfile.Table()).
					Where(profileCols.UserId, userID).
					Data(do.UserProfile{DefaultAddressId: gdb.Raw("NULL")}).
					Update()
				if clearErr != nil {
					return gerror.Wrap(clearErr, "clear profile default address id failed")
				}
			}
		}

		// 读取新地址的快照。
		if scanErr := tx.Model(dao.UserAddress.Table()).
			Where(addrCols.AddressId, newAddressID).
			Where(addrCols.UserId, userID).
			Scan(&newAddressEntity); scanErr != nil {
			return gerror.Wrap(scanErr, "query replacement address failed")
		}

		addressBookVersion, err = getAddressBookVersionTx(ctx, tx, userID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &v1.ReplaceMyAddressRes{
		SourceAddressId:    req.GetSourceAddressId(),
		NewAddress:         toProtoAddress(&newAddressEntity),
		AddressBookVersion: addressBookVersion,
	}, nil
}

// GetAddressSnapshotById 返回用于订单冗余的地址快照。
func (s *sUserProfile) GetAddressSnapshotById(ctx context.Context, req *v1.GetAddressSnapshotByIdReq) (res *v1.GetAddressSnapshotByIdRes, err error) {
	if req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "address_id is required")
	}
	item, err := s.getAddressEntityByID(ctx, req.GetAddressId())
	if err != nil {
		return nil, err
	}
	return &v1.GetAddressSnapshotByIdRes{
		Snapshot: &v1.AddressSnapshot{
			SourceAddressId:      item.AddressId,
			SourceAddressVersion: item.AddressVersion,
			UserId:               item.UserId,
			ReceiverName:         item.ReceiverName,
			ReceiverPhone:        item.ReceiverPhone,
			CountryCode:          item.CountryCode,
			ProvinceCode:         item.ProvinceCode,
			ProvinceName:         item.ProvinceName,
			CityCode:             item.CityCode,
			CityName:             item.CityName,
			DistrictCode:         item.DistrictCode,
			DistrictName:         item.DistrictName,
			Street:               item.Street,
			Detail:               item.Detail,
			PostalCode:           item.PostalCode,
			Latitude:             item.Latitude,
			Longitude:            item.Longitude,
			SnapshotAt:           timestamppb.Now(),
		},
	}, nil
}

// ResolveUserIdByPhone 根据手机号尝试反查 user_id（当前实现基于地址表 receiver_phone）。
// 返回时默认只给掩码手机号，避免明文泄露。
func (s *sUserProfile) ResolveUserIdByPhone(ctx context.Context, req *v1.ResolveUserIdByPhoneReq) (res *v1.ResolveUserIdByPhoneRes, err error) {
	phone := normalizePhoneLike(req.GetPhone())
	if phone == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "phone is required")
	}

	cols := dao.UserAddress.Columns()
	var addr entity.UserAddress
	if scanErr := dao.UserAddress.Ctx(ctx).
		Where(cols.ReceiverPhone, phone).
		Where(cols.Status, consts.AddressStatusActive).
		OrderDesc(cols.UpdatedAt).
		OrderDesc(cols.AddressId).
		Scan(&addr); scanErr != nil {
		return nil, gerror.Wrap(scanErr, "query address by phone failed")
	}
	if addr.UserId == 0 {
		return &v1.ResolveUserIdByPhoneRes{
			Found:       false,
			PhoneMasked: maskPhone(phone),
		}, nil
	}

	out := &v1.ResolveUserIdByPhoneRes{
		Found:       true,
		UserId:      addr.UserId,
		PhoneMasked: maskPhone(phone),
	}
	if req.GetIncludeProfileSummary() {
		profile, getErr := s.getProfileEntity(ctx, addr.UserId)
		if getErr != nil {
			return nil, getErr
		}
		out.ProfileSummary = &v1.UserProfileSummary{
			UserId:      profile.UserId,
			DisplayName: profile.DisplayName,
			Avatar: &v1.ImageRef{
				AssetId: profile.AvatarAssetId,
				Url:     profile.AvatarUrl,
			},
		}
	}
	return out, nil
}

// ensureProfile 确保 user_profile 行存在，避免后续更新空记录。
func (s *sUserProfile) ensureProfile(ctx context.Context, userID uint64) error {
	_, err := dao.UserProfile.Ctx(ctx).Data(do.UserProfile{
		UserId:            userID,
		ProfileVersion:    1,
		DisplayNameSource: consts.DisplayNameSourceSystemInit,
	}).InsertIgnore()
	if err != nil {
		return gerror.Wrap(err, "ensure profile row failed")
	}
	return nil
}

// getProfileEntity 查询用户资料实体。
func (s *sUserProfile) getProfileEntity(ctx context.Context, userID uint64) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	cols := dao.UserProfile.Columns()
	if err := dao.UserProfile.Ctx(ctx).Where(cols.UserId, userID).Scan(&profile); err != nil {
		return nil, gerror.Wrap(err, "query profile failed")
	}
	if profile.UserId == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "profile not found")
	}
	return &profile, nil
}

// listAddressEntities 查询地址列表。
func (s *sUserProfile) listAddressEntities(
	ctx context.Context,
	userID uint64,
	includeDeleted bool,
	page int,
	pageSize int,
) ([]entity.UserAddress, error) {
	cols := dao.UserAddress.Columns()
	model := dao.UserAddress.Ctx(ctx).Where(cols.UserId, userID)
	if !includeDeleted {
		model = model.Where(cols.Status, consts.AddressStatusActive)
	}

	var addresses []entity.UserAddress
	err := model.
		Page(page, pageSize).
		OrderDesc(cols.IsDefault).
		OrderDesc(cols.UpdatedAt).
		OrderDesc(cols.AddressId).
		Scan(&addresses)
	if err != nil {
		return nil, gerror.Wrap(err, "query addresses failed")
	}
	return addresses, nil
}

// getMyAddressEntity 查询当前用户名下的一条地址。
func (s *sUserProfile) getMyAddressEntity(ctx context.Context, userID uint64, addressID uint64) (*entity.UserAddress, error) {
	var item entity.UserAddress
	cols := dao.UserAddress.Columns()
	if err := dao.UserAddress.Ctx(ctx).
		Where(cols.AddressId, addressID).
		Where(cols.UserId, userID).
		Scan(&item); err != nil {
		return nil, gerror.Wrap(err, "query address failed")
	}
	if item.AddressId == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "address not found")
	}
	return &item, nil
}

// getAddressEntityByID 按 address_id 查询地址。
func (s *sUserProfile) getAddressEntityByID(ctx context.Context, addressID uint64) (*entity.UserAddress, error) {
	var item entity.UserAddress
	cols := dao.UserAddress.Columns()
	if err := dao.UserAddress.Ctx(ctx).Where(cols.AddressId, addressID).Scan(&item); err != nil {
		return nil, gerror.Wrap(err, "query address failed")
	}
	if item.AddressId == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "address not found")
	}
	return &item, nil
}

// userIDFromContext 从上下文提取 user_id。
// 优先读业务 context（HTTP 中间件注入），其次读 gRPC metadata。
func userIDFromContext(ctx context.Context) (uint64, error) {
	// 优先读取业务上下文值（HTTP 中间件可注入）。
	if v := ctx.Value(ctxkey.UserIDKey{}); v != nil {
		userID := gconv.Uint64(v)
		if userID > 0 {
			return userID, nil
		}
	}
	// 回退到 gRPC metadata 以兼容 RPC 场景。
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"user_id", "userid", "x-user-id", "uid"} {
		if v := md.Get(key); v != nil {
			userID := gconv.Uint64(v)
			if userID > 0 {
				return userID, nil
			}
		}
	}
	return 0, gerror.NewCode(gcode.CodeInvalidParameter, "missing user_id in grpc metadata")
}

// normalizeProfileMaskPath 归一化 profile 更新路径。
func normalizeProfileMaskPath(path string) (string, error) {
	switch strings.TrimSpace(path) {
	case "display_name", "displayName":
		return "display_name", nil
	case "avatar":
		return "avatar", nil
	case "gender":
		return "gender", nil
	case "birthday":
		return "birthday", nil
	case "locale":
		return "locale", nil
	case "timezone":
		return "timezone", nil
	case "marketing_opt_in", "marketingOptIn":
		return "marketing_opt_in", nil
	case "ext":
		return "ext", nil
	default:
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "unsupported update path: "+path)
	}
}

// normalizeAddressMaskPath 归一化 address 更新路径。
func normalizeAddressMaskPath(path string) (string, error) {
	switch strings.TrimSpace(path) {
	case "label":
		return "label", nil
	case "receiver_name", "receiverName":
		return "receiver_name", nil
	case "receiver_phone", "receiverPhone":
		return "receiver_phone", nil
	case "country_code", "countryCode":
		return "country_code", nil
	case "province_code", "provinceCode":
		return "province_code", nil
	case "province_name", "provinceName":
		return "province_name", nil
	case "city_code", "cityCode":
		return "city_code", nil
	case "city_name", "cityName":
		return "city_name", nil
	case "district_code", "districtCode":
		return "district_code", nil
	case "district_name", "districtName":
		return "district_name", nil
	case "street":
		return "street", nil
	case "detail":
		return "detail", nil
	case "postal_code", "postalCode":
		return "postal_code", nil
	case "latitude":
		return "latitude", nil
	case "longitude":
		return "longitude", nil
	case "ext":
		return "ext", nil
	default:
		return "", gerror.NewCode(gcode.CodeInvalidParameter, "unsupported update path: "+path)
	}
}

// toProtoProfile 把 UserProfile 实体映射为对外返回结构。
func toProtoProfile(p *entity.UserProfile) *v1.UserProfile {
	if p == nil {
		return nil
	}
	out := &v1.UserProfile{
		UserId:           p.UserId,
		DisplayName:      p.DisplayName,
		Avatar:           &v1.ImageRef{AssetId: p.AvatarAssetId, Url: p.AvatarUrl},
		Gender:           v1.Gender(p.Gender),
		Locale:           p.Locale,
		Timezone:         p.Timezone,
		MarketingOptIn:   p.MarketingOptIn == 1,
		DefaultAddressId: p.DefaultAddressId,
		ProfileVersion:   p.ProfileVersion,
		Ext:              parseExt(p.Ext),
		CreatedAt:        toProtoTimestamp(p.CreatedAt),
		UpdatedAt:        toProtoTimestamp(p.UpdatedAt),
	}
	if p.Birthday != nil && !p.Birthday.IsZero() {
		out.Birthday = p.Birthday.Layout("2006-01-02")
	}
	return out
}

// toProtoAddressList 批量转换地址实体切片。
func toProtoAddressList(items []entity.UserAddress) []*v1.UserAddress {
	if len(items) == 0 {
		return nil
	}
	out := make([]*v1.UserAddress, 0, len(items))
	for i := range items {
		out = append(out, toProtoAddress(&items[i]))
	}
	return out
}

// toProtoAddress 把 UserAddress 实体映射为对外返回结构。
func toProtoAddress(a *entity.UserAddress) *v1.UserAddress {
	if a == nil {
		return nil
	}
	return &v1.UserAddress{
		AddressId:             a.AddressId,
		UserId:                a.UserId,
		Status:                v1.AddressStatus(a.Status),
		Label:                 a.Label,
		ReceiverName:          a.ReceiverName,
		ReceiverPhone:         a.ReceiverPhone,
		CountryCode:           a.CountryCode,
		ProvinceCode:          a.ProvinceCode,
		ProvinceName:          a.ProvinceName,
		CityCode:              a.CityCode,
		CityName:              a.CityName,
		DistrictCode:          a.DistrictCode,
		DistrictName:          a.DistrictName,
		Street:                a.Street,
		Detail:                a.Detail,
		PostalCode:            a.PostalCode,
		IsDefault:             a.IsDefault == 1,
		Latitude:              a.Latitude,
		Longitude:             a.Longitude,
		AddressVersion:        a.AddressVersion,
		ReplacedFromAddressId: a.ReplacedFromAddressId,
		Ext:                   parseExt(a.Ext),
		CreatedAt:             toProtoTimestamp(a.CreatedAt),
		UpdatedAt:             toProtoTimestamp(a.UpdatedAt),
	}
}

// marshalExt 将扩展字段 map 序列化为 JSON 字符串。
func marshalExt(ext map[string]string) (string, error) {
	if len(ext) == 0 {
		return "{}", nil
	}
	extJSON, err := json.Marshal(ext)
	if err != nil {
		return "", gerror.WrapCode(gcode.CodeInvalidParameter, err, "ext marshal failed")
	}
	return string(extJSON), nil
}

// parseExt 解析 JSON 扩展字段，解析失败时返回 nil 以保证兼容旧数据。
func parseExt(raw string) map[string]string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	ext := make(map[string]string)
	if err := json.Unmarshal([]byte(raw), &ext); err != nil {
		return nil
	}
	if len(ext) == 0 {
		return nil
	}
	return ext
}

// toProtoTimestamp 把 gtime 转成 protobuf timestamp。
func toProtoTimestamp(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(time.Unix(0, t.TimestampNano()))
}

// setDefaultAddressTx 在事务中设置默认地址并同步地址簿版本号。
func setDefaultAddressTx(ctx context.Context, tx gdb.TX, userID, addressID uint64) error {
	// 验证目标地址存在且处于活跃状态。
	addrCols := dao.UserAddress.Columns()
	count, err := tx.Model(dao.UserAddress.Table()).
		Where(addrCols.AddressId, addressID).
		Where(addrCols.UserId, userID).
		Where(addrCols.Status, consts.AddressStatusActive).
		Count()
	if err != nil {
		return gerror.Wrap(err, "check address failed")
	}
	if count == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "address not found")
	}

	// 清除所有地址的默认标记。
	_, err = tx.Model(dao.UserAddress.Table()).
		Where(addrCols.UserId, userID).
		Where(addrCols.Status, consts.AddressStatusActive).
		Data(do.UserAddress{IsDefault: 0}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "clear default address flags failed")
	}
	// 标记目标地址为默认。
	_, err = tx.Model(dao.UserAddress.Table()).
		Where(addrCols.AddressId, addressID).
		Where(addrCols.UserId, userID).
		Data(do.UserAddress{IsDefault: 1}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "set default address flag failed")
	}

	// 同步画像的 default_address_id 与地址簿版本。
	profileCols := dao.UserProfile.Columns()
	_, err = tx.Model(dao.UserProfile.Table()).
		Where(profileCols.UserId, userID).
		Data(do.UserProfile{
			DefaultAddressId:   addressID,
			AddressBookVersion: gdb.Raw(profileCols.AddressBookVersion + " + 1"),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "set profile default address id failed")
	}
	return nil
}

// clearDefaultAddressTx 在事务里清除默认地址标记，并推进 address_book_version。
func clearDefaultAddressTx(ctx context.Context, tx gdb.TX, userID uint64) error {
	addrCols := dao.UserAddress.Columns()
	_, err := tx.Model(dao.UserAddress.Table()).
		Where(addrCols.UserId, userID).
		Where(addrCols.Status, consts.AddressStatusActive).
		Data(do.UserAddress{IsDefault: gdb.Raw("NULL")}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "clear default address flags failed")
	}
	profileCols := dao.UserProfile.Columns()
	_, err = tx.Model(dao.UserProfile.Table()).
		Where(profileCols.UserId, userID).
		Data(do.UserProfile{
			DefaultAddressId:   gdb.Raw("NULL"),
			AddressBookVersion: gdb.Raw(profileCols.AddressBookVersion + " + 1"),
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "clear profile default address id failed")
	}
	return nil
}

// getAddressBookVersionTx 读取地址簿版本号。
func getAddressBookVersionTx(ctx context.Context, tx gdb.TX, userID uint64) (uint64, error) {
	var profile entity.UserProfile
	cols := dao.UserProfile.Columns()
	if err := tx.Model(dao.UserProfile.Table()).
		Where(cols.UserId, userID).
		Scan(&profile); err != nil {
		return 0, gerror.Wrap(err, "query profile version failed")
	}
	return profile.AddressBookVersion, nil
}

// normalizePhoneLike 归一化手机号字符串，便于匹配。
func normalizePhoneLike(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	raw = strings.ReplaceAll(raw, " ", "")
	raw = strings.ReplaceAll(raw, "-", "")
	raw = strings.ReplaceAll(raw, "(", "")
	raw = strings.ReplaceAll(raw, ")", "")
	return raw
}

// maskPhone 对手机号做脱敏展示。
func maskPhone(phone string) string {
	phone = normalizePhoneLike(phone)
	n := len(phone)
	if n <= 4 {
		return "****"
	}
	if n <= 7 {
		return phone[:1] + "****" + phone[n-1:]
	}
	return phone[:3] + "****" + phone[n-4:]
}

// toInterfaceSliceUint64 把 []uint64 转为 ORM 的 []any 入参。
func toInterfaceSliceUint64(ids []uint64) []any {
	out := make([]any, 0, len(ids))
	for _, id := range ids {
		out = append(out, id)
	}
	return out
}

// formatMaskPaths 用于日志打印 FieldMask 路径。
func formatMaskPaths(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", paths)
}
