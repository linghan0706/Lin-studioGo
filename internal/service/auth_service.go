package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/utils"
	"Lin_studio/internal/repository"
	"errors"
)

// AuthService 认证服务接口
type AuthService interface {
	Login(username, password string) (*domain.UserResponse, string, string, error)
	Register(user *domain.User) (*domain.UserResponse, error)
	RefreshToken(refreshToken string) (string, string, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Login 用户登录
func (s *authService) Login(username, password string) (*domain.UserResponse, string, string, error) {
	// 查询用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, "", "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return nil, "", "", errors.New("用户名或密码错误")
	}

	// 检查账户状态
	if user.Status != "active" {
		return nil, "", "", errors.New("账户已被禁用")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", errors.New("生成令牌失败")
	}

	// 生成刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", errors.New("生成刷新令牌失败")
	}

	// 更新最后登录时间
	_ = s.userRepo.UpdateLastLogin(user)

	// 转换为响应数据
	userResponse := user.ToResponse()

	return &userResponse, token, refreshToken, nil
}

// Register 用户注册
func (s *authService) Register(user *domain.User) (*domain.UserResponse, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(user.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 创建用户
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败: " + err.Error())
	}

	// 转换为响应数据
	userResponse := user.ToResponse()

	return &userResponse, nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新令牌
	userID, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("无效的刷新令牌")
	}

	// 查询用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", "", errors.New("用户不存在")
	}

	// 检查账户状态
	if user.Status != "active" {
		return "", "", errors.New("账户已被禁用")
	}

	// 生成新的JWT令牌
	newToken, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", errors.New("生成令牌失败")
	}

	// 生成新的刷新令牌
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.New("生成刷新令牌失败")
	}

	return newToken, newRefreshToken, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// 查询用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证当前密码
	if !user.CheckPassword(currentPassword) {
		return errors.New("当前密码错误")
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return errors.New("设置密码失败: " + err.Error())
	}

	// 更新用户
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("更新用户失败: " + err.Error())
	}

	return nil
} 