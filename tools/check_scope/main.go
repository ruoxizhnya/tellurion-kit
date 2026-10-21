// Command check_scope 核对 tellurion-kit 的包边界（KT-01 ③-③ 的 CI 护栏）。
//
// 四个红向，少一个就是假护栏：
//
//	业务词命中   —— 非测试文件里**声明出来的**标识符（函数 / 方法 / 类型 / 接口方法 / 结构体字段 / const / var）
//	              或**非标准库 import 路径的段**命中 docs/SCOPE.md §2 词表
//	越界依赖     —— import 了 core（module 首段 `tellurion`）或任何兄弟仓（`github.com/ruoxizhnya/tellurion-<非 kit>`）
//	清单没覆盖   —— 本仓出现了 SCOPE.md 十三包之外、也不在任何登记包子树里的 Go 包目录
//	结构红       —— SCOPE.md 找不到 scope:begin/end 区段、表头分隔行被删、或十三包少任何一行
//	              （「把清单删了 / 缩了来让检查通过」这条路堵死，所以这三条是 fatal 而不是可忽略的红）
//
// 不参与核对（这三条限制在 docs/SCOPE.md §2 如实登记，不假装覆盖）：
//
//	注释与字符串字面量 —— 注释要说清调用方是谁（「供回测循环使用」是正当的），扫它会把文档本身变成违规物；
//	*_test.go          —— 库对外的形状由非测试文件决定；
//	package main       —— 自查工具自己必须写出这些词才能拦住它们（词表就在本文件里）。
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	scopeBegin = "<!-- scope:begin -->"
	scopeEnd   = "<!-- scope:end -->"
	kitModule  = "github.com/ruoxizhnya/tellurion-kit"
	kitOwner   = "ruoxizhnya"
)

// businessWords 是 docs/SCOPE.md §2 那张表的机械形态，值 = 出处。
// 冻结来源：本仓 AGENTS.md §3 禁入清单（backtest/strategy/order/factor/market/portfolio）
// + QT-05-B ③ 点名的 core 专用类型（screen/sector/supply_chain/macro/mining）
// + KT-01 补充的同族数据形状（ohlcv/bar/signal）。改词表先改 SCOPE.md 与出处那张单。
var businessWords = map[string]string{
	"backtest":    "AGENTS §3",
	"strategy":    "AGENTS §3",
	"order":       "AGENTS §3",
	"factor":      "AGENTS §3",
	"market":      "AGENTS §3",
	"portfolio":   "AGENTS §3",
	"screen":      "QT-05-B ③",
	"sector":      "QT-05-B ③",
	"supplychain": "QT-05-B ③",
	"macro":       "QT-05-B ③",
	"mining":      "QT-05-B ③",
	"ohlcv":       "KT-01 补充",
	"bar":         "KT-01 补充",
	"signal":      "KT-01 补充",
}

// scopePackages 是 SCOPE.md §3 必须逐包列全的十三行（QT-05 ③-②：12 个横切包 + 新增 db）。
// 名单冻结在工具里是刻意的：清单少一行就 fatal，否则「删一行」就等于「少登记一个包」。
var scopePackages = []string{
	"logging", "errors", "id", "decimal", "statistics",
	"httpclient", "httpserver", "auth", "authz", "observability",
	"clock", "sandbox", "db",
}

var (
	tableSep = regexp.MustCompile(`^[\s:|-]+$`)
	// 切词的第一层：下划线、连字符、其它非字母数字字符。驼峰与数字边界在 splitCamel 里处理。
	tokenSplit = regexp.MustCompile(`[^A-Za-z0-9]+`)
)

type options struct {
	root  string
	scope string
	out   io.Writer
}

func main() {
	var opts options
	flag.StringVar(&opts.root, "root", ".", "仓库根目录（包路径相对于它计算）")
	flag.StringVar(&opts.scope, "scope", "docs/SCOPE.md", "逐包准入 / 禁入清单（KT-01 的载体）")
	flag.Parse()

	opts.out = os.Stdout
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "check_scope:", err)
		os.Exit(1)
	}
}

// goFile 是解析过一次的文件：包名、AST 与位置信息都从它取，不重复解析。
type goFile struct {
	name string // 相对仓库根的 slash 路径
	fset *token.FileSet
	ast  *ast.File
}

type violation struct {
	pos string
	why string
}

