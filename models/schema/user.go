package schema

//用户表
type User struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Status     int    `json:"status"`
	CreatedOn  uint   `json:"created_on"`
	ModifiedOn uint   `json:"modified_on"`
	DeletedOn  uint   `json:"deleted_on"`
	Secret     string `json:"secret"`
}

//登录
type AuthSwag struct {
	Username    string `json:"username"` //登录账户
	Password    string `json:"password"` //登录密码
	CaptchaCode string `json:"captcha_code"`
	CaptchaId   string `json:"captcha_id"`
}

//修改密码
type PasswordSwag struct {
	OldPassword string `json:"old_password"` //旧密码
	NewPassword string `json:"new_password"` //新密码
}

//注册
type Reg struct {
	Username      string `json:"username" binding:"required"`        //用户名
	Password      string `json:"password"  binding:"required"`       //密码
	PasswordAgain string `json:"password_again" binding:"required" ` //确认密码
	CaptchaCode   string `json:"captcha_code" binding:"required"`    //验证码
	CaptchaId     string `json:"captcha_id"  binding:"required"`     //验证码Id
}
