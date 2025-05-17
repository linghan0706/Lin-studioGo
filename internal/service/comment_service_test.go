package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentRepository 模拟评论仓储
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	comment.ID = 1 // 模拟数据库分配ID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return args.Error(0)
}

func (m *MockCommentRepository) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentRepository) FindAll(ctx context.Context, filter repository.CommentFilter) ([]domain.Comment, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Comment), args.Get(1).(int64), args.Error(2)
}

func (m *MockCommentRepository) FindReplies(ctx context.Context, parentID uint) ([]domain.Comment, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func (m *MockCommentRepository) CountReplies(ctx context.Context, parentID uint) (int64, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockCommentRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCommentRepository) IncrementLikes(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// 更多方法实现...

// TestCreateComment 测试创建评论
func TestCreateComment(t *testing.T) {
	// 创建模拟仓储
	mockCommentRepo := new(MockCommentRepository)
	mockUserRepo := new(MockUserRepository)

	// 设置模拟行为
	userID := uint(1)
	mockUserRepo.On("FindByID", mock.Anything, userID).Return(&domain.User{
		ID:       userID,
		Username: "testuser",
		Role:     "user",
	}, nil)

	mockCommentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil)

	// 创建服务实例
	commentService := NewCommentService(mockCommentRepo, mockUserRepo)

	// 测试创建评论
	ctx := context.Background()
	comment, err := commentService.CreateComment(
		ctx,
		"测试评论内容",
		"article",
		uint(1),
		&userID,
		"",
		"",
		nil,
	)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "测试评论内容", comment.Content)
	assert.Equal(t, "article", comment.ItemType)
	assert.Equal(t, uint(1), comment.ItemID)
	assert.Equal(t, &userID, comment.UserID)
	assert.Equal(t, "pending", comment.Status) // 普通用户的评论默认为待审核状态

	// 验证调用
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

// TestCreateCommentByAdmin 测试管理员创建的评论应自动批准
func TestCreateCommentByAdmin(t *testing.T) {
	// 创建模拟仓储
	mockCommentRepo := new(MockCommentRepository)
	mockUserRepo := new(MockUserRepository)

	// 设置模拟行为
	userID := uint(1)
	mockUserRepo.On("FindByID", mock.Anything, userID).Return(&domain.User{
		ID:       userID,
		Username: "admin",
		Role:     "admin",
	}, nil)

	mockCommentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil)

	// 创建服务实例
	commentService := NewCommentService(mockCommentRepo, mockUserRepo)

	// 测试创建评论
	ctx := context.Background()
	comment, err := commentService.CreateComment(
		ctx,
		"管理员评论",
		"article",
		uint(1),
		&userID,
		"",
		"",
		nil,
	)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "approved", comment.Status) // 管理员的评论应自动批准

	// 验证调用
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

// TestGetComments 测试获取评论列表
func TestGetComments(t *testing.T) {
	// 创建模拟仓储
	mockCommentRepo := new(MockCommentRepository)
	mockUserRepo := new(MockUserRepository)

	// 模拟数据
	mockComments := []domain.Comment{
		{
			ID:       1,
			Content:  "评论1",
			UserID:   nil,
			ItemType: "article",
			ItemID:   1,
			Status:   "approved",
		},
		{
			ID:       2,
			Content:  "评论2",
			UserID:   nil,
			ItemType: "article",
			ItemID:   1,
			Status:   "approved",
		},
	}

	// 设置模拟行为
	mockCommentRepo.On("FindAll", mock.Anything, mock.AnythingOfType("repository.CommentFilter")).
		Return(mockComments, int64(2), nil)
	mockCommentRepo.On("CountReplies", mock.Anything, uint(1)).Return(int64(0), nil)
	mockCommentRepo.On("CountReplies", mock.Anything, uint(2)).Return(int64(0), nil)

	// 创建服务实例
	commentService := NewCommentService(mockCommentRepo, mockUserRepo)

	// 测试获取评论
	ctx := context.Background()
	status := "approved"
	comments, pagination, err := commentService.GetComments(ctx, "article", 1, nil, 1, 10, &status)

	// 断言
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, int64(2), pagination.Total)

	// 验证调用
	mockCommentRepo.AssertExpectations(t)
}

// 更多测试用例... 