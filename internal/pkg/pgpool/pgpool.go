/*
 * @Author: Monve
 * @Date: 2023-07-25 03:00:01
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 12:17:19
 * @FilePath: /web-service-gin/internal/pkg/pgpool/pgpool.go
 */
package pgpool

import (
	"database/sql"
	"fmt"
	"sync"
	"web-service-gin/configs/env"

	_ "github.com/lib/pq" // 导入PostgreSQL驱动，这里使用匿名导入
)

type DBPool struct {
	mu     sync.Mutex
	dbConn *sql.DB
}

var pool *DBPool

var db_url string

func Init() {
	db_url = env.PostgresUrl
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// 连接测试
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("【PostgreSQL】连接成功")

	pool = &DBPool{
		dbConn: db,
	}
}

func GetDBConn() *sql.DB {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		fmt.Println("error", "error when GetDBConn", err)
	}
	pool.dbConn = db

	return db

}
