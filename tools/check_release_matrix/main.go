// Command check_release_matrix 核对版本纪律三件套（KT-02 ③/⑤）：
// README 的版本规则表与兼容矩阵、CHANGELOG 的版本节、git tag。
//
// 六个红向，少一个就是假护栏：
//
//	结构红     —— README 缺 rules / matrix 标记，规则表缺 MAJOR/MINOR/PATCH 任一行，或矩阵没有任何数据行
//	发版没更新矩阵 —— CHANGELOG 里出现了矩阵没声明的版本序列
//	矩阵行没账 —— 矩阵声明了一个序列，CHANGELOG 里却没有任何该序列的版本节
//	tag 没账   —— 每个 v* tag 必须同时有 CHANGELOG 节与矩阵行（tag 不 push 等于没发版，反之亦然）
//	前提缺失   —— AGENTS.md 里没有 GOPRIVATE 那一行（本仓转私有后 CI/构建必然拉不到）
//	空账       —— CHANGELOG 一条条目都没有
//
// 口径：兼容矩阵按**版本序列** `vX.Y` 建行（同一序列的补丁版共用一行，源头见 QT-19 ③-③），
// CHANGELOG 按具体 tag `vX.Y.Z` 建节。两者靠序列对齐，不靠字符串相等。
//
// 与 contracts 的同名工具（CT-02）为什么各写一份而不共用：本仓在依赖图的**最底层**，
// 发版时要知道「contracts / core / agent 三方各 pin 到哪」，那份矩阵的形状与使用方都与契约仓不同；
// 而 `kit` 不能 import 任何应用仓、别仓也无法 require 本仓之外的公共工具（R2：单仓可自证）。
//
// 不参与核对：别的仓的 require 是否落在矩阵范围内 —— 那是跨仓事实，本仓 checkout 不到，
// 归 tellurion:docs/TASKS.md QT-09 的 CI 步骤。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	rulesBegin    = "<!-- rules:begin -->"
	rulesEnd      = "<!-- rules:end -->"
	matrixBegin   = "<!-- matrix:begin -->"
	matrixEnd     = "<!-- matrix:end -->"
	goprivateMark = "GOPRIVATE=github.com/ruoxizhnya/"
	unreleased    = "unreleased"
)

var (
	// 矩阵首格：vX.Y 或 vX.Y.Z 都收，统一折算成序列 vX.Y。
	cellVersion = regexp.MustCompile(`^v\d+\.\d+(?:\.\d+)?$`)
	// 首格像版本号但形状不合法（打错的行不能静默当信息行放过）。
	looksVersionish = regexp.MustCompile(`^v\d`)
	// CHANGELOG 版本节标题：只认完整 vX.Y.Z（tag 的形状）。
	sectionVersion = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	// 标题尾巴：` - YYYY-MM-DD`。
	releaseDate = regexp.MustCompile(`^- \d{4}-\d{2}-\d{2}$`)
	// `## [vX.Y.Z] - 2026-11-01` 与 `## vX.Y.Z - 2026-11-01` 两种写法都解析（内容被看见比格式洁癖重要）。
	heading = regexp.MustCompile(`^##\s*(?:\[([^\]]+)\]|(\S+))\s*(.*)$`)
	// markdown 表格分隔行：|---|---|
	tableSep = regexp.MustCompile(`^[\s:|-]+$`)
)

type options struct {
	root      string
	readme    string
	changelog string
	agents    string
	listTags  func() ([]string, error)
	out       io.Writer
}

func main() {
	var opts options
	flag.StringVar(&opts.root, "root", ".", "仓库根目录（三份文件相对于它计算）")
	flag.StringVar(&opts.readme, "readme", "README.md", "版本规则表与兼容矩阵所在文件")
	flag.StringVar(&opts.changelog, "changelog", "CHANGELOG.md", "版本账本")
	flag.StringVar(&opts.agents, "agents", "AGENTS.md", "私有化前提所在文件")
	flag.Parse()

	opts.out = os.Stdout
	opts.listTags = gitTags(opts.root)
	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, "check_release_matrix:", err)
		os.Exit(1)
	}
}

