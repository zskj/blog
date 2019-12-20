package routers

import (
	"blog/middleware"
	"blog/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"

	_ "blog/docs"
	"blog/middleware/jwt"
	"blog/routers/api"
	"blog/routers/api/v1"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())      //日志
	r.Use(middleware.Cors()) // 跨域请求
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode) //设置运行模式

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
		//注册文章路由
		RegisterArticleRouter(apiv1)
	}

	r.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})
	return r
}
