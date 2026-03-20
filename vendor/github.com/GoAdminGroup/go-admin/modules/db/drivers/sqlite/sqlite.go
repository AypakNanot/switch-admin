// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

// Package sqlite 提供纯 Go SQLite 驱动支持（不需要 CGO）
// 使用 modernc.org/sqlite 替代 github.com/mattn/go-sqlite3
package sqlite

import (
	"database/sql"
	"database/sql/driver"

	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动，不需要 CGO
)

// 驱动别名注册 - 将 "sqlite" 包装为 "sqlite3"
// 因为 modernc.org/sqlite 注册为 "sqlite"，但 GoAdmin 使用 "sqlite3"
type sqlite3Driver struct{}

func (d sqlite3Driver) Open(dsn string) (driver.Conn, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// 直接返回内部驱动的连接
	return db.Driver().Open(dsn)
}

func init() {
	sql.Register("sqlite3", sqlite3Driver{})
}
