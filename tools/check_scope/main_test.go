package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureRoot = "testdata/ok"

func check(root, scope string) (string, error) {
	var buf bytes.Buffer
	err := run(options{root: root, scope: scope, out: &buf})
	return buf.String(), err
}

func checkOK(t *testing.T) (string, error) {
	t.Helper()
	return check(fixtureRoot, "SCOPE.md")
}

// TestCleanFixtureIsGreen 同时钉住两条正腿：
// 同仓内部 import（sandbox/runner → kit/logging）与标准库 os/signal 的段，都不许被误红。
func TestCleanFixtureIsGreen(t *testing.T) {
	out, err := checkOK(t)
	if err != nil {
		t.Fatalf("干净夹具必须绿，got %v\n%s", err, out)
	}
	if !strings.Contains(out, "✅") {
		t.Errorf("绿时要有结论行，got:\n%s", out)
	}
}

func TestDeclaredTypeWithBusinessWordIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "id/ident.go", "package id\n\ntype OrderID string\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("声明业务类型必须红，got\n%s", out)
	}
	if !strings.Contains(out, `类型 "OrderID" 命中业务词 "order"`) {
		t.Errorf("红因要点名 OrderID 与词表出处，got:\n%s", out)
	}
	// 位置要指到具体文件与行，否则报错了没人能找到该改哪儿。
	if !strings.Contains(out, "id/ident.go:3") {
		t.Errorf("红因要给出文件:行，got:\n%s", out)
	}
}

func TestFunctionNameWithBusinessWordIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "observability/metrics.go",
		"package observability\n\nfunc ObserveBacktest(name string, seconds float64) {}\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("业务函数名必须红，got\n%s", out)
	}
	if !strings.Contains(out, "命中业务词 \"backtest\"") {
		t.Errorf("红因要点名 backtest，got:\n%s", out)
	}
}

func TestStructFieldWithBusinessWordIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "statistics/bar.go",
		"package statistics\n\ntype Series struct {\n\tBarCount int\n\tMean float64\n}\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("业务字段名必须红，got\n%s", out)
	}
	if !strings.Contains(out, `字段 "BarCount"`) {
		t.Errorf("红因要点名字段 BarCount，got:\n%s", out)
	}
}

func TestInterfaceMethodWithBusinessWordIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "sandbox/runner/iface.go",
		"package runner\n\ntype Emitter interface {\n\tEmitSignal(v float64) error\n}\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("接口方法名必须红，got\n%s", out)
	}
	if !strings.Contains(out, `接口方法 "EmitSignal"`) {
		t.Errorf("红因要点名接口方法，got:\n%s", out)
	}
}

func TestPackageLevelVarWithBusinessWordIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "decimal/factor.go", "package decimal\n\nvar FactorCache = map[string]int{}\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("包级变量名必须红，got\n%s", out)
	}
	if !strings.Contains(out, "命中业务词 \"factor\"") {
		t.Errorf("红因要点名 factor，got:\n%s", out)
	}
}

// TestSupplyChainCompoundIsRed 钉住相邻双词拼接那条：SupplyChain 切词是 [supply chain]，
// 单段整词比对抓不到，而把 supply 或 chain 单独列进词表会误红正当的链表类型。
func TestSupplyChainCompoundIsRed(t *testing.T) {
	if got := wordHit("SupplyChainSchemaDDL"); got != "supplychain" {
		t.Errorf("SupplyChain 必须命中 supplychain，got %q", got)
	}
	if got := wordHit("supply_chain"); got != "supplychain" {
		t.Errorf("snake 写法也必须命中，got %q", got)
	}
}

func TestChainAloneIsNotBusinessWord(t *testing.T) {
	// 反向腿：chain 单词不是业务词，否则 kit 里任何链表实现都会被拦。
	if got := wordHit("Chain"); got != "" {
		t.Errorf("Chain 不该命中词表，got %q", got)
	}
	if got := wordHit("LinkedChain"); got != "" {
		t.Errorf("LinkedChain 不该命中词表，got %q", got)
	}
}

// TestFactorialIsNotBusinessWord 钉住整词比对：factor 是词表里的词，factorial 不是。
func TestFactorialIsNotBusinessWord(t *testing.T) {
	if got := wordHit("Factorial"); got != "" {
		t.Errorf("Factorial 不该误命中 factor，got %q", got)
	}
	if got := wordHit("SampleFactorial"); got != "" {
		t.Errorf("SampleFactorial 不该误命中，got %q", got)
	}
	if got := wordHit("FactorDefinition"); got != "factor" {
		t.Errorf("FactorDefinition 必须命中 factor，got %q", got)
	}
}

func TestCoreImportIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "db/conn.go",
		"package db\n\nimport \"tellurion/pkg/storage\"\n\nvar Store = storage.NewPostgresStore\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("import core 必须红，got\n%s", out)
	}
	if !strings.Contains(out, "越界依赖：core 的 module 首段是 `tellurion`") {
		t.Errorf("红因要说明越界依赖，got:\n%s", out)
	}
}

func TestSiblingRepoImportIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "decimal/money.go",
		"package decimal\n\nimport \"github.com/ruoxizhnya/tellurion-contracts/domain\"\n\nvar Zero = domain.Decimal{}\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("反向依赖 contracts 必须红，got\n%s", out)
	}
	if !strings.Contains(out, "tellurion-contracts 是同项目兄弟仓") {
		t.Errorf("红因要点名兄弟仓，got:\n%s", out)
	}
}

func TestThirdPartyBusinessPathSegmentIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "httpclient/market.go",
		"package httpclient\n\nimport \"example.com/vendor/market/feed\"\n\nvar Topic = feed.All\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("第三方路径段带业务词也要红，got\n%s", out)
	}
	if !strings.Contains(out, "路径段命中业务词 \"market\"") {
		t.Errorf("红因要指出路径段，got:\n%s", out)
	}
}

func TestUnregisteredPackageIsRed(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "cache/cache.go", "package cache\n\nfunc Get(k string) string { return \"\" }\n")

	out, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatalf("新增没登记的包必须红，got\n%s", out)
	}
	if !strings.Contains(out, "./cache：这个包没登记进") {
		t.Errorf("红因要点名未登记目录，got:\n%s", out)
	}
}

func TestMissingScopeRowIsFatal(t *testing.T) {
	root := copyFixture(t)
	src := readFile(t, filepath.Join(fixtureRoot, "SCOPE.md"))
	writeFile(t, root, "SCOPE.md", strings.Replace(src, "| `db` | 连接 / DSN / 迁移机制 | 任何表结构 | — |\n", "", 1))

	_, err := check(root, "SCOPE.md")
	if err == nil {
		t.Fatal("清单少一行必须 fatal（否则「删一行」等于「少登记一个包」）")
	}
	if !strings.Contains(err.Error(), "少了 db") {
		t.Errorf("错误要点名缺失的 db 包，got %v", err)
	}
}

func TestMissingMarkersIsFatal(t *testing.T) {
	root := copyFixture(t)
	src := readFile(t, filepath.Join(fixtureRoot, "SCOPE.md"))
	writeFile(t, root, "SCOPE.md", strings.Replace(src, scopeBegin, "（清单在下面）", 1))

	if _, err := check(root, "SCOPE.md"); err == nil {
		t.Fatal("删掉标记来让检查通过必须 fatal")
	}
}

func TestMissingTableSeparatorIsFatal(t *testing.T) {
	root := copyFixture(t)
	src := readFile(t, filepath.Join(fixtureRoot, "SCOPE.md"))
	writeFile(t, root, "SCOPE.md", strings.Replace(src, "|---|---|---|---|\n", "", 1))

	if _, err := check(root, "SCOPE.md"); err == nil {
		t.Fatal("删掉表头分隔行必须 fatal")
	}
}

func TestMainPackageIsNotScanned(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "cmd/scratch/main.go",
		"package main\n\nfunc runBacktest() {}\n\ntype OrderRow struct{ ID string }\n")

	out, err := check(root, "SCOPE.md")
	if err != nil {
		t.Fatalf("package main 不参与核对（词表就写在工具里），got %v\n%s", err, out)
	}
}

func TestTestFilesAreNotScanned(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "logging/logger_test.go",
		"package logging\n\nimport \"testing\"\n\ntype OrderID string\n\nfunc TestX(t *testing.T) { _ = OrderID(\"\") }\n")

	out, err := check(root, "SCOPE.md")
	if err != nil {
		t.Fatalf("*_test.go 不参与核对（库的形状由非测试文件决定），got %v\n%s", err, out)
	}
}

// TestCommentsAreNotScanned 钉住 §2 登记的覆盖面限制：
// 注释要说清调用方是谁，扫注释会把文档本身变成违规物。
func TestCommentsAreNotScanned(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "clock/clock.go",
		"// Package clock 提供时间原语，供回测循环与实盘共用。\npackage clock\n\n// Mode 是时钟模式（backtest / live）。\ntype Mode string\n")

	out, err := check(root, "SCOPE.md")
	if err != nil {
		t.Fatalf("注释不参与核对，got %v\n%s", err, out)
	}
}

func TestTwoPackageNamesInOneDirIsFatal(t *testing.T) {
	root := copyFixture(t)
	writeFile(t, root, "logging/other.go", "package logx\n\nfunc Extra() {}\n")

	if _, err := check(root, "SCOPE.md"); err == nil {
		t.Fatal("一个目录两个包名必须 fatal（否则会静默少扫一半）")
	}
}

func TestTokensSplitting(t *testing.T) {
	cases := map[string][]string{
		"OrderID":              {"order", "id"},
		"ObserveBacktest":      {"observe", "backtest"},
		"MaskAPIKey":           {"mask", "api", "key"},
		"SupplyChainSchemaDDL": {"supply", "chain", "schema", "ddl"},
		"TimeoutMS":            {"timeout", "ms"},
		"Bar2":                 {"bar", "2"},
		"HTTPServer":           {"http", "server"},
	}
	for in, want := range cases {
		if got := tokens(in); strings.Join(got, " ") != strings.Join(want, " ") {
			t.Errorf("tokens(%q) = %v，期望 %v", in, got, want)
		}
	}
}

func TestStdlibPathDetection(t *testing.T) {
	for _, p := range []string{"os/signal", "net/http", "syscall", "go/ast"} {
		if !isStdlibPath(p) {
			t.Errorf("%s 应判为标准库路径", p)
		}
	}
	for _, p := range []string{"github.com/ruoxizhnya/tellurion-kit/logging", "example.com/vendor/market/feed"} {
		if isStdlibPath(p) {
			t.Errorf("%s 不该判为标准库路径", p)
		}
	}
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func copyFixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	err := filepath.Walk(fixtureRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(fixtureRoot, path)
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
	if err != nil {
		t.Fatal(err)
	}
	return dst
}