func run(opts options) error {
	out := opts.out
	if out == nil {
		out = os.Stdout
	}

	registered, err := parseScope(filepath.Join(opts.root, filepath.FromSlash(opts.scope)))
	if err != nil {
		return err
	}

	dirs, err := goDirs(opts.root)
	if err != nil {
		return err
	}

	var red []violation
	scanned, skipped := 0, 0
	for _, dir := range dirs {
		files, pkgName, err := parseDir(filepath.Join(opts.root, filepath.FromSlash(dir)), dir)
		if err != nil {
			return err
		}
		if pkgName == "main" {
			skipped++
			continue
		}
		scanned++

		if dir != "." {
			if top := strings.Split(dir, "/")[0]; !registered[top] {
				red = append(red, violation{
					pos: "./" + dir,
					why: fmt.Sprintf("这个包没登记进 %s 的逐包清单（新增包必须先进 §3，见 KT-01 ③-①）", opts.scope),
				})
			}
		}

		for _, f := range files {
			red = append(red, checkImports(f)...)
			red = append(red, checkDecls(f)...)
		}
	}

	sort.Slice(red, func(i, j int) bool {
		if red[i].pos != red[j].pos {
			return red[i].pos < red[j].pos
		}
		return red[i].why < red[j].why
	})

	fmt.Fprintf(out, "包边界核对：%d 个库包已扫 / %d 个 main 包跳过｜SCOPE.md 登记 %d 包\n",
		scanned, skipped, len(registered))
	if len(red) > 0 {
		for _, v := range red {
			fmt.Fprintf(out, "❌ %s：%s\n", v.pos, v.why)
		}
		return fmt.Errorf("包边界违规（%d 项）", len(red))
	}
	fmt.Fprintln(out, "✅ 业务词零命中、无越界 import、十三包清单齐且覆盖现有包")
	return nil
}

// parseScope 读登记区，返回「已登记包名」集合。
// 区段缺失 / 分隔行缺失 / 十三包少任何一行都返回错误，而不是记一条可讨论的红。
func parseScope(path string) (map[string]bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inside, sawSep bool
	var cells []string
	for _, line := range strings.Split(string(raw), "\n") {
		t := strings.TrimSpace(line)
		switch t {
		case scopeBegin:
			inside = true
			continue
		case scopeEnd:
			inside = false
			continue
		}
		if !inside || !strings.HasPrefix(t, "|") {
			continue
		}
		row := splitCells(t)
		if len(row) == 0 {
			continue
		}
		if allSep(row) {
			sawSep = true
			continue
		}
		cells = append(cells, row[0])
	}
	if !sawSep {
		return nil, fmt.Errorf("%s 的登记区没有 markdown 表格分隔行（删掉表头会让第一行数据被当表头跳过，护栏随之变松）", filepath.Base(path))
	}
	if len(cells) > 0 {
		cells = cells[1:] // 表头
	}
	registered := map[string]bool{}
	for _, c := range cells {
		if name := strings.Trim(c, "` "); name != "" {
			registered[name] = true
		}
	}
	var missing []string
	for _, p := range scopePackages {
		if !registered[p] {
			missing = append(missing, p)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%s 的逐包清单少了 %s（十三包 = QT-05 ③-② 的 12 个横切包 + 新增 db）",
			filepath.Base(path), strings.Join(missing, " / "))
	}
	return registered, nil
}

func splitCells(line string) []string {
	parts := strings.Split(strings.Trim(line, "|"), "|")
	cells := make([]string, 0, len(parts))
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	return cells
}

func allSep(cells []string) bool {
	for _, c := range cells {
		if c == "" || !tableSep.MatchString(c) {
			return false
		}
	}
	return len(cells) > 0
}

// goDirs 返回 root 下含非测试 .go 文件的目录（slash 相对路径；根包记为 "."）。
func goDirs(root string) ([]string, error) {
	set := map[string]bool{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if p != root && (name == "testdata" || name == "vendor" ||
				strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		rel, err := filepath.Rel(root, filepath.Dir(p))
		if err != nil {
			return err
		}
		set[filepath.ToSlash(rel)] = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	dirs := make([]string, 0, len(set))
	for dir := range set {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	return dirs, nil
}

// parseDir 解析目录里的非测试 Go 文件，返回文件列表与该目录的包名。
// relDir 只用于拼出「相对仓库根」的位置串（报错要能指到人能改的地方）。
// 一个目录两种包名在这里直接报错：那种仓 `go build` 本来就过不去，但静态扫不允许静默少扫一半。
func parseDir(absDir, relDir string) ([]goFile, string, error) {
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return nil, "", err
	}
	var files []goFile
	pkgName := ""
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filepath.Join(absDir, name), nil, parser.AllErrors)
		if err != nil {
			return nil, "", fmt.Errorf("解析 %s/%s: %w", relDir, name, err)
		}
		if pkgName == "" {
			pkgName = f.Name.Name
		} else if pkgName != f.Name.Name {
			return nil, "", fmt.Errorf("%s 目录里有两个包名（%s 与 %s）", relDir, pkgName, f.Name.Name)
		}
		rel := name
		if relDir != "." {
			rel = relDir + "/" + name
		}
		files = append(files, goFile{name: filepath.ToSlash(rel), fset: fset, ast: f})
	}
	return files, pkgName, nil
}

