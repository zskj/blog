package routers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "blog/docs"
	"blog/middleware/jwt"
	"blog/pkg/setting"
	"blog/routers/api"
	"blog/routers/api/v1"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger()) //日志
	r.Use(Cors())       // 跨域请求

	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)

	r.POST("/auth", api.Auth)                                            //获取登录token
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //api注释文档

	apiv1 := r.Group("/api/v1")
	pub := apiv1.Group("/pub")
	{
		//获取验证码Id
		pub.GET("/login/captchaid", v1.GetCaptcha) 
		//获取验证码图片
		pub.GET("login/captcha", v1.ResCaptcha)
	}

	apiv1.Use(jwt.JWT()) //token 验证
	{
		//标签
		tag := apiv1.Group("/tags")
		{
			//列表
			tag.GET("", v1.GetTags)
			//新建
			tag.POST("", v1.AddTag)
			//更新
			tag.PUT(":id", v1.EditTag)
			//删除
			tag.DELETE(":id", v1.DeleteTag)
		}
		//文章
		article := apiv1.Group("/articles")
		{
			//列表
			article.GET("", v1.GetArticles)
			//文章ById
			article.GET(":id", v1.GetArticle)
			//新建
			article.POST("", v1.AddArticle)
			//更新ById
			article.PUT(":id", v1.EditArticle)
			//删除ById
			article.DELETE(":id", v1.DeleteArticle)
		}

	}

	r.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})
	return r
}

// 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
