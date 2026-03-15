package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite"
)

func main() {
	// 直接打开 SQLite 数据库
	db, err := sql.Open("sqlite3", "data/admin.db")
	if err != nil {
		log.Fatalf("打开数据库失败：%v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}

	fmt.Println("数据库连接成功！")

	// 读取 SQL 文件
	sqlContent, err := ioutil.ReadFile("scripts/migrate_v2.sql")
	if err != nil {
		log.Fatalf("读取 SQL 文件失败：%v", err)
	}

	// 预处理 SQL 内容，移除注释
	cleanedSQL := removeComments(string(sqlContent))

	// 分割 SQL 语句（按分号分割）
	statements := strings.Split(cleanedSQL, ";")

	executed := 0
	skipped := 0
	errors := 0

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := db.Exec(stmt)
		if err != nil {
			// 忽略已存在的错误（如表已存在）
			errMsg := err.Error()
			if strings.Contains(errMsg, "already exists") ||
				strings.Contains(errMsg, "duplicate column") ||
				strings.Contains(errMsg, "UNIQUE constraint failed") ||
				strings.Contains(errMsg, "duplicate key") {
				fmt.Printf("跳过（已存在）: %s\n", getStatementType(stmt))
				skipped++
				continue
			}
			log.Printf("执行 SQL 失败：%v\n  语句：%s\n", err, getStatementType(stmt))
			errors++
			continue
		}
		executed++
		fmt.Printf("执行成功：%s\n", getStatementType(stmt))
	}

	fmt.Printf("\n==========================\n")
	fmt.Printf("迁移完成！\n")
	fmt.Printf("  成功：%d 条\n", executed)
	fmt.Printf("  跳过：%d 条\n", skipped)
	fmt.Printf("  错误：%d 条\n", errors)
	fmt.Printf("==========================\n")
}

func getStatementType(stmt string) string {
	stmt = strings.TrimSpace(stmt)
	words := strings.Fields(stmt)
	if len(words) > 2 {
		return fmt.Sprintf("%s %s", words[0], words[1])
	}
	if len(words) > 0 {
		return words[0]
	}
	return "UNKNOWN"
}

// removeComments 移除 SQL 注释
func removeComments(sql string) string {
	// 移除单行注释 -- ...
	re := regexp.MustCompile(`(?m)--.*$`)
	sql = re.ReplaceAllString(sql, "")

	// 移除多行注释 /* ... */
	re = regexp.MustCompile(`/\*[\s\S]*?\*/`)
	sql = re.ReplaceAllString(sql, "")

	return sql
}
