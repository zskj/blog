### Required
- go version 1.13.4
- mysql

###Ready
Create a **blog database** and import [SQL](https://github.com/zskj/blog/blob/master/docs/sql/blog.sql)

### Conf

You should modify `conf/app.ini`

```
[database]
Type = mysql
User = root
Password =
Host = 127.0.0.1:3306
Name = blog
TablePrefix = blog_
```
### Install
```
$ git clone https://github.com/zskj/blog.git
$ cd $path/blog
$ go mod tidy
$ go mod vendor
$ swag init
$ go run main.go
```
### Project information and existing API
```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /auth                     --> blog/routers/api.Auth (4 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (4 handlers)
[GIN-debug] GET    /api/v1/tags              --> blog/routers/api/v1.GetTags (5 handlers)
[GIN-debug] POST   /api/v1/tags              --> blog/routers/api/v1.AddTag (5 handlers)
[GIN-debug] PUT    /api/v1/tags/:id          --> blog/routers/api/v1.EditTag (5 handlers)
[GIN-debug] DELETE /api/v1/tags/:id          --> blog/routers/api/v1.DeleteTag (5 handlers)
[GIN-debug] GET    /api/v1/articles          --> blog/routers/api/v1.GetArticles (5 handlers)
[GIN-debug] GET    /api/v1/articles/:id      --> blog/routers/api/v1.GetArticle (5 handlers)
[GIN-debug] POST   /api/v1/articles          --> blog/routers/api/v1.AddArticle (5 handlers)
[GIN-debug] PUT    /api/v1/articles/:id      --> blog/routers/api/v1.EditArticle (5 handlers)
[GIN-debug] DELETE /api/v1/articles/:id      --> blog/routers/api/v1.DeleteArticle (5 handlers)
[GIN-debug] GET    /test                     --> blog/routers.InitRouter.func1 (4 handlers)

```
### Swaggo

> http://127.0.0.1:8000/swagger/index.html

![demo](https://github.com/zskj/blog/blob/master/docs/screenshots/swagger.png)
## 项目结构概览
```
├── conf 配置文件
├── docs：api文档swagger
│   └── sql：sql执行语句  
├── middleware：中间件
│   └── jwt：认证中间件
├── model：引用数据库模型
├── pkg：第三方包和公共模块
│   ├── app：gin engine
│   ├── e： 错误编码和错误信息
│   ├── logging：日志模块
│   ├── setting：go-ini包
│   └── util：工具库 
└── routers：路由处理
│    └── api：controller 逻辑梳理
│        └── v1：controller逻辑处理 
└── service：逻辑处理
└── runtime：日志，文件缓存存放
└── test：单元测试
└── main.go：入口文件 
```






