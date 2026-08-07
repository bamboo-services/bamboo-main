// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

// Package main 提供旧库（xf_index）数据迁移工具。
//
// 用法：
//
//	go run ./cmd/import-old-data
//
// 默认读取 data/backup.sql（pg_dump custom 格式），先经 pg_restore 转成
// plain SQL，再解析 COPY 数据，最后用与运行时一致的雪花算法（xSnowflake，
// datacenter=7, node=3，业务基因取自 pkg/constants）生成主键写入新库。
// 数据库连接参数可用环境变量 DATABASE_HOST/PORT/USER/PASS/NAME 覆盖，
// 缺省为 localhost:5432 bamboo_main/bamboo_main/bamboo_main。
//
// 注意：脚本会清空 bm_link_color / bm_link_group / bm_link_friend 三张表
// 后重新导入，请在开发/迁移环境使用。
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ----------------------------------------------------------------------------
// 雪花算法位布局常量（复制自 bamboo-base-go common v1.2.0 snowflake/snowflake.go）
// epoch = 2023-07-25 00:00:00 UTC；布局：1位符号 + 41位毫秒时间戳 + 6位基因
// + 3位数据中心 + 3位节点 + 10位序列号。
// datacenter=7 / node=3 与当前库种子数据（bm_system / bm_system_user）反推一致。
// ----------------------------------------------------------------------------
const (
	snowflakeEpoch  int64 = 1690214400000
	timestampShift        = 22
	geneShift             = 16
	datacenterShift       = 13
	nodeShift             = 10
	sequenceBits          = 10
	maxSequence           = 1023
	datacenterID          = 7
	nodeID                = 3
	// builtinGroupCount 内置分组数量（已失效分组不参与排序持久化），迁移时旧分组排序值直接沿用。
	builtinGroupCount = 0
)

// snowflakeID 按 xSnowflake 位布局构造指定时间的雪花 ID。
func snowflakeID(ts time.Time, gene xSnowflake.Gene, seq int64) xSnowflake.SnowflakeID {
	ms := ts.UnixMilli()
	v := ((ms - snowflakeEpoch) << timestampShift) |
		(int64(gene) << geneShift) |
		(datacenterID << datacenterShift) |
		(nodeID << nodeShift) |
		(seq & maxSequence)
	return xSnowflake.SnowflakeID(v)
}

// ----------------------------------------------------------------------------
// COPY 数据解析
// ----------------------------------------------------------------------------

var copyHeaderRe = regexp.MustCompile(`(?m)^COPY public\.(\w+) \(([^)]*)\) FROM stdin;\n`)

// copyBlock 一段 COPY 数据块。
type copyBlock struct {
	table string
	cols  []string
	rows  [][]string
}

// parseCopyBlocks 解析 pg_restore --data-only 输出的 plain SQL 中的 COPY 段。
func parseCopyBlocks(sql string) []copyBlock {
	var blocks []copyBlock
	rest := sql
	for {
		loc := copyHeaderRe.FindStringSubmatchIndex(rest)
		if loc == nil {
			break
		}
		table := rest[loc[2]:loc[3]]
		cols := strings.Split(rest[loc[4]:loc[5]], ", ")
		bodyStart := loc[1]
		// 数据行从头部结尾之后开始，到 ^\.$ 结束
		body := rest[bodyStart:]
		end := strings.Index(body, "\n\\.")
		if end < 0 {
			break
		}
		data := body[:end]
		rows := parseCopyRows(data)
		blocks = append(blocks, copyBlock{table: table, cols: cols, rows: rows})
		rest = rest[bodyStart+end+3:]
	}
	return blocks
}

// parseCopyRows 解析 COPY 数据行（tab 分隔，\N 表示 NULL，反斜杠转义）。
func parseCopyRows(data string) [][]string {
	var rows [][]string
	sc := bufio.NewScanner(strings.NewReader(data))
	sc.Buffer(make([]byte, 1024*1024), 64*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		for i := range fields {
			fields[i] = copyUnescape(fields[i])
		}
		rows = append(rows, fields)
	}
	return rows
}

// copyUnescape 还原 COPY 文本格式的转义序列。
func copyUnescape(s string) string {
	if s == `\N` {
		return "\x00NULL\x00" // 哨兵值，调用处转 NULL
	}
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case 'v':
			b.WriteByte('\v')
		case '\\':
			b.WriteByte('\\')
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// ----------------------------------------------------------------------------
// 旧数据行结构
// ----------------------------------------------------------------------------

type oldColor struct {
	id          int64
	name        string
	displayName string
	color       string
	hasSelect   bool
	createdAt   time.Time
	updatedAt   *time.Time
}

type oldLocation struct {
	id          int64
	sort        int
	name        string
	displayName string
	description *string
	reveal      bool
	createdAt   time.Time
	updatedAt   *time.Time
}

type oldLink struct {
	id              int64
	webmasterEmail  *string
	siteName        string
	siteURL         string
	siteLogo        *string
	cdnLogoURL      *string
	siteDesc        *string
	siteRSSURL      *string
	location        *int64
	color           *int64
	webmasterRemark *string
	remark          *string
	status          int
	createdAt       time.Time
	updatedAt       *time.Time
}

const timeLayout = "2006-01-02 15:04:05.999999"

func parseTime(s string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(timeLayout, s, loc)
}

func nullableStr(f string) *string {
	if f == "\x00NULL\x00" || f == "" {
		return nil
	}
	return &f
}

// seqGen 同毫秒内递增雪花序列号（跨毫秒自动归零）。
type seqGen struct {
	lastMs int64
	seq    int64
}

// next 返回指定毫秒对应的序列号；进入新毫秒时归零重计。
func (g *seqGen) next(ms int64) int64 {
	if ms != g.lastMs {
		g.lastMs, g.seq = ms, 0
	} else {
		g.seq++
	}
	return g.seq
}

// orUpdated 更新时间缺省回退为创建时间（旧库 updated_at 可能为空）。
func orUpdated(created time.Time, updated *time.Time) time.Time {
	if updated != nil {
		return *updated
	}
	return created
}

// parseNullableTime 解析可为空的时间字段：NULL 哨兵/空串返回 nil，解析失败静默返回 nil。
func parseNullableTime(f string) *time.Time {
	if f == "\x00NULL\x00" || f == "" {
		return nil
	}
	if t, err := parseTime(f); err == nil {
		return &t
	}
	return nil
}

// ----------------------------------------------------------------------------
// 主流程
// ----------------------------------------------------------------------------

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "导入失败:", err)
		os.Exit(1)
	}
}

