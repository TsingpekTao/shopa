package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

func bootstrapConfig(ctx context.Context) {
	fileAdapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return
	}

	for _, dir := range discoverConfigDirs() {
		_ = fileAdapter.AddPath(dir)
	}
	fileAdapter.SetFileName(resolveConfigFileName())

	if filePath, err := fileAdapter.GetFilePath(); err == nil && strings.TrimSpace(filePath) != "" {
		g.Log().Info(ctx, fmt.Sprintf("agent-svc config loaded from %s", filePath))
	}
}

func resolveConfigFileName() string {
	env := firstNonEmpty(
		os.Getenv("SHOPA_ENV"),
		os.Getenv("APP_ENV"),
		os.Getenv("GF_APP_ENV"),
		"local",
	)
	return fmt.Sprintf("config.%s.yaml", strings.ToLower(strings.TrimSpace(env)))
}

func discoverConfigDirs() []string {
	seen := make(map[string]struct{})
	dirs := make([]string, 0, 4)

	addDir := func(dir string) {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			return
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return
		}
		info, statErr := os.Stat(abs)
		if statErr != nil || !info.IsDir() {
			return
		}
		if _, ok := seen[abs]; ok {
			return
		}
		seen[abs] = struct{}{}
		dirs = append(dirs, abs)
	}

	addDir(filepath.Join("manifest", "config"))
	addDir(filepath.Join("agent-svc", "manifest", "config"))

	if exePath, err := os.Executable(); err == nil {
		baseDir := filepath.Dir(exePath)
		addDir(filepath.Join(baseDir, "manifest", "config"))
		addDir(filepath.Join(baseDir, "agent-svc", "manifest", "config"))
	}

	return dirs
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
