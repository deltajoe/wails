package build

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v2/internal/project"
	"github.com/wailsapp/wails/v2/internal/staticanalysis"
)

func TestSyncAssetsToEmbedPathsCopiesFromCustomDir(t *testing.T) {
	tmpDir := t.TempDir()

	customFrontend := filepath.Join(tmpDir, "web")
	customDist := filepath.Join(customFrontend, "dist")
	if err := os.MkdirAll(customDist, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(customDist, "index.html"), []byte("<h1>test</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}

	defaultDist := filepath.Join(tmpDir, "frontend", "dist")
	if err := os.MkdirAll(defaultDist, 0o755); err != nil {
		t.Fatal(err)
	}

	embedDetails := []*staticanalysis.EmbedDetails{
		{BaseDir: tmpDir, EmbedPath: "frontend/dist", All: true},
	}

	if err := syncAssetsToEmbedPaths(customFrontend, tmpDir, embedDetails); err != nil {
		t.Fatalf("syncAssetsToEmbedPaths failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(defaultDist, "index.html"))
	if err != nil {
		t.Fatalf("expected index.html in default embed dir: %v", err)
	}
	if string(data) != "<h1>test</h1>" {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestSyncFrontendAssetsNoOpForDefaultDir(t *testing.T) {
	tmpDir := t.TempDir()

	frontendDist := filepath.Join(tmpDir, "frontend", "dist")
	if err := os.MkdirAll(frontendDist, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontendDist, "index.html"), []byte("<h1>original</h1>"), 0o644); err != nil {
		t.Fatal(err)
	}

	opts := &Options{
		ProjectData: &project.Project{
			Path:        tmpDir,
			FrontendDir: "frontend",
		},
	}

	if err := syncFrontendAssetsToEmbedDir(opts); err != nil {
		t.Fatalf("syncFrontendAssetsToEmbedDir failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(frontendDist, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "<h1>original</h1>" {
		t.Errorf("content should be unchanged for default frontend dir")
	}
}

func TestSyncFrontendAssetsNoOpForNilProject(t *testing.T) {
	opts := &Options{}
	if err := syncFrontendAssetsToEmbedDir(opts); err != nil {
		t.Errorf("should be no-op for nil project data: %v", err)
	}
}

func TestCopyDir(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	subDir := filepath.Join(srcDir, "assets")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "app.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "index.html"), []byte("<html></html>"), 0o644); err != nil {
		t.Fatal(err)
	}

	targetDir := filepath.Join(dstDir, "output")
	if err := copyDir(srcDir, targetDir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(targetDir, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "<html></html>" {
		t.Errorf("unexpected index.html content: %s", data)
	}

	data, err = os.ReadFile(filepath.Join(targetDir, "assets", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "console.log(1)" {
		t.Errorf("unexpected app.js content: %s", data)
	}
}