func run() error {
	source := "data/backup.sql"
	if len(os.Args) > 1 {
		source = os.Args[1]
	}

	// 1. pg_restore 转 plain SQL
	sqlText, err := dumpPlainSQL(source)
	if err != nil {
		return err
	}

	// 2. 解析 COPY 块
	blocks := parseCopyBlocks(sqlText)
	byTable := map[string]copyBlock{}
	for _, b := range blocks {
		byTable[b.table] = b
	}
	for _, t := range []string{"xf_color", "xf_location", "xf_link_list"} {
		if _, ok := byTable[t]; !ok {
			return fmt.Errorf("备份中缺少数据表 %s", t)
		}
	}

	// 3. 解析为结构体
	colors, err := parseColors(byTable["xf_color"])
	if err != nil {
		return err
	}
	groups, err := parseLocations(byTable["xf_location"])
	if err != nil {
		return err
	}
	links, err := parseLinks(byTable["xf_link_list"])
	if err != nil {
		return err
	}
	sort.Slice(links, func(i, j int) bool { return links[i].id < links[j].id })

	// 4. 连接数据库
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		envOr("DATABASE_HOST", "localhost"),
		envOr("DATABASE_PORT", "5432"),
		envOr("DATABASE_USER", "bamboo_main"),
		envOr("DATABASE_PASS", "bamboo_main"),
		envOr("DATABASE_NAME", "bamboo_main"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 5. 事务导入
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM bm_link_friend").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM bm_link_group").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM bm_link_color").Error; err != nil {
			return err
		}

		// 5.1 颜色：旧 id -> 新雪花 ID
		colorMap := make(map[int64]xSnowflake.SnowflakeID, len(colors))
		gen := &seqGen{}
		for _, c := range colors {
			nid := snowflakeID(c.createdAt, constants.GeneLinkColor, gen.next(c.createdAt.UnixMilli()))
			colorMap[c.id] = nid
			primary := "#" + strings.TrimSpace(c.color)
			created := c.createdAt
			updated := orUpdated(created, c.updatedAt)
			if err := tx.Exec(
				`INSERT INTO bm_link_color (id, created_at, updated_at, name, primary_color, sub_color, hover_color, sort_order, status)
				 VALUES (?, ?, ?, ?, ?, NULL, NULL, ?, ?)`,
				nid, created, updated, c.displayName, primary, c.id, c.hasSelect,
			).Error; err != nil {
				return fmt.Errorf("插入颜色 %s 失败: %w", c.displayName, err)
			}
		}

		// 5.2 分组：旧 id -> 新雪花 ID
		groupMap := make(map[int64]xSnowflake.SnowflakeID, len(groups))
		gen = &seqGen{}
		for _, g := range groups {
			nid := snowflakeID(g.createdAt, constants.GeneLinkGroup, gen.next(g.createdAt.UnixMilli()))
			groupMap[g.id] = nid
			created := g.createdAt
			updated := orUpdated(created, g.updatedAt)
			var desc any
			if g.description != nil {
				desc = *g.description
			}
			if err := tx.Exec(
				`INSERT INTO bm_link_group (id, created_at, updated_at, name, description, sort_order, status)
				 VALUES (?, ?, ?, ?, ?, ?, ?)`,
				nid, created, updated, g.displayName, desc, g.sort+builtinGroupCount, g.reveal,
			).Error; err != nil {
				return fmt.Errorf("插入分组 %s 失败: %w", g.displayName, err)
			}
		}

		// 5.3 查询 super_admin（作者友链归属）
		var adminID xSnowflake.SnowflakeID
		var adminUser struct {
			ID xSnowflake.SnowflakeID
		}
		if err := tx.Raw(`SELECT id FROM bm_system_user WHERE username = 'super_admin' LIMIT 1`).Scan(&adminUser).Error; err != nil {
			return fmt.Errorf("查询 super_admin 失败: %w", err)
		}
		adminID = adminUser.ID

		// 5.4 友链（组内按旧 id 顺序赋 sort_order）
		curGroupSort := map[int64]int{}
		gen = &seqGen{}
		for _, l := range links {
			nid := snowflakeID(l.createdAt, constants.GeneLink, gen.next(l.createdAt.UnixMilli()))
			created := l.createdAt
			updated := orUpdated(created, l.updatedAt)
			var groupID, colorID, userID any
			if l.location != nil {
				if v, ok := groupMap[*l.location]; ok {
					groupID = v
				}
			}
			if l.color != nil {
				if v, ok := colorMap[*l.color]; ok {
					colorID = v
				}
			}
			if l.webmasterEmail != nil && *l.webmasterEmail == "gm@x-lf.cn" {
				userID = adminID
			}
			isFailure := 0
			if l.location != nil && *l.location == 2 {
				isFailure = 1
			}
			curSort := 0
			if l.location != nil {
				curGroupSort[*l.location]++
				curSort = curGroupSort[*l.location]
			}
			if err := tx.Exec(
				`INSERT INTO bm_link_friend (
					id, created_at, updated_at, name, url, avatar, rss, description, email,
					user_id, group_id, color_id, sort_order, status, is_failure, level,
					fail_reason, apply_remark, review_remark, screenshot_url, screenshot_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, NULL, ?, ?, NULL, NULL)`,
				nid, created, updated, l.siteName, l.siteURL, l.siteLogo, l.siteRSSURL, l.siteDesc, l.webmasterEmail,
				userID, groupID, colorID, curSort, l.status, isFailure, l.webmasterRemark, l.remark,
			).Error; err != nil {
				return fmt.Errorf("插入友链 %s 失败: %w", l.siteName, err)
			}
		}

		fmt.Printf("导入完成：颜色 %d 条、分组 %d 条、友链 %d 条\n", len(colors), len(groups), len(links))
		return nil
	})
}

// ----------------------------------------------------------------------------
// 解析辅助
// ----------------------------------------------------------------------------

func parseColors(b copyBlock) ([]oldColor, error) {
	var out []oldColor
	for _, row := range b.rows {
		if len(row) < 7 {
			return nil, fmt.Errorf("xf_color 行字段数不足: %d", len(row))
		}
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, err
		}
		created, err := parseTime(row[5])
		if err != nil {
			return nil, fmt.Errorf("解析 xf_color 时间失败: %w", err)
		}
		out = append(out, oldColor{
			id: id, name: row[1], displayName: row[2], color: row[3],
			hasSelect: row[4] == "t", createdAt: created, updatedAt: parseNullableTime(row[6]),
		})
	}
	return out, nil
}

func parseLocations(b copyBlock) ([]oldLocation, error) {
	var out []oldLocation
	for _, row := range b.rows {
		if len(row) < 8 {
			return nil, fmt.Errorf("xf_location 行字段数不足: %d", len(row))
		}
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, err
		}
		sort, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}
		created, err := parseTime(row[6])
		if err != nil {
			return nil, fmt.Errorf("解析 xf_location 时间失败: %w", err)
		}
		out = append(out, oldLocation{
			id: id, sort: sort, name: row[2], displayName: row[3],
			description: nullableStr(row[4]), reveal: row[5] == "t",
			createdAt: created, updatedAt: parseNullableTime(row[7]),
		})
	}
	return out, nil
}

func parseLinks(b copyBlock) ([]oldLink, error) {
	var out []oldLink
	for _, row := range b.rows {
		if len(row) < 21 {
			return nil, fmt.Errorf("xf_link_list 行字段数不足: %d", len(row))
		}
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			return nil, err
		}
		created, err := parseTime(row[18])
		if err != nil {
			return nil, fmt.Errorf("解析 xf_link_list 时间失败: %w", err)
		}
		status, err := strconv.Atoi(row[16])
		if err != nil {
			return nil, err
		}
		var location, color *int64
		if row[11] != "\x00NULL\x00" && row[11] != "" {
			if v, err := strconv.ParseInt(row[11], 10, 64); err == nil {
				location = &v
			}
		}
		if row[13] != "\x00NULL\x00" && row[13] != "" {
			if v, err := strconv.ParseInt(row[13], 10, 64); err == nil {
				color = &v
			}
		}
		out = append(out, oldLink{
			id: id, webmasterEmail: nullableStr(row[1]), siteName: row[3], siteURL: row[4],
			siteLogo: nullableStr(row[5]), cdnLogoURL: nullableStr(row[6]), siteDesc: nullableStr(row[7]),
			siteRSSURL: nullableStr(row[8]), location: location, color: color,
			webmasterRemark: nullableStr(row[14]), remark: nullableStr(row[15]),
			status: status, createdAt: created, updatedAt: parseNullableTime(row[19]),
		})
	}
	return out, nil
}

// dumpPlainSQL 调用 pg_restore 把 custom 格式备份转为 plain SQL 数据。
func dumpPlainSQL(source string) (string, error) {
	pgRestore, err := findPgRestore()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(pgRestore, "--data-only", "-f", "-", abs)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pg_restore 执行失败: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

// findPgRestore 查找 pg_restore 可执行文件。
func findPgRestore() (string, error) {
	if p := os.Getenv("PG_BIN"); p != "" {
		cand := filepath.Join(p, "pg_restore")
		if info, err := os.Stat(cand); err == nil && !info.IsDir() {
			return cand, nil
		}
	}
	if p, err := exec.LookPath("pg_restore"); err == nil {
		return p, nil
	}
	// macOS Homebrew libpq 常见路径
	for _, p := range []string{
		"/opt/homebrew/Cellar/libpq/18.3/bin/pg_restore",
		"/opt/homebrew/Cellar/libpq/18.2/bin/pg_restore",
		"/usr/local/opt/libpq/bin/pg_restore",
	} {
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("未找到 pg_restore，请安装 libpq 或设置 PG_BIN 环境变量")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
