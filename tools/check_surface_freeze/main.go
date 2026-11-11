// Command check_surface_freeze 冻结本仓导出面的「形状指纹」（KIT-SURFACE-FREEZE）。
//
// 与 check_scope 的分工：那条核**包边界**（业务词 / 越界 import / 未登记包）；这条核**导出面形状**
// 有没有动 —— 改名、改签名、改字段、改底层类型、改常量字面量都要过基线 update，一次评审看得见的 diff。
// 本仓在依赖图最底层，一次 breaking 同时打到 contracts / core / agent，所以导出面必须有机器判据。
//
// 指纹格式（一行一个符号，确定性排序，去重）：
//
//	dir/Name(params) results              函数（含导出类型的导出方法）
//	dir/Type struct{字段 类型; ...}        导出结构体
//	dir/Type interface{方法签名; ...}      导出接口
//	dir/Type <底层类型>                    其它具名类型（含泛型类型参数）
//	dir/Name <类型>                       有显式类型的 const/var
//	dir/Name = <字面量>                   无显式类型、能取到字面量/表达式值的 const/var
//
// 三个红向，少一个就是假护栏：
//
//	基线缺失 —— docs/api-surface.txt 不存在、读不了，或一行指纹都没有
//	空扫红   —— 一个库包都没扫到（目录删空 / 改名），不许静默绿
//	形状漂移 —— 任何一行增删改都红，先列 `-`（基线有、代码没有）再列 `+`（代码有、基线没有）
//
// 不参与扫描：package main（tools/ 的护栏与命令不是导出面）、*_test.go、testdata/、vendor/、
// `.` / `_` 开头目录。用 go/parser + go/ast 渲染，不依赖 go/types —— 单仓 checkout 就能跑，
// 不需要解析依赖图。实现思路对照 tellurion-contracts 的 check_surface_freeze（同形状、独立一份）。
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ignoredDirs = map[string]bool{
	"testdata": true, // 夹具不是导出面
	"vendor":   true, // vendored 依赖的符号不是本仓的承诺
}

type options struct {
	root     string
	baseline string
	update   bool
	out      io.Writer
}

type packageDir struct {
	dir   string // slash 相对目录，根包记为 "."
	files []string
}

func main() {
	var opts options
	flag.StringVar(&opts.root, "root", ".", "仓库根目录（包目录相对于它计算）")
	flag.StringVar(&opts.baseline, "baseline", "docs/api-surface.txt", "冻结基线路径（相对仓库根）")
	flag.BoolVar(&opts.update, "update", false, "重新生成基线（默认模式只比对，不改文件）")
	flag.Parse()

	opts.out = os.Stdout
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "check_surface_freeze:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	out := opts.out
	if out == nil {
		out = os.Stdout
	}

	dirs, err := goDirs(opts.root)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return fmt.Errorf("一个 Go 包都没扫到（%s 下没有非测试 .go 文件）—— 护栏拒绝空绿", opts.root)
	}

	fset := token.NewFileSet()
	var lines []string
	scanned := 0
	for _, pkgDir := range dirs {
		pkgLines, isMain, err := scanPackage(opts.root, pkgDir, fset)
		if err != nil {
			return err
		}
		if isMain {
			continue // 命令/工具目录不是库的导出面
		}
		scanned++
		lines = append(lines, pkgLines...)
	}
	if scanned == 0 {
		return fmt.Errorf("扫描到 %d 个目录，但没有任何库包（全是 package main）—— 护栏拒绝空绿", len(dirs))
	}
	lines = dedupe(lines)

	baselinePath := filepath.Join(opts.root, filepath.FromSlash(opts.baseline))
	if opts.update {
		if err := os.WriteFile(baselinePath, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
			return fmt.Errorf("写基线 %s: %w", opts.baseline, err)
		}
		fmt.Fprintf(out, "导出面冻结基线已更新：%s（%d 个库包 / %d 行指纹）\n", opts.baseline, scanned, len(lines))
		return nil
	}

	want, err := readBaseline(baselinePath)
	if err != nil {
		return err
	}
	if len(want) == 0 {
		return fmt.Errorf("基线 %s 一行指纹都没有 —— 空基线等于没冻结；先跑 -update 生成", opts.baseline)
	}

	removed, added := diff(want, lines)
	fmt.Fprintf(out, "导出面冻结核对：%d 个库包 ｜ 代码 %d 行指纹 ｜ 基线 %d 行\n", scanned, len(lines), len(want))
	if len(removed) == 0 && len(added) == 0 {
		fmt.Fprintf(out, "✅ 导出面形状与 %s 逐行一致\n", opts.baseline)
		return nil
	}
	if len(removed) > 0 {
		fmt.Fprintf(out, "❌ 基线里有、代码里没有（%d）—— 删符号/改签名/改字段都算，确认后跑 -update：\n", len(removed))
		for _, line := range removed {
			fmt.Fprintln(out, "   -", line)
		}
	}
	if len(added) > 0 {
		fmt.Fprintf(out, "❌ 代码里有、基线里没有（%d）—— 新增符号要先登记 docs/SCOPE.md §3 与 CHANGELOG，再跑 -update：\n", len(added))
		for _, line := range added {
			fmt.Fprintln(out, "   +", line)
		}
	}
	return fmt.Errorf("导出面形状漂移：%d 行消失 + %d 行新增", len(removed), len(added))
}

