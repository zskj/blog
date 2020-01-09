package routers

import (
	"blog/controller/api/v1"
	"github.com/gin-gonic/gin"
)

//注册文章路由
func RegisterArticleRouter(router *gin.RouterGroup) {
	articleRouter := router.Group("/articles")
	{
			//列表
			articleRouter.GET("", v1.GetArticles)
			//文章ById
			articleRouter.GET(":id", v1.GetArticle)
			//新建
			articleRouter.POST("", v1.AddArticle)
			//更新ById
			articleRouter.PUT(":id", v1.EditArticle)
			//删除ById
			articleRouter.DELETE(":id", v1.DeleteArticle)

	}
}
