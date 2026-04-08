package media

import (
	"testing"

	"github.com/TsingpekTao/shopa/media-svc/internal/consts"
	"github.com/TsingpekTao/shopa/media-svc/internal/model/entity"
)

func TestDefaultScenePolicyForMissingSceneBuyerAvatarUsesPrivateACL(t *testing.T) {
	policy, ok := defaultScenePolicyForMissingScene("buyer_avatar")
	if !ok {
		t.Fatalf("expected buyer_avatar fallback policy to exist")
	}
	if policy.AclType != consts.ACLTypePrivate {
		t.Fatalf("expected buyer_avatar fallback acl_type=%d, got %d", consts.ACLTypePrivate, policy.AclType)
	}
}

func TestShouldPreferSignedReadURLUsesSignedURLForSellerMediaOSSAssets(t *testing.T) {
	asset := &entity.MediaAsset{
		AssetId:         32,
		SceneCode:       "seller_media",
		AclType:         consts.ACLTypePublicRead,
		StorageProvider: "oss",
		Bucket:          "shopa-catalog-1",
		ObjectKey:       "seller_media/example.png",
	}

	if !shouldPreferSignedReadURL(asset) {
		t.Fatalf("expected seller_media oss asset to prefer signed read url")
	}
}

func TestShouldPreferSignedReadURLKeepsPublicURLForOtherPublicOSSAssets(t *testing.T) {
	asset := &entity.MediaAsset{
		AssetId:         46,
		SceneCode:       "marketing_banner",
		AclType:         consts.ACLTypePublicRead,
		StorageProvider: "oss",
		Bucket:          "shopa-public-1",
		ObjectKey:       "marketing/banner.png",
	}

	if shouldPreferSignedReadURL(asset) {
		t.Fatalf("expected non-seller_media oss asset to keep public read url")
	}
}

func TestBuildOSSPublicURLEncodesObjectKeyPath(t *testing.T) {
	url := buildOSSPublicURL(
		"oss-cn-beijing.aliyuncs.com",
		"shopa-catalog-1",
		"seller_media/20260403/example 图片.png",
	)

	expected := "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/seller_media/20260403/example%20%E5%9B%BE%E7%89%87.png"
	if url != expected {
		t.Fatalf("expected encoded oss public url %q, got %q", expected, url)
	}
}
