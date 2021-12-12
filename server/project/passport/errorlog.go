package passport

import "errors"

//
//注册模块错误
//

//用户名被注册错误
var errTheSameUserNameExpection = errors.New("用户名被注册")

//邮箱被注册错误
var errTheSameUserEmailExpection = errors.New("邮箱被注册")

//
//登录模块错误
//

//密码错误
var errWrongUserNameOrPasswd = errors.New("账号或密码错误")
