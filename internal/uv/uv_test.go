package uv

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetUvDownloadURLLinuxWithMirror(t *testing.T) {
	arch := mapGoArchToUVArch(runtime.GOARCH)
	expected := "https://pypi.tuna.tsinghua.edu.cn/uv/releases/latest/download/uv-" + arch + "-unknown-linux-gnu.tar.gz"
	if got := getUvDownloadURL("https://pypi.tuna.tsinghua.edu.cn/"); got != expected {
		t.Fatalf("unexpected mirror URL: %s", got)
	}
}

func TestGetUvDownloadURLLinuxUsesOfficialForEmptyMirror(t *testing.T) {
	arch := mapGoArchToUVArch(runtime.GOARCH)
	expected := "https://github.com/astral-sh/uv/releases/latest/download/uv-" + arch + "-unknown-linux-gnu.tar.gz"
	if got := getUvDownloadURL(""); got != expected {
		t.Fatalf("unexpected official URL: %s", got)
	}
}

func TestGetUvDownloadURLLinuxIgnoresAstralShMirror(t *testing.T) {
	arch := mapGoArchToUVArch(runtime.GOARCH)
	expected := "https://github.com/astral-sh/uv/releases/latest/download/uv-" + arch + "-unknown-linux-gnu.tar.gz"
	got := getUvDownloadURL("https://astral.sh")
	if !strings.Contains(got, expected) {
		t.Fatalf("expected official mirror fallback, got %s", got)
	}
}
