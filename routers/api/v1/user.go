package v1

import (
	"blog/models"
	"github.com/dchest/captcha"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"blog/pkg/app"
	"blog/pkg/e"
	"blog/pkg/util"
)

type auth struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaCode string `json:"captcha_code"`
	CaptchaId   string `json:"captcha_id"`
}

type currentUser struct {
	Id string `json:"id"`

	Username string `json:"username"`
}

// @Summary   注册用户
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   models.Reg   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/reg  [POST]
func Reg(c *gin.Context) {

	appG := app.Gin{C: c}
	var reqInfo models.Reg
	var data interface{}
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}
	valid := validation.Validation{}
	valid.Required(reqInfo.Username, "username").Message("请输入用户名")
	valid.Required(reqInfo.Password, "password").Message("输入密码")
	valid.Required(reqInfo.PasswordAgain, "password_again").Message("输入去确认密码")
	valid.MaxSize(reqInfo.Username, 50, "username").Message("密码最长为50字符")
	valid.MaxSize(reqInfo.Password, 50, "password").Message("密码最长为50字符")
	valid.MaxSize(reqInfo.PasswordAgain, 50, "password_again").Message("确认密码最长为50字符")
	valid.Required(reqInfo.CaptchaCode, "captcha_code").Message("请您输入验证码")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, valid.Errors)
		return
	}
	if !captcha.VerifyString(reqInfo.CaptchaId, reqInfo.CaptchaCode) {
		appG.Response(http.StatusInternalServerError, e.ERROR_CAPTCHA_FAIL, data)
		return
	}
	if _, isExist := models.FindUserByUsername(reqInfo.Username); isExist != gorm.ErrRecordNotFound {
		appG.Response(http.StatusInternalServerError, e.ERROR_REPEAT_NAME, data)
		return
	}
	var newUser models.User
	newUser.Username = reqInfo.Username
	newUser.Password = util.EncodeMD5(reqInfo.Password)
	newUser.Status = 1
	newUser.Secret = util.RandStringBytesMaskImprSrcUnsafe(5)
	newUser.CreatedOn = int(time.Now().Unix())
	newUser.ModifiedOn = int(time.Now().Unix())
	userId, isSuccess := models.NewUser(&newUser)
	if userId > 0 {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"id": userId})
		return
	}
	appG.Response(http.StatusOK, e.ERROR_ADD_FAIL, isSuccess)

}

// @Summary   用户登录 获取token 信息
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Param   body  body   models.AuthSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router /api/v1/auth  [POST]
func Auth(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo auth
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}
	valid.Required(reqInfo.Username, "username").Message("请输入用户名")
	valid.Required(reqInfo.Password, "password").Message("输入密码")
	valid.MaxSize(reqInfo.Username, 50, "username").Message("最长为50字符")
	valid.MaxSize(reqInfo.Password, 50, "password").Message("用户名最长为50字符")
	valid.Required(reqInfo.CaptchaCode, "captcha_code").Message("请您输入验证码")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_FAIL, valid.Errors)
		return
	}
	if !captcha.VerifyString(reqInfo.CaptchaId, reqInfo.CaptchaCode) {
		appG.Response(http.StatusInternalServerError, e.ERROR_CAPTCHA_FAIL, valid.Errors)
		return
	}

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	isExist, user, err := models.LoginCheck(reqInfo.Username, util.EncodeMD5(reqInfo.Password))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(user)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

// @Summary 获取登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/currentuser   [GET]
func CurrentUser(c *gin.Context) {
	var code int
	var data interface{}
	var user currentUser
	appG := app.Gin{C: c}
	code = e.SUCCESS
	Authorization := c.GetHeader("Authorization") //在header中存放token
	token := strings.Split(Authorization, " ")
	//token := c.Query("token")
	if Authorization == "" {
		code = e.INVALID_PARAMS
	} else {
		claims, err := util.ParseToken(token[0])
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			default:
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
		}
		user.Id = claims.Id
		user.Username = claims.Audience
	}

	if code != e.SUCCESS {
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	} else {
		appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"data": user,
		})
	}

}

// @Summary 刷新token
// @Tags 用户管理
// @Accept json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/refreshtoken  [GET]
func RefreshToken(c *gin.Context) {
	var data interface{}
	var code int
	appG := app.Gin{C: c}
	code = e.SUCCESS
	Authorization := c.GetHeader("Authorization") //在header中存放token
	if Authorization == "" {
		code = e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": data,
		})
	}
	token, err := util.RefreshToken(Authorization)
	if err != nil {
		code = e.INVALID_PARAMS
		appG.Response(http.StatusOK, code, map[string]interface{}{
			"data": err,
		})
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})

}