// checkImports 核对 import 路径：越界依赖与路径段业务词。
func checkImports(f goFile) []violation {
	var red []violation
	for _, spec := range f.ast.Imports {
		path := strings.Trim(spec.Path.Value, `"`)
		pos := f.pos(spec.Pos())
		if why := importProblem(path); why != "" {
			red = append(red, violation{pos: fmt.Sprintf("%s:%d import %q", f.name, pos, path), why: why})
			continue
		}
		if isStdlibPath(path) {
			continue // os/signal 这类标准库段不是交易信号
		}
		for _, seg := range strings.Split(path, "/") {
			if w := wordHit(seg); w != "" {
				red = append(red, violation{
					pos: fmt.Sprintf("%s:%d import %q", f.name, pos, path),
					why: fmt.Sprintf("路径段命中业务词 %q（出处：%s）", w, businessWords[w]),
				})
				break
			}
		}
	}
	return red
}

// checkDecls 核对声明出来的标识符名。
// 形参与局部变量不扫：`order` 在普通英语里是「次序」，扫它会误红；这条限制登记在 SCOPE.md §2。
func checkDecls(f goFile) []violation {
	var red []violation
	add := func(kind, name string, pos token.Pos) {
		w := wordHit(name)
		if w == "" {
			return
		}
		red = append(red, violation{
			pos: fmt.Sprintf("%s:%d", f.name, f.pos(pos)),
			why: fmt.Sprintf("%s %q 命中业务词 %q（出处：%s）—— 业务概念不许进 kit，见 docs/SCOPE.md §2/§3",
				kind, name, w, businessWords[w]),
		})
	}
	ast.Inspect(f.ast, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			add("函数", x.Name.Name, x.Name.Pos())
		case *ast.TypeSpec:
			add("类型", x.Name.Name, x.Name.Pos())
		case *ast.ValueSpec:
			for _, id := range x.Names {
				add("声明", id.Name, id.Pos())
			}
		case *ast.StructType:
			for _, fld := range x.Fields.List {
				for _, id := range fld.Names {
					add("字段", id.Name, id.Pos())
				}
			}
		case *ast.InterfaceType:
			for _, m := range x.Methods.List {
				for _, id := range m.Names {
					add("接口方法", id.Name, id.Pos())
				}
			}
		}
		return true
	})
	return red
}

func (f goFile) pos(p token.Pos) int {
	return f.fset.Position(p).Line
}

// importProblem 返回越界依赖的说明；不是越界则空串。
func importProblem(path string) string {
	if path == "" {
		return ""
	}
	segs := strings.Split(path, "/")
	if segs[0] == "tellurion" {
		return "越界依赖：core 的 module 首段是 `tellurion`，kit 是依赖图底座，不得 import 任何应用仓实现包"
	}
	if len(segs) >= 3 && segs[0] == "github.com" && segs[1] == kitOwner && strings.HasPrefix(segs[2], "tellurion-") {
		if "github.com/"+kitOwner+"/"+segs[2] == kitModule {
			return "" // 本仓自己的包，允许
		}
		return fmt.Sprintf("越界依赖：%s 是同项目兄弟仓，不变量是 kit ← contracts ← core，反向一条都不许有", segs[2])
	}
	return ""
}

// isStdlibPath：标准库路径的首段永远不含点（`os/signal`、`net/http`）；含点的是第三方或本仓 module。
func isStdlibPath(path string) bool {
	return !strings.Contains(strings.Split(path, "/")[0], ".")
}

// wordHit 切词后与词表比对；命中返回键，否则空串。
// 两条比法：单段整词（OrderID -> order），以及**相邻两段拼接**（SupplyChain -> supply+chain = supplychain）。
// 拼接这条是为了 `supply_chain` 这类复合词：只把 "supply" 或 "chain" 单独列进词表会误红
// （kit 里 `Chain` 是链表的正当名字），整词又抓不到 SupplyChain —— 相邻拼接同时解决这两头。
func wordHit(name string) string {
	toks := tokens(name)
	for i, tok := range toks {
		if _, ok := businessWords[tok]; ok {
			return tok
		}
		if i+1 < len(toks) {
			joined := tok + toks[i+1]
			if _, ok := businessWords[joined]; ok {
				return joined
			}
		}
	}
	return ""
}

// tokens 按驼峰、下划线、数字边界切词并小写。
// OrderID -> [order id]；ObserveBacktest -> [observe backtest]；MaskAPIKey -> [mask api key]；
// Factorial -> [factorial]（不会误命中 factor，这正是整词比对的意义）。
func tokens(s string) []string {
	var out []string
	for _, p := range nonEmpty(tokenSplit.Split(s, -1)) {
		for _, w := range splitCamel(p) {
			out = append(out, strings.ToLower(w))
		}
	}
	return out
}

func nonEmpty(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func splitCamel(s string) []string {
	runes := []rune(s)
	var words []string
	start := 0
	for i := 1; i < len(runes); i++ {
		prev, cur := runes[i-1], runes[i]
		boundary := (!isUpper(prev) && isUpper(cur)) ||
			(isUpper(prev) && isUpper(cur) && i+1 < len(runes) && isLower(runes[i+1])) ||
			(isDigit(prev) != isDigit(cur))
		if boundary {
			words = append(words, s[start:i])
			start = i
		}
	}
	return append(words, s[start:])
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }
func isDigit(r rune) bool { return r >= '0' && r <= '9' }