// goDirs 列出 root 下含非测试 .go 文件的目录（相对 root，slash 形式；根包记为 "."）。
func goDirs(root string) ([]packageDir, error) {
	var dirs []packageDir
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if path != root && (ignoredDirs[name] || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")) {
				return filepath.SkipDir
			}
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		relDir, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		relFile, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		key := filepath.ToSlash(relDir)
		if n := len(dirs); n == 0 || dirs[n-1].dir != key {
			// WalkDir 按字典序遍历，同一目录的文件天然连续。
			dirs = append(dirs, packageDir{dir: key})
		}
		dirs[len(dirs)-1].files = append(dirs[len(dirs)-1].files, filepath.ToSlash(relFile))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

// scanPackage 解析一个目录的非测试文件，返回指纹行与「是否 package main」。
func scanPackage(root string, pkgDir packageDir, fset *token.FileSet) ([]string, bool, error) {
	var lines []string
	isMain := false
	pkgNames := map[string]bool{}
	for _, path := range pkgDir.files {
		// 语法错误必须当场红：解析失败的文件里所有符号会「凭空消失」，护栏静默空过。
		file, err := parser.ParseFile(fset, filepath.Join(root, filepath.FromSlash(path)), nil, parser.SkipObjectResolution)
		if err != nil {
			return nil, false, fmt.Errorf("解析 %s: %w", path, err)
		}
		pkgNames[file.Name.Name] = true
		if file.Name.Name == "main" {
			isMain = true
			continue
		}
		lines = append(lines, fileLines(fset, pkgDir.dir, file)...)
	}
	if len(pkgNames) > 1 && !isMain {
		names := make([]string, 0, len(pkgNames))
		for name := range pkgNames {
			names = append(names, name)
		}
		sort.Strings(names)
		return nil, false, fmt.Errorf("%s 里有多个 package（%s）：导出面按包冻结，一个目录只能有一个包", pkgDir.dir, strings.Join(names, ", "))
	}
	return lines, isMain, nil
}

// fileLines 收集一个文件里的导出符号指纹。
func fileLines(fset *token.FileSet, dir string, file *ast.File) []string {
	var lines []string
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name == nil || !d.Name.IsExported() {
				continue
			}
			if d.Recv == nil || len(d.Recv.List) == 0 {
				lines = append(lines, dir+"/"+d.Name.Name+renderSignature(fset, d.Type))
				continue
			}
			recv := receiverName(d.Recv.List[0].Type)
			if !ast.IsExported(recv) {
				continue
			}
			// 指针接收者与值接收者影响接口满足关系，指纹里必须区分。
			recvText := recv
			if _, ok := d.Recv.List[0].Type.(*ast.StarExpr); ok {
				recvText = "(*" + recv + ")"
			}
			lines = append(lines, dir+"/"+recvText+"."+d.Name.Name+renderSignature(fset, d.Type))
		case *ast.GenDecl:
			lines = append(lines, genDeclLines(fset, dir, d)...)
		}
	}
	return lines
}

func genDeclLines(fset *token.FileSet, dir string, decl *ast.GenDecl) []string {
	var lines []string
	switch decl.Tok {
	case token.TYPE:
		for _, spec := range decl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || !ts.Name.IsExported() {
				continue
			}
			name := dir + "/" + ts.Name.Name + renderTypeParams(fset, ts.TypeParams)
			if ts.Assign.IsValid() {
				lines = append(lines, name+" = "+renderExpr(fset, ts.Type))
				continue
			}
			lines = append(lines, name+" "+renderExpr(fset, ts.Type))
		}
	case token.CONST, token.VAR:
		for _, spec := range decl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if !name.IsExported() {
					continue
				}
				key := dir + "/" + name.Name
				switch {
				case vs.Type != nil:
					lines = append(lines, key+" "+renderExpr(fset, vs.Type))
				case i < len(vs.Values):
					if typ := compositeType(fset, vs.Values[i]); typ != "" {
						// 复合字面量只冻结它的类型（数据内容是种子数据，不是签名）；
						// 否则一行指纹会把整张表抄进来，diff 不可读。
						lines = append(lines, key+" "+typ)
						continue
					}
					lines = append(lines, key+" = "+renderExpr(fset, vs.Values[i]))
				case len(vs.Values) > 0:
					lines = append(lines, key+" = "+renderExpr(fset, vs.Values[0]))
				default:
					// 整组靠上一行表达式/iota 推导：取不到值就只留名字。
					lines = append(lines, key)
				}
			}
		}
	}
	return lines
}

// compositeType 取复合字面量的类型（`[]T{...}` → `[]T`，`&T{...}` → `T`）；
// 不是复合字面量或类型省略（嵌套字面量）时返回空串。
func compositeType(fset *token.FileSet, expr ast.Expr) string {
	if unary, ok := expr.(*ast.UnaryExpr); ok && unary.Op == token.AND {
		return compositeType(fset, unary.X)
	}
	lit, ok := expr.(*ast.CompositeLit)
	if !ok || lit.Type == nil {
		return ""
	}
	return renderExpr(fset, lit.Type)
}

// renderSignature 渲染 `(params) results`；结果为零个时不写，一个直接写类型，
// 多个写成括号列表（与 go doc 的形状一致）。
func renderSignature(fset *token.FileSet, ft *ast.FuncType) string {
	if ft == nil {
		return ""
	}
	sig := "(" + strings.Join(fieldTypes(fset, ft.Params), ", ") + ")"
	results := fieldTypes(fset, ft.Results)
	switch len(results) {
	case 0:
		return sig
	case 1:
		return sig + " " + results[0]
	default:
		return sig + " (" + strings.Join(results, ", ") + ")"
	}
}

// fieldTypes 按类型渲染参数/结果列表（丢弃参数名）：参数名不属于调用形状，
// 但类型、指针、切片、泛型、可变参数全部保留。
func fieldTypes(fset *token.FileSet, fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}
	var out []string
	for _, f := range fields.List {
		text := renderExpr(fset, f.Type)
		if len(f.Names) == 0 {
			out = append(out, text)
			continue
		}
		for range f.Names {
			out = append(out, text)
		}
	}
	return out
}

// renderTypeParams 渲染泛型参数表（`[K comparable, V any]`）；非泛型返回空串。
func renderTypeParams(fset *token.FileSet, params *ast.FieldList) string {
	if params == nil || len(params.List) == 0 {
		return ""
	}
	var parts []string
	for _, f := range params.List {
		text := renderExpr(fset, f.Type)
		if len(f.Names) == 0 {
			parts = append(parts, text)
			continue
		}
		names := make([]string, 0, len(f.Names))
		for _, n := range f.Names {
			names = append(names, n.Name)
		}
		parts = append(parts, strings.Join(names, ", ")+" "+text)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// renderExpr 按语法树渲染类型表达式；结构体与接口用 `;` 分隔字段/方法，保证一行可读。
func renderExpr(fset *token.FileSet, e ast.Expr) string {
	switch t := e.(type) {
	case nil:
		return ""
	case *ast.Ident:
		return t.Name
	case *ast.BasicLit:
		return t.Value
	case *ast.SelectorExpr:
		return renderExpr(fset, t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + renderExpr(fset, t.X)
	case *ast.ArrayType:
		length := ""
		if t.Len != nil {
			length = renderExpr(fset, t.Len)
		}
		return "[" + length + "]" + renderExpr(fset, t.Elt)
	case *ast.Ellipsis:
		return "..." + renderExpr(fset, t.Elt)
	case *ast.MapType:
		return "map[" + renderExpr(fset, t.Key) + "]" + renderExpr(fset, t.Value)
	case *ast.ChanType:
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + renderExpr(fset, t.Value)
		case ast.RECV:
			return "<-chan " + renderExpr(fset, t.Value)
		default:
			return "chan " + renderExpr(fset, t.Value)
		}
	case *ast.FuncType:
		return "func" + renderSignature(fset, t)
	case *ast.StructType:
		return renderStruct(fset, t)
	case *ast.InterfaceType:
		return renderInterface(fset, t)
	case *ast.IndexExpr:
		return renderExpr(fset, t.X) + "[" + renderExpr(fset, t.Index) + "]"
	case *ast.IndexListExpr:
		parts := make([]string, 0, len(t.Indices))
		for _, idx := range t.Indices {
			parts = append(parts, renderExpr(fset, idx))
		}
		return renderExpr(fset, t.X) + "[" + strings.Join(parts, ", ") + "]"
	case *ast.ParenExpr:
		return "(" + renderExpr(fset, t.X) + ")"
	case *ast.UnaryExpr:
		return t.Op.String() + renderExpr(fset, t.X)
	case *ast.BinaryExpr:
		return renderExpr(fset, t.X) + " " + t.Op.String() + " " + renderExpr(fset, t.Y)
	}
	// 兜底：剩余节点交给 go/printer，再把换行/制表压成单空格。
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, e); err != nil {
		return "?"
	}
	return strings.Join(strings.Fields(buf.String()), " ")
}

func renderStruct(fset *token.FileSet, st *ast.StructType) string {
	if st.Fields == nil || len(st.Fields.List) == 0 {
		return "struct{}"
	}
	parts := make([]string, 0, len(st.Fields.List))
	for _, f := range st.Fields.List {
		text := renderExpr(fset, f.Type)
		if len(f.Names) == 0 {
			parts = append(parts, text) // 内嵌字段：类型即字段
			continue
		}
		names := make([]string, 0, len(f.Names))
		for _, n := range f.Names {
			names = append(names, n.Name)
		}
		parts = append(parts, strings.Join(names, ", ")+" "+text)
	}
	return "struct{" + strings.Join(parts, "; ") + "}"
}

func renderInterface(fset *token.FileSet, it *ast.InterfaceType) string {
	if it.Methods == nil || len(it.Methods.List) == 0 {
		return "interface{}"
	}
	parts := make([]string, 0, len(it.Methods.List))
	for _, m := range it.Methods.List {
		if len(m.Names) == 0 {
			parts = append(parts, renderExpr(fset, m.Type)) // 内嵌接口 / 类型集
			continue
		}
		sig := renderExpr(fset, m.Type)
		if ft, ok := m.Type.(*ast.FuncType); ok {
			sig = renderSignature(fset, ft)
		}
		for _, n := range m.Names {
			// 方法写成 `Name(args) results`：方法名后面直接跟签名，不留 `func` 前缀。
			parts = append(parts, n.Name+sig)
		}
	}
	return "interface{" + strings.Join(parts, "; ") + "}"
}

// receiverName 取接收者的具名类型（剥掉指针与泛型下标）。
func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return receiverName(t.X)
	case *ast.IndexListExpr:
		return receiverName(t.X)
	}
	return ""
}

// readBaseline 读基线：空行不算指纹（文件尾换行）；文件不存在即红。
func readBaseline(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读冻结基线 %s: %w（缺失即红；先跑 -update 生成）", path, err)
	}
	var lines []string
	for _, line := range strings.Split(strings.ReplaceAll(string(raw), "\r\n", "\n"), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// diff 返回只在基线里的行与只在代码里的行（各自升序）。
func diff(baseline, current []string) (removed, added []string) {
	want := make(map[string]bool, len(baseline))
	for _, line := range baseline {
		want[line] = true
	}
	have := make(map[string]bool, len(current))
	for _, line := range current {
		have[line] = true
	}
	for _, line := range baseline {
		if !have[line] {
			removed = append(removed, line)
		}
	}
	for _, line := range current {
		if !want[line] {
			added = append(added, line)
		}
	}
	sort.Strings(removed)
	sort.Strings(added)
	return removed, added
}

func dedupe(lines []string) []string {
	seen := make(map[string]bool, len(lines))
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if seen[line] {
			continue
		}
		seen[line] = true
		out = append(out, line)
	}
	sort.Strings(out)
	return out
}
