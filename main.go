/*
 * @Author: Monve
 * @Date: 2023-07-24 10:35:58
 * @LastEditors: Monve
 * @LastEditTime: 2023-07-25 09:33:40
 * @FilePath: /web-service-gin/main.go
 */
package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 将login返回的token填在这里.

import (
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"

	docs "web-service-gin/docs"

	"web-service-gin/controllers/record"
	"web-service-gin/controllers/token"
	"web-service-gin/controllers/user"
	"web-service-gin/utils/auth"
	"web-service-gin/utils/pgpool"
	"web-service-gin/utils/redis"
)

func main() {
	redis.Init()
	pgpool.Init()
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		user_g := v1.Group("/user")
		{
			user_g.POST("/login", user.LoginHandler)                          //登陆
			user_g.POST("/logout", auth.AuthMiddleware(), user.LogoutHandler) //登出
		}
		token_g := v1.Group("/token")
		{
			token_g.GET("/detail", auth.AuthMiddleware(), redis.CacheMiddleware, token.DetialHandler) //查询合约币
			token_g.POST("/add", auth.AuthMiddleware(), auth.RoleMiddleware([]string{"owener"}),      //添加
				token.AddHandler)
			token_g.POST("/edit", auth.AuthMiddleware(), auth.RoleMiddleware([]string{"owener"}), //编辑
				token.EditHandler)
			token_g.POST("/delete", auth.AuthMiddleware(), auth.RoleMiddleware([]string{"owener"}), //删除
				token.DeleteHandler)
		}
		record_g := v1.Group("/record")
		{
			record_g.GET("/stats", record.StatsHandler) //查询接口访问量
		}
	}

	// 添加 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 运行服务，监听端口8080
	r.Run(":8080")
}
