package v1

import (
	"blog/pkg/app"
	"blog/pkg/e"
	"bytes"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetCaptcha 获取验证码信息
// @Tags 公共接口
// @Summary 获取验证码信息
// @Success 200 {string} gin.Context.JSON
// @Router /api/v1/pub/captchaid [get]
func GetCaptcha(c *gin.Context) {
	id := captcha.NewLen(4) //todo 这里要把二维码长度放到配置文件
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"captchaId": id,
	})
}

// GetCaptcha 获取验证码信息
// @Tags 公共接口
// @Summary 响应图形验证码
// @Param captchaId query string true "验证码ID"
// @Success 200 "图形验证码"
// @Produce image/png
// @Failure 400 {string} gin.Context.JSON "{code:400,data:null,msg:无效的请求参数}"
// @Router /api/v1/pub/captcha [get]
func ResCaptcha(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Query("captchaId")
	if id == "" {
		appG.Response(http.StatusBadRequest, e.ERROR_NOT_EXIST_CAPTCHAID, nil)
		return
	}
	var w bytes.Buffer
	error := captcha.WriteImage(&w, id, 100, 50)
	if error != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_CAPTCHAID, nil)
		return
	}
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.Writer.Header().Set("Content-Type", "image/png")
	http.ServeContent(c.Writer, c.Request, id+".png", time.Time{}, bytes.NewReader(w.Bytes()))
	return
}


