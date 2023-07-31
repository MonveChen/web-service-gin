###  题目：

##### 使用golang实现一个服务，给其他服务提供一个token信息查询的接口，需要返回token的symbol和decimal

要求：

 	1.需要实现已查询的token信息存储到数据库里面(需要对token信息进行校验，非法token地址请求不允许通过)，并且查询接口实现缓存（需要使用合理的数据结构）
	2.对API的访问进行限制，需要实现jwt token验证的逻辑，通过jwt信息确认访问者的身份（通过key/sercert获取token），按照不同的身份对用户进行不同的限流(需要考虑多节点的问题)
	3.需要实现手动管理和更新/移除token信息的接口
	4.实现接口访问统计，以天为维度统计每一个访问者每天的访问次数，实现按（时间区间）查询访问情况的API





### 记录：

##### 建表

```sql
CREATE TABLE [IF NOT EXISTS] public.contact_token (
  	id SERIAL PRIMARY KEY,
   	chainId VARCHAR(20) NOT NULL,
		token VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NULL,
    decimals SMALLINT NULL
);

CREATE UNIQUE INDEX idx_main ON public.contact_token (chainId, token);



```



##### swagger

http://localhost:8080/swagger/index.html