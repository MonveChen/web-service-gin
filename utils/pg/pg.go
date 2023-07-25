/*
 * @Author: Monve
 * @Date: 2023-07-24 16:31:58
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-24 16:39:54
 * @FilePath: /web-service-gin/utils/pg/pg.go
 */

package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // 导入PostgreSQL驱动，这里使用匿名导入
)

var Db *sql.DB

func Init() {
	// 数据库连接字符串
	var err error
	Db, err = sql.Open("postgres", "user=root password=123456 dbname=postgres host=server1.dibiaozuitu.com port=5432 sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	// 连接测试
	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("【PostgreSQL】连接成功")
}

func Example() {
	// 查询数据
	rows, err := Db.Query("SELECT id, name, age FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		var id int
		var name string
		var age int

		err = rows.Scan(&id, &name, &age)
		if err != nil {
			panic(err)
		}

		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	// 检查查询时是否出现错误
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}
