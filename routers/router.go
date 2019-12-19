package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"

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
