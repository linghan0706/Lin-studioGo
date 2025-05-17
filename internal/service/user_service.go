package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	
	"gorm.io/datatypes"
)

// UserService 用户服务接口
type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uint, bio string, socialLinks, contactInfo map[string]interface{}) (*domain.User, error)
	UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error)
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// GetProfile 获取用户信息
func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	
	return user, nil
}

// UpdateProfile 更新用户信息
func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uint, bio string, socialLinks, contactInfo map[string]interface{}) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	
	// 更新字段
	user.Bio = sql.NullString{
		String: bio,
		Valid: bio != "",
	}
	
	// 将map转换为JSON
	socialLinksJSON, err := json.Marshal(socialLinks)
	if err != nil {
		return nil, err
	}
	user.SocialLinks = datatypes.JSON(socialLinksJSON)
	
	contactInfoJSON, err := json.Marshal(contactInfo)
	if err != nil {
		return nil, err
	}
	user.ContactInfo = datatypes.JSON(contactInfoJSON)
	
	// 保存更新
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// UploadAvatar 上传头像
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	
	if user == nil {
		return "", errors.New("用户不存在")
	}
	
	// 确保上传目录存在
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}
	
	// 创建文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, user.UpdatedAt.Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)
	
	// 保存文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	// 将源文件内容复制到目标文件
	if _, err = dst.ReadFrom(src); err != nil {
		return "", err
	}
	
	// 更新用户头像
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	user.Avatar = avatarURL
	
	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}
	
	return avatarURL, nil
} 