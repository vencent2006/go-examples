package logic

import (
	"liwenzhou/bluebell/dao/mysql"
	"liwenzhou/bluebell/models"
	"liwenzhou/bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) error {
	// 1. 判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3. 保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