func run(opts options) error {
	out := opts.out
	if out == nil {
		out = os.Stdout
	}
	if opts.listTags == nil {
		opts.listTags = gitTags(opts.root)
	}

	var red []string

	rulesRows, err := tableRows(filepath.Join(opts.root, filepath.FromSlash(opts.readme)), rulesBegin, rulesEnd)
	if err != nil {
		return fmt.Errorf("README 版本规则表：%w", err)
	}
	if missing := missingRuleKeys(rulesRows); len(missing) > 0 {
		red = append(red, fmt.Sprintf("版本规则表缺行：%s（MAJOR/MINOR/PATCH 三行缺一不可）", strings.Join(missing, " ")))
	}

	matrixRows, err := tableRows(filepath.Join(opts.root, filepath.FromSlash(opts.readme)), matrixBegin, matrixEnd)
	if err != nil {
		return fmt.Errorf("README 兼容矩阵：%w", err)
	}
	series, malformed, dup, err := matrixSeries(matrixRows)
	if err != nil {
		return fmt.Errorf("README 兼容矩阵：%w", err)
	}
	for _, cell := range malformed {
		if cell == "" {
			red = append(red, "矩阵有首格为空的数据行（版本序列那一列不能空）")
			continue
		}
		red = append(red, fmt.Sprintf("矩阵首格 %q 像版本号但形状不合法（应为 vX.Y 或 vX.Y.Z）", cell))
	}
	for _, s := range dup {
		red = append(red, fmt.Sprintf("矩阵里序列 %s 有多行（一个序列只能一行，补丁版更新同一行）", s))
	}

	sections, items, badTitles, err := changelogSections(filepath.Join(opts.root, filepath.FromSlash(opts.changelog)))
	if err != nil {
		return fmt.Errorf("CHANGELOG：%w", err)
	}
	for _, title := range badTitles {
		red = append(red, fmt.Sprintf("CHANGELOG 标题 %q 不合规范（版本节要写成 `## [vX.Y.Z] - YYYY-MM-DD`）", title))
	}
	if items == 0 {
		red = append(red, "CHANGELOG 一条条目都没有（⑤ 要求 ≥ 1 条）")
	}

	// 发版没更新矩阵 / 矩阵行没账：两个方向都算红。
	for _, sec := range sections {
		if s := seriesOf(sec); !series[s] {
			red = append(red, fmt.Sprintf("CHANGELOG 有 %s 这一版，矩阵却没有 %s 那一行 —— 发版必须同步兼容矩阵", sec, s))
		}
	}
	for s := range series {
		if !sectionCovers(sections, s) {
			red = append(red, fmt.Sprintf("矩阵声明了序列 %s，CHANGELOG 里却没有该序列的任何版本节", s))
		}
	}

	tags, err := opts.listTags()
	if err != nil {
		return fmt.Errorf("读取 tag 列表失败：%w（CI 需要完整 tag 历史；本仓 checkout 步骤须 fetch-depth: 0）", err)
	}
	for _, tag := range tags {
		if !hasVersion(sections, tag) {
			red = append(red, fmt.Sprintf("tag %s 在 CHANGELOG 里没有对应版本节", tag))
		}
		if s := seriesOf(tag); s != "" && !series[s] {
			red = append(red, fmt.Sprintf("tag %s 的版本序列 %s 不在兼容矩阵里", tag, s))
		}
	}

	if err := checkPrerequisite(filepath.Join(opts.root, filepath.FromSlash(opts.agents))); err != nil {
		red = append(red, err.Error())
	}

	fmt.Fprintf(out, "版本纪律核对：规则表 %d 行 ｜ 矩阵序列 %d ｜ CHANGELOG 版本节 %d（条目 %d）｜ tag %d\n",
		len(rulesRows), len(series), len(sections), items, len(tags))
	if len(red) > 0 {
		sort.Strings(red)
		for _, r := range red {
			fmt.Fprintln(out, "❌", r)
		}
		return fmt.Errorf("版本纪律不一致（%d 项）", len(red))
	}
	fmt.Fprintln(out, "✅ 兼容矩阵 ↔ CHANGELOG ↔ tag 一致，规则表与 GOPRIVATE 前提在册")
	return nil
}

// sectionText 返回两个标记之间（不含标记）的内容行。
// 缺任一模板、顺序颠倒、或区段里没有任何内容都算错——「把表删了就绿」必须是红。
func sectionText(path, begin, end string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var text []string
	inside := false
	for _, line := range strings.Split(string(raw), "\n") {
		switch strings.TrimSpace(line) {
		case begin:
			inside = true
			continue
		case end:
			if !inside {
				return nil, fmt.Errorf("出现 %s 但没有配对的 %s", end, begin)
			}
			inside = false
			continue
		}
		if inside {
			text = append(text, line)
		}
	}
	if inside {
		return nil, fmt.Errorf("有 %s 但没有配对的 %s", begin, end)
	}
	if len(text) == 0 {
		return nil, fmt.Errorf("区段 %s…%s 之间没有任何内容", begin, end)
	}
	return text, nil
}

// tableRows 解析区段里的 markdown 表格，返回数据行（跳过表头与分隔行）。
// 没有分隔行就报错：删掉表头会让第一行数据被当成表头跳过，护栏随之悄悄变松。
func tableRows(path, begin, end string) ([][]string, error) {
	text, err := sectionText(path, begin, end)
	if err != nil {
		return nil, err
	}
	var rows [][]string
	sawSep := false
	for _, line := range text {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := splitCells(line)
		if len(cells) == 0 {
			continue
		}
		if allSep(cells) {
			sawSep = true
			continue
		}
		rows = append(rows, cells)
	}
	if !sawSep {
		return nil, fmt.Errorf("%s 区段里没有 markdown 表格（要表头行 + |---| 分隔行）", begin)
	}
	// 第一行数据是表头（markdown 语义），不参与内容核对。
	if len(rows) > 0 {
		rows = rows[1:]
	}
	return rows, nil
}

