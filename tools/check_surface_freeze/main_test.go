package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	fixtureRoot = "testdata/ok"
	baselineRel = "docs/api-surface.txt"
)

func check(root string, update bool) (string, error) {
	var buf bytes.Buffer
	err := run(options{root: root, baseline: baselineRel, update: update, out: &buf})
	return buf.String(), err
}

func checkOK(t *testing.T) (string, error) {
	t.Helper()
	return check(fixtureRoot, false)
}

func TestCleanFixtureIsGreen(t *testing.T) {
	out, err := checkOK(t)
	if err != nil {
		t.Fatalf("干净夹具必须绿，got %v\n%s", err, out)
	}
	if !strings.Contains(out, "✅") {
		t.Errorf("绿时要有结论行，got:\n%s", out)
	}
}

// 夹具里的 cmd/ 是 package main：命令不在导出面里，基线不能出现它的行。
func TestMainPackageIsNotFrozen(t *testing.T) {
	out, err := checkOK(t)
	if err != nil {
		t.Fatalf("夹具必须绿，got %v\n%s", err, out)
	}
	raw, err := os.ReadFile(filepath.Join(fixtureRoot, filepath.FromSlash(baselineRel)))
	if err != nil {
		t.Fatalf("读基线: %v", err)
	}
	if strings.Contains(string(raw), "cmd/") {
		t.Errorf("package main 的符号不许进基线，got:\n%s", raw)
	}
}

func TestAddedExportedSymbolIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "demo/extra.go", "package demo\n\nfunc Unregistered() {}\n")

	out, err := check(root, false)
	if err == nil {
		t.Fatalf("新增未登记导出符号必须红，got\n%s", out)
	}
	if !strings.Contains(out, "+ demo/Unregistered()") {
		t.Errorf("红因要列新增行（+ 前缀 + 完整指纹），got:\n%s", out)
	}
}

func TestRemovedExportedSymbolIsRed(t *testing.T) {
	root := copyFixture(t)
	if err := os.Remove(filepath.Join(root, "demo", "extra.go")); err != nil {
		t.Fatalf("删夹具文件: %v", err)
	}

	out, err := check(root, false)
	if err == nil {
		t.Fatalf("删掉已登记导出符号必须红，got\n%s", out)
	}
	if !strings.Contains(out, "- demo/OldThing()") {
		t.Errorf("红因要列消失行（- 前缀 + 完整指纹），got:\n%s", out)
	}
}

func TestSignatureChangeIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "demo/demo.go", `package demo

type Widget struct {
	Name  string
	Size  int
	Extra bool
}

func (w Widget) Label() string { return w.Name }

func (w *Widget) Resize(n int) error { return nil }
`)

	out, err := check(root, false)
	if err == nil {
		t.Fatalf("改结构体字段必须红，got\n%s", out)
	}
	// 结构体形状变了：旧行消失 + 新行新增，两侧都要点出来。
	if !strings.Contains(out, "- demo/Widget struct{Name string; Size int}") {
		t.Errorf("要列出旧结构体形状，got:\n%s", out)
	}
	if !strings.Contains(out, "+ demo/Widget struct{Name string; Size int; Extra bool}") {
		t.Errorf("要列出新结构体形状，got:\n%s", out)
	}
}

func TestMissingBaselineIsRed(t *testing.T) {
	root := copyFixture(t)
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(baselineRel))); err != nil {
		t.Fatalf("删基线: %v", err)
	}

	out, err := check(root, false)
	if err == nil {
		t.Fatalf("基线缺失必须红（删基线不能等于放行），got\n%s", out)
	}
	if !strings.Contains(err.Error(), "读冻结基线") {
		t.Errorf("红因要点名基线缺失，got %v", err)
	}
}

func TestEmptyBaselineIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, baselineRel, "\n")

	out, err := check(root, false)
	if err == nil {
		t.Fatalf("空基线必须红，got\n%s", out)
	}
	if !strings.Contains(err.Error(), "一行指纹都没有") {
		t.Errorf("红因要点名空基线，got %v", err)
	}
}

func TestNoLibraryPackageIsRed(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "cmd/tool/main.go", "package main\n\nfunc main() {}\n")
	writeFile(t, root, baselineRel, "placeholder\n")

	out, err := check(root, false)
	if err == nil {
		t.Fatalf("没有任何库包必须红（拒绝空绿），got\n%s", out)
	}
	if !strings.Contains(err.Error(), "没有任何库包") {
		t.Errorf("红因要点名空扫，got %v", err)
	}
}

// -update 是「先登记再基线」那条腿：新增符号后跑一次 update，随后比对必须绿。
func TestUpdateRewritesBaseline(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "demo/extra.go", "package demo\n\nfunc Unregistered() {}\n")

	if out, err := check(root, false); err == nil {
		t.Fatalf("更新前必须红，got\n%s", out)
	}

	out, err := check(root, true)
	if err != nil {
		t.Fatalf("update 必须成功，got %v\n%s", err, out)
	}
	if !strings.Contains(out, "已更新") {
		t.Errorf("update 要有确认行，got:\n%s", out)
	}
	if out, err := check(root, false); err != nil {
		t.Fatalf("update 后必须绿，got %v\n%s", err, out)
	}
}

// copyFixture 把夹具复制到 t.TempDir()：每条用例都在自己的副本上增删符号，互不干扰。
func copyFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	err := filepath.WalkDir(fixtureRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(fixtureRoot, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(root, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
	if err != nil {
		t.Fatalf("复制夹具: %v", err)
	}
	return root
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("建目录 %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("写 %s: %v", rel, err)
	}
}
