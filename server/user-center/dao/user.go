package dao

import (
	"context"
	"gorm.io/gorm"
	"user-center/model"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

// ExistOrNotByUserName 根据userName查询数据库中是否存在该名字
func (dao *UserDao) ExistOrNotByUserName(userName string) (user *model.User, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.User{}).Where("username=?", userName).First(&user).Count(&count).Error
	if count == 0 || err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return user, true, nil
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(user *model.User) (err error) {
	return dao.DB.Model(&model.User{}).Create(user).Error
}

// GetUserById 根据id获取user
func (dao *UserDao) GetUserById(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	return
}

// IsFollow 查询是否关注该用户
func (dao *UserDao) IsFollow(uId, followId uint) (bool, error) {
	var user model.User
	err := dao.Where("id = ?", uId).Preload("Follows", "id = ?", followId).Find(&user).Error
	if err != nil {
		return false, err
	}
	if len(user.Follows) > 0 {
		return true, nil
	}
	return false, nil
}
