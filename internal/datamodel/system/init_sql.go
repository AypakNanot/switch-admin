package system

import _ "embed"

//go:embed data/init.sql
var embeddedInitSQL string

// GetEmbeddedInitSQL 返回嵌入的初始化 SQL 脚本
func GetEmbeddedInitSQL() string {
	return embeddedInitSQL
}
