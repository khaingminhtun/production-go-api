package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/khaingminhtun/production-go-api/internal/features/user"
	redisinfra "github.com/khaingminhtun/production-go-api/internal/infrastructure/redis"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
)

// ============================================================
// Mock User Repository
// ============================================================

type MockRegistrationUserRepository struct {
	mock.Mock
}

func (m *MockRegistrationUserRepository) Create(ctx context.Context, user *user.User) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationUserRepository) GetByID(ctx context.Context, id uint) (*user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationUserRepository) List(ctx context.Context, offset, limit int) ([]user.User, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationUserRepository) Update(ctx context.Context, user *user.User) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationUserRepository) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	args := m.Called(ctx, email)

	var u *user.User

	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}

	return u, args.Error(1)
}

func (m *MockRegistrationUserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*user.User, error) {
	args := m.Called(ctx, username)

	var u *user.User

	if args.Get(0) != nil {
		u = args.Get(0).(*user.User)
	}

	return u, args.Error(1)
}

// ============================================================
// Mock Redis Store
// ============================================================

type MockRegistrationStore struct {
	mock.Mock
}

func (m *MockRegistrationStore) Get(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationStore) Exists(ctx context.Context, key string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRegistrationStore) Set(
	ctx context.Context,
	key string,
	value string,
	ttl time.Duration,
) error {
	args := m.Called(ctx, key, value, ttl)

	return args.Error(0)
}

func (m *MockRegistrationStore) Delete(
	ctx context.Context,
	key string,
) error {
	args := m.Called(ctx, key)

	return args.Error(0)
}

// ============================================================
// Mock Email Publisher
// ============================================================

type MockEmailPublisher struct {
	mock.Mock
}

func (m *MockEmailPublisher) Consume(ctx context.Context, consumerName string, count int, block time.Duration) ([]redisinfra.EmailJobMessage, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockEmailPublisher) Ack(ctx context.Context, messageID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockEmailPublisher) EnsureConsumerGroup(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockEmailPublisher) Publish(
	ctx context.Context,
	job redisinfra.EmailJob,
) error {
	args := m.Called(ctx, job)

	return args.Error(0)
}

// ============================================================
// Helper
// ============================================================

func userNotFoundError() error {
	return apperror.New(
		apperror.CodeUserNotFound,
		"user not found",
		nil,
	)
}

// ============================================================
// Register - Success
// ============================================================

func TestService_Register_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	// --------------------------------------------------------
	// User email does not exist
	// --------------------------------------------------------

	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// --------------------------------------------------------
	// Username does not exist
	// --------------------------------------------------------

	userRepo.
		On(
			"GetByUsername",
			ctx,
			"john",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// --------------------------------------------------------
	// Redis SET succeeds
	// --------------------------------------------------------

	redisStore.
		On(
			"Set",
			ctx,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			10*time.Minute,
		).
		Return(nil)

	// --------------------------------------------------------
	// Email queue succeeds
	// --------------------------------------------------------

	emailQueue.
		On(
			"Publish",
			ctx,
			mock.AnythingOfType("redis.EmailJob"),
		).
		Return(nil)

	// --------------------------------------------------------
	// Execute
	// --------------------------------------------------------

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: " john ",
			Email:    " TEST@example.com ",
			Password: "Password123!",
		},
	)

	// --------------------------------------------------------
	// Assertions
	// --------------------------------------------------------

	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotEmpty(t, result.RegistrationID)

	require.Equal(
		t,
		"Verification code sent to your email",
		result.Message,
	)

	// Verify all expected calls happened.
	userRepo.AssertExpectations(t)
	redisStore.AssertExpectations(t)
	emailQueue.AssertExpectations(t)
}

// ============================================================
// Register - Email Already Exists
// ============================================================

func TestService_Register_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	existingUser := &user.User{
		ID:    1,
		Email: "test@example.com",
	}

	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(existingUser, nil)

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: "john",
			Email:    "test@example.com",
			Password: "Password123!",
		},
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserAlreadyExists,
		),
	)

	userRepo.AssertExpectations(t)

	// Registration must stop here.
	redisStore.AssertNotCalled(
		t,
		"Set",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)

	emailQueue.AssertNotCalled(
		t,
		"Publish",
		mock.Anything,
		mock.Anything,
	)
}

// ============================================================
// Register - Username Already Exists
// ============================================================

