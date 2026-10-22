package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixtureRoot 是一个「已发版且三件套齐」的最小样本仓：
// 矩阵有 v0.1 那一行、CHANGELOG 有 v0.1.0 那一节、tag 列表由用例注入。
const fixtureRoot = "testdata/ok"

// baselineRoot 是本仓今天的真实形状：零发版、零 tag、矩阵只有基线那一行。
// 它必须绿——否则护栏会把「还没发版」误判成「纪律没遵守」。
const baselineRoot = "testdata/baseline"

func check(root string, tags []string) (string, error) {
	var buf bytes.Buffer
	err := run(options{
		root:      root,
		readme:    "README.md",
		changelog: "CHANGELOG.md",
		agents:    "AGENTS.md",
		listTags:  func() ([]string, error) { return tags, nil },
		out:       &buf,
	})
	return buf.String(), err
}

func checkFixture(t *testing.T, tags []string) (string, error) {
	t.Helper()
	return check(fixtureRoot, tags)
}

func TestConsistentFixtureIsGreen(t *testing.T) {
	out, err := checkFixture(t, []string{"v0.1.0"})
	if err != nil {
		t.Fatalf("一致的夹具必须绿，got %v\n%s", err, out)
	}
	if !strings.Contains(out, "✅") {
		t.Errorf("绿时要有结论行，got:\n%s", out)
	}
}

func TestBaselineWithoutReleaseIsGreen(t *testing.T) {
	out, err := check(baselineRoot, nil)
	if err != nil {
		t.Fatalf("尚未发版的基线必须绿（零发版不是违规），got %v\n%s", err, out)
	}
}

// TestReleaseWithoutMatrixRowIsRed 是 KT-02 ⑤ 的破坏验证本体：
// 发一版（CHANGELOG 补一节）但忘了更新兼容矩阵 → 必须红，且点名缺哪个序列。
func TestReleaseWithoutMatrixRowIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "CHANGELOG.md"),
		changelogOf(t, root)+"\n## [v0.2.0] - 2026-12-01\n\n- 新增 `statistics/Quantile`。\n")

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("发版没更新矩阵必须红\n%s", out)
	}
	if !strings.Contains(out, "CHANGELOG 有 v0.2.0 这一版，矩阵却没有 v0.2 那一行") {
		t.Errorf("红因要点名缺失的矩阵行，got:\n%s", out)
	}
}

func TestMatrixRowWithoutChangelogIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")),
			"<!-- matrix:end -->",
			"| `v0.9` | 已发版 | — | — | — |\n<!-- matrix:end -->", 1))

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("矩阵凭空多一行必须红，got\n%s", out)
	}
	if !strings.Contains(out, "矩阵声明了序列 v0.9") {
		t.Errorf("红因要点名多余序列 v0.9，got:\n%s", out)
	}
}

func TestTagWithoutChangelogSectionIsRed(t *testing.T) {
	out, err := checkFixture(t, []string{"v0.1.0", "v0.3.0"})
	if err == nil {
		t.Fatalf("打了 tag 却没写账必须红，got\n%s", out)
	}
	if !strings.Contains(out, "tag v0.3.0 在 CHANGELOG 里没有对应版本节") {
		t.Errorf("红因要点名 tag v0.3.0，got:\n%s", out)
	}
}

func TestTagWithoutMatrixRowIsRed(t *testing.T) {
	out, err := checkFixture(t, []string{"v0.1.0", "v0.2.0"})
	if err == nil {
		t.Fatalf("tag 落在矩阵没声明的序列上必须红，got\n%s", out)
	}
	for _, want := range []string{"tag v0.2.0 的版本序列 v0.2 不在兼容矩阵里", "tag v0.2.0 在 CHANGELOG 里没有对应版本节"} {
		if !strings.Contains(out, want) {
			t.Errorf("红因缺 %q，got:\n%s", want, out)
		}
	}
}

func TestMissingMatrixMarkersIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")), matrixBegin, "（矩阵在下面）", 1))

	if _, err := check(root, []string{"v0.1.0"}); err == nil {
		t.Fatal("删掉标记来「让检查通过」必须红")
	}
}

func TestEmptyMatrixDataRowsIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")),
			"| `main`（未发版） | 开发基线，不做兼容承诺 | 任意 | 任意 | 任意 |\n| `v0.1` | 已发版 | `v0.1.0` | `v0.1.0` | `v0.1.0` |\n", "", 1))

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("把矩阵行清空必须红，got\n%s", out)
	}
	if !strings.Contains(err.Error(), "没有任何数据行") {
		t.Errorf("错误要指出矩阵没有数据行，got %v", err)
	}
}

// TestMissingTableSeparatorIsFatal 钉住「表头约定」：分隔行没了就无法区分表头与数据，
// 删表头会让第一行数据被当表头跳过、护栏悄悄变松——这条路要堵死。
func TestMissingTableSeparatorIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")), "|---|---|---|---|---|\n", "", 1))

	if _, err := check(root, []string{"v0.1.0"}); err == nil {
		t.Fatal("矩阵缺 |---| 分隔行必须红")
	}
}

func TestMissingRuleRowIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")),
			"| `PATCH` | 文档 / 注释 / 无行为修正 | 措辞订正 |\n", "", 1))

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("规则表少一行必须红，got\n%s", out)
	}
	if !strings.Contains(out, "PATCH") {
		t.Errorf("红因要点名缺失的 PATCH 行，got:\n%s", out)
	}
}

func TestGoprivatePrerequisiteMissingIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "AGENTS.md"), "# AGENTS（夹具）\n\n（前提那行被我删了）\n")

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("私有化前提缺失必须红，got\n%s", out)
	}
	if !strings.Contains(out, goprivateMark) {
		t.Errorf("红因要指出缺 GOPRIVATE 前提，got:\n%s", out)
	}
}

func TestChangelogWithoutEntriesIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "CHANGELOG.md"), "# CHANGELOG（夹具）\n\n## [Unreleased]\n")

	out, err := check(root, nil)
	if err == nil {
		t.Fatalf("CHANGELOG 空账必须红（⑤ 要求 ≥ 1 条），got\n%s", out)
	}
	if !strings.Contains(out, "一条条目都没有") {
		t.Errorf("红因要点名空账，got:\n%s", out)
	}
}

func TestReleaseTitleWithoutDateIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "CHANGELOG.md"),
		strings.Replace(changelogOf(t, root), "## [v0.1.0] - 2026-11-01", "## [v0.1.0]", 1))

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("版本节缺日期必须红，got\n%s", out)
	}
	if !strings.Contains(out, "不合规范") {
		t.Errorf("红因要指出标题格式，got:\n%s", out)
	}
}

func TestMalformedMatrixCellIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")),
			"| `v0.1` | 已发版", "| `v0.1.0-beta` | 已发版", 1))

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("矩阵首格形状打错必须红，got\n%s", out)
	}
	if !strings.Contains(out, "v0.1.0-beta") {
		t.Errorf("红因要点名打错的首格，got:\n%s", out)
	}
}

// TestDuplicateMatrixSeriesIsRed：v0.1 与 v0.1.0 是同一个序列，写两行会让「更新哪一行」重新变模糊。
func TestDuplicateMatrixSeriesIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, filepath.Join(root, "README.md"),
		strings.Replace(readFile(t, filepath.Join(root, "README.md")),
			"<!-- matrix:end -->", "| `v0.1.0` | 已发版 | `v0.1.0` | `v0.1.0` | `v0.1.0` |\n<!-- matrix:end -->", 1))

	out, err := check(root, []string{"v0.1.0"})
	if err == nil {
		t.Fatalf("同一序列两行必须红，got\n%s", out)
	}
	if !strings.Contains(out, "序列 v0.1 有多行") {
		t.Errorf("红因要点名重复序列，got:\n%s", out)
	}
}

// TestListTagsFailureIsFatal 钉住「读不到 tag 就红」：
// 静默按零 tag 走下去，CI 少 fetch-depth 时会永远绿，那是假护栏。
func TestListTagsFailureIsFatal(t *testing.T) {
	err := run(options{
		root:      fixtureRoot,
		readme:    "README.md",
		changelog: "CHANGELOG.md",
		agents:    "AGENTS.md",
		listTags:  func() ([]string, error) { return nil, errors.New("no git history") },
		out:       io.Discard,
	})
	if err == nil {
		t.Fatal("读 tag 失败必须返回错误")
	}
	if !strings.Contains(err.Error(), "fetch-depth") {
		t.Errorf("错误要指出 CI 取 tag 的前提，got %v", err)
	}
}

func changelogOf(t *testing.T, root string) string {
	t.Helper()
	return strings.TrimRight(readFile(t, filepath.Join(root, "CHANGELOG.md")), "\n")
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func copyFixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	if err := copyTree(fixtureRoot, dst); err != nil {
		t.Fatal(err)
	}
	return dst
}

func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, raw, 0o600)
	})
}
