package util

import "errors"

// ValidateUser 用户数据验证
func ValidateUser(username, password string) error {
	// 验证用户名不能为空，而且长度最多32个字符
	if username == "" || len(username) > 32 {
		return errors.New("用户名格式出错")
	}

	// 验证密码必须大于等于5字符，小于32个字符
	if len(password) < 5 || len(password) > 32 {
		return errors.New("密码格式出错")
	}

	return nil
}