func TestService_Register_UsernameAlreadyExists(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	// Email does not exist.
	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// Username already exists.
	userRepo.
		On(
			"GetByUsername",
			ctx,
			"john",
		).
		Return(
			&user.User{
				ID:       1,
				Username: "john",
			},
			nil,
		)

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: "john",
			Email:    "test@example.com",
			Password: "Password123!",
		},
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserAlreadyExists,
		),
	)

	userRepo.AssertExpectations(t)

	// Registration must stop.
	redisStore.AssertNotCalled(
		t,
		"Set",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)

	emailQueue.AssertNotCalled(
		t,
		"Publish",
		mock.Anything,
		mock.Anything,
	)
}

// ============================================================
// Register - GetByEmail Unexpected Error
// ============================================================

func TestService_Register_GetByEmailError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	repoErr := errors.New("database connection failed")

	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(nil, repoErr)

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: "john",
			Email:    "test@example.com",
			Password: "Password123!",
		},
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.ErrorIs(t, err, repoErr)

	userRepo.AssertExpectations(t)

	redisStore.AssertNotCalled(
		t,
		"Set",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)

	emailQueue.AssertNotCalled(
		t,
		"Publish",
		mock.Anything,
		mock.Anything,
	)
}

// ============================================================
// Register - GetByUsername Unexpected Error
// ============================================================

func TestService_Register_GetByUsernameError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	repoErr := errors.New("database connection failed")

	// Email does not exist.
	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// Username lookup fails.
	userRepo.
		On(
			"GetByUsername",
			ctx,
			"john",
		).
		Return(nil, repoErr)

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: "john",
			Email:    "test@example.com",
			Password: "Password123!",
		},
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.ErrorIs(t, err, repoErr)

	userRepo.AssertExpectations(t)

	redisStore.AssertNotCalled(
		t,
		"Set",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)

	emailQueue.AssertNotCalled(
		t,
		"Publish",
		mock.Anything,
		mock.Anything,
	)
}

// ============================================================
// Register - Redis Set Error
// ============================================================

func TestService_Register_RedisSetError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	// Email does not exist.
	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// Username does not exist.
	userRepo.
		On(
			"GetByUsername",
			ctx,
			"john",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	redisErr := errors.New("redis unavailable")

	redisStore.
		On(
			"Set",
			ctx,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			10*time.Minute,
		).
		Return(redisErr)

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: "john",
			Email:    "test@example.com",
			Password: "Password123!",
		},
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.ErrorIs(t, err, redisErr)

	userRepo.AssertExpectations(t)
	redisStore.AssertExpectations(t)

	// Email must not be queued when Redis fails.
	emailQueue.AssertNotCalled(
		t,
		"Publish",
		mock.Anything,
		mock.Anything,
	)
}

// ============================================================
// Register - Email Queue Error
// ============================================================

func TestService_Register_EmailQueueError(t *testing.T) {
	ctx := context.Background()

	userRepo := new(MockRegistrationUserRepository)
	redisStore := new(MockRegistrationStore)
	emailQueue := new(MockEmailPublisher)

	svc := NewService(
		userRepo,
		redisStore,
		emailQueue,
	)

	// Email does not exist.
	userRepo.
		On(
			"GetByEmail",
			ctx,
			"test@example.com",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// Username does not exist.
	userRepo.
		On(
			"GetByUsername",
			ctx,
			"john",
		).
		Return(
			nil,
			userNotFoundError(),
		)

	// Redis registration succeeds.
	redisStore.
		On(
			"Set",
			ctx,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			10*time.Minute,
		).
		Return(nil)

	emailErr := errors.New("email queue unavailable")

	// Email queue fails.
	emailQueue.
		On(
			"Publish",
			ctx,
			mock.AnythingOfType("redis.EmailJob"),
		).
		Return(emailErr)

	// Service should clean up Redis.
	redisStore.
		On(
			"Delete",
			ctx,
			mock.AnythingOfType("string"),
		).
		Return(nil)

	result, err := svc.Register(
		ctx,
		RegisterRequest{
			Username: "john",
			Email:    "test@example.com",
			Password: "Password123!",
		},
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.ErrorIs(t, err, emailErr)

	userRepo.AssertExpectations(t)
	redisStore.AssertExpectations(t)
	emailQueue.AssertExpectations(t)
}