func splitCells(line string) []string {
	trimmed := strings.Trim(line, "|")
	parts := strings.Split(trimmed, "|")
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

func missingRuleKeys(rows [][]string) []string {
	want := []string{"MAJOR", "MINOR", "PATCH"}
	have := map[string]bool{}
	for _, r := range rows {
		have[firstCell(r)] = true
	}
	var missing []string
	for _, k := range want {
		if !have[k] {
			missing = append(missing, k)
		}
	}
	return missing
}

// firstCell 去掉反引号与空格，返回首格内容（无首格时返回空串）。
func firstCell(row []string) string {
	if len(row) == 0 {
		return ""
	}
	return strings.Trim(row[0], "` ")
}

// matrixSeries 返回矩阵声明的版本序列集合；另返回形状不合法的首格与重复序列。
func matrixSeries(rows [][]string) (map[string]bool, []string, []string, error) {
	if len(rows) == 0 {
		return nil, nil, nil, fmt.Errorf("矩阵没有任何数据行（至少要有一行说明当前基线）")
	}
	series := map[string]bool{}
	var malformed, dup []string
	seen := map[string]int{}
	for _, r := range rows {
		cell := firstCell(r)
		if cell == "" {
			malformed = append(malformed, cell)
			continue
		}
		if !cellVersion.MatchString(cell) {
			// 非版本行（如 `main`（未发版））是允许的现状说明行；
			// 但以 v+数字开头的多半是打错的版本号，必须点名而不是放过。
			if looksVersionish.MatchString(cell) {
				malformed = append(malformed, cell)
			}
			continue
		}
		s := seriesOf(cell)
		series[s] = true
		seen[s]++
		if seen[s] == 2 {
			dup = append(dup, s)
		}
	}
	return series, malformed, dup, nil
}

// seriesOf 把 vX.Y(.Z) 折算成序列 vX.Y；不是版本形状则返回空串。
func seriesOf(v string) string {
	v = strings.Trim(v, "` ")
	if !cellVersion.MatchString(v) {
		return ""
	}
	parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
	return "v" + parts[0] + "." + parts[1]
}

// changelogSections 返回版本节里的版本号（升序）、条目行数、以及标题不合规的原文。
func changelogSections(path string) ([]string, int, []string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, nil, err
	}
	var versions, badTitles []string
	items := 0
	seen := map[string]bool{}
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "- ") {
			items++
		}
		// 只看二级标题；`### Added` 这类子节标题留给内容，不参与版本节识别。
		if strings.HasPrefix(line, "###") {
			continue
		}
		m := heading.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		label := m[1]
		if label == "" {
			label = m[2]
		}
		tail := strings.TrimSpace(m[3])
		if strings.EqualFold(label, unreleased) {
			continue
		}
		if !sectionVersion.MatchString(label) {
			// 非版本的 `## xxx` 标题（如 `## 用法`）不属于版本节，跳过。
			if looksVersionish.MatchString(label) {
				badTitles = append(badTitles, strings.TrimSpace(line))
			}
			continue
		}
		if !releaseDate.MatchString(tail) {
			badTitles = append(badTitles, strings.TrimSpace(line))
		}
		if !seen[label] {
			seen[label] = true
			versions = append(versions, label)
		}
	}
	sort.Strings(versions)
	return versions, items, badTitles, nil
}

func hasVersion(versions []string, tag string) bool {
	for _, v := range versions {
		if v == strings.TrimSpace(tag) {
			return true
		}
	}
	return false
}

// sectionCovers 判断版本节集合里是否有落在该序列中的版本。
func sectionCovers(versions []string, s string) bool {
	for _, v := range versions {
		if seriesOf(v) == s {
			return true
		}
	}
	return false
}

// checkPrerequisite 钉住私有化前提那一行：转私有后没配 GOPRIVATE，CI 与本地 go get 必然红。
func checkPrerequisite(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.Contains(line, goprivateMark) {
			return nil
		}
	}
	return fmt.Errorf("%s 里没有 %s 那一行（私有化前提必须成文）", filepath.Base(path), goprivateMark+"*")
}

// gitTags 从仓库读取 v* tag 列表；读不到就返回错误（静默跳过会让 tag 这条腿变成假绿）。
func gitTags(root string) func() ([]string, error) {
	return func() ([]string, error) {
		cmd := exec.Command("git", "-C", root, "tag", "-l", "v*")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		var tags []string
		for _, line := range strings.Split(string(out), "\n") {
			if t := strings.TrimSpace(line); t != "" {
				tags = append(tags, t)
			}
		}
		sort.Strings(tags)
		return tags, nil
	}
}
