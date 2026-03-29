package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type productSpec struct {
	SpuNo     string
	SourceURL string
}

func main() {
	endpoint := getenvOr("OSS_ENDPOINT", "oss-cn-beijing.aliyuncs.com")
	accessKeyID := strings.TrimSpace(os.Getenv("OSS_ACCESS_KEY_ID"))
	accessKeySecret := strings.TrimSpace(os.Getenv("OSS_ACCESS_KEY_SECRET"))
	if accessKeyID == "" || accessKeySecret == "" {
		panic("missing OSS_ACCESS_KEY_ID or OSS_ACCESS_KEY_SECRET")
	}

	newBucketName := getenvOr("OSS_BUCKET", "shopa-catalog-1")
	keyPrefix := "mall/products/20260329/"

	products := []productSpec{
		{SpuNo: "P2001", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/4/4d/Puma_Jago_Zig_Zag_Running_Shoe.jpg"},
		{SpuNo: "P2002", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/4/4c/Jabra_True_Wireless_Earbuds_Lineup_%28Sep_2023%29.jpg"},
		{SpuNo: "P2003", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/0/0c/Air_Fryer_5458.jpg"},
		{SpuNo: "P2004", SourceURL: "https://live.staticflickr.com/2933/14594868619_18f368e83a.jpg"},
		{SpuNo: "P2005", SourceURL: "https://live.staticflickr.com/65535/48927346277_581e257742_b.jpg"},
		{SpuNo: "P2006", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/4/4a/Silicon_vs_GaN_30W_USB-C_chargers.jpg"},
		{SpuNo: "P2007", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/9/9f/Tasty_Coffee_pourover_V60_bloom_2025.jpg"},
		{SpuNo: "P2008", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/a/a3/Leather_handbag_by_Les_cuirs_d%27Agathe_%28DSC07738%29.jpg"},
		{SpuNo: "P2009", SourceURL: "https://live.staticflickr.com/4850/30787534027_4657683ed9_b.jpg"},
		{SpuNo: "P2010", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/2/2c/Digital_Body_Scale.jpg"},
		{SpuNo: "P2011", SourceURL: "https://upload.wikimedia.org/wikipedia/commons/b/b0/Terracotta_Polo_Ralph_Lauren_hooded_Jacket.jpg"},
		{SpuNo: "P2012", SourceURL: "https://live.staticflickr.com/5263/5756031128_d157bbaed7_b.jpg"},
		{SpuNo: "P2013", SourceURL: "https://live.staticflickr.com/3779/9300690265_eb43cfd164_b.jpg"},
		{SpuNo: "P2014", SourceURL: "https://live.staticflickr.com/3641/3541401244_95232bc83d_b.jpg"},
		{SpuNo: "P2015", SourceURL: "https://live.staticflickr.com/2695/4152119647_0f82b054cf_b.jpg"},
	}

	httpClient := &http.Client{
		Timeout: 70 * time.Second,
	}

	ossClient, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	must(err)

	newBucket, err := ossClient.Bucket(newBucketName)
	must(err)

	urlMap := make(map[string]string, len(products))
	for idx, p := range products {
		fmt.Printf("[%d/%d] preparing %s\n", idx+1, len(products), p.SpuNo)
		img, err := fetchImage(httpClient, p.SourceURL)
		must(err)
		fmt.Printf("picked source for %s: %s\n", p.SpuNo, p.SourceURL)

		cropped := cropCenterSquare(img)
		jpgBytes, err := toJPEG(cropped, 90)
		must(err)

		objectKey := keyPrefix + strings.ToLower(p.SpuNo) + ".jpg"
		err = newBucket.PutObject(objectKey, bytes.NewReader(jpgBytes))
		must(err)

		// 5-year signed URL for private bucket read.
		signedURL, err := newBucket.SignURL(objectKey, oss.HTTPGet, int64((24 * time.Hour * 365 * 5).Seconds()))
		must(err)
		signedURL = strings.Replace(signedURL, "http://", "https://", 1)
		urlMap[p.SpuNo] = signedURL
		fmt.Printf("uploaded %s -> %s\n", p.SpuNo, objectKey)
	}

	err = writeSQL(urlMap)
	must(err)
	fmt.Println("generated SQL: tmp_catalog_spu_image_update.sql")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func fetchImage(client *http.Client, imageURL string) (image.Image, error) {
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		f, err := os.Open(imageURL)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}
		return img, nil
	}

	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "shopa-catalog-image-migrator/1.0")
	resp, err := doRequestWithRetry(client, req, 4)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("download image status=%d", resp.StatusCode)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func doRequestWithRetry(client *http.Client, req *http.Request, attempts int) (*http.Response, error) {
	if attempts <= 0 {
		attempts = 1
	}
	var lastErr error
	for i := 1; i <= attempts; i++ {
		cloned := req.Clone(req.Context())
		resp, err := client.Do(cloned)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(time.Duration(i) * 1200 * time.Millisecond)
	}
	return nil, lastErr
}

func cropCenterSquare(src image.Image) image.Image {
	b := src.Bounds()
	w := b.Dx()
	h := b.Dy()
	side := w
	if h < side {
		side = h
	}
	startX := b.Min.X + (w-side)/2
	startY := b.Min.Y + (h-side)/2
	dst := image.NewRGBA(image.Rect(0, 0, side, side))
	draw.Draw(dst, dst.Bounds(), src, image.Pt(startX, startY), draw.Src)
	return dst
}

func toJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSQL(urlMap map[string]string) error {
	spus := make([]string, 0, len(urlMap))
	for spu := range urlMap {
		spus = append(spus, spu)
	}
	sort.Strings(spus)

	var b strings.Builder
	b.WriteString("SET NAMES utf8mb4;\n")
	for _, spu := range spus {
		u := strings.ReplaceAll(urlMap[spu], "'", "''")
		b.WriteString(fmt.Sprintf(
			"INSERT INTO catalog_spu_image (spu_no, image_url, source, is_primary, sort_order, status, created_at, updated_at)\n"+
				"VALUES ('%s','%s','oss',1,1,1,NOW(3),NOW(3))\n"+
				"ON DUPLICATE KEY UPDATE image_url=VALUES(image_url), source='oss', is_primary=1, sort_order=1, status=1, updated_at=NOW(3);\n",
			spu, u,
		))
	}
	return os.WriteFile("tmp_catalog_spu_image_update.sql", []byte(b.String()), 0644)
}

func getenvOr(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
