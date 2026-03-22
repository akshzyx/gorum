package auth

import (
	"context"
	"errors"
	"time"

	"github.com/akshzyx/gorum/internal/util"
	"github.com/oklog/ulid/v2"
)

// Service handles auth business logic
type Service struct {
	repo        AuthRepository
	emailSender EmailSender
	tokenExpiry time.Duration
}

func NewService(repo AuthRepository, emailSender EmailSender) *Service {
	return &Service{
		repo:        repo,
		emailSender: emailSender,
		tokenExpiry: 24 * time.Hour,
	}
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) (SignupResponse, error) {
	// 1. Check duplicate email
	_, err := s.repo.GetUserByEmail(ctx, req.Email)

	switch {
	case err == nil:
		// user exists
		return SignupResponse{}, ErrEmailExists

	case errors.Is(err, ErrUserNotFound):
		// user does not exist → OK to continue

	default:
		// real DB error
		return SignupResponse{}, err
	}

	// 2. Hash password
	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return SignupResponse{}, err
	}

	// 3. Generate user ID
	userID := ulid.Make().String()

	// 4. Prepare verification token
	token := ulid.Make().String()
	vtoken := VerificationToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.tokenExpiry),
		Used:      false,
		CreatedAt: time.Now(),
	}

	// 5. Transactional create
	arg := CreateUserTxParams{
		ID:                userID,
		Username:          req.Username,
		Email:             req.Email,
		PasswordHash:      hashed,
		VerificationToken: vtoken,
	}

	createdUser, err := s.repo.CreateUserTx(ctx, arg)
	if err != nil {
		return SignupResponse{}, err
	}

	// 6. Send activation email
	if err := s.emailSender.SendVerificationEmail(req.Email, token); err != nil {
		return SignupResponse{}, err
	}

	return SignupResponse{
		UserID:   createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
	}, nil
}

func (s *Service) Activate(ctx context.Context, req ActivateRequest) error {
	// 1. Get token
	t, err := s.repo.GetVerificationToken(ctx, req.Token)
	if err != nil {
		return ErrTokenNotFound
	}

	// 2. Check expiration
	if time.Now().After(t.ExpiresAt) {
		return ErrTokenExpired
	}

	// 3. Check if already used
	if t.Used {
		return ErrTokenUsed
	}

	// 4. Mark token used
	if err := s.repo.MarkVerificationTokenUsed(ctx, req.Token); err != nil {
		return err
	}

	// 5. Activate user
	return s.repo.ActivateUser(ctx, t.UserID)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	if !util.CheckPasswordHash(req.Password, u.PasswordHash) {
		return LoginResponse{}, ErrInvalidCredentials
	}

	if !u.IsVerified {
		return LoginResponse{}, ErrEmailNotVerified
	}

	token, err := util.GenerateJWT(u.ID)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{Token: token}, nil
}

// resend activation
func (s *Service) ResendActivation(ctx context.Context, req ResendActivationRequest) error {
	// 1. Get user
	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return ErrUserNotFound
	}

	// 2. Already verified?
	if u.IsVerified {
		return ErrAlreadyVerified
	}

	// 3. New token
	token := ulid.Make().String()

	vtoken := VerificationToken{
		Token:     token,
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(s.tokenExpiry),
		Used:      false,
		CreatedAt: time.Now(),
	}

	// 4. Save new token
	if err := s.repo.CreateVerificationToken(ctx, vtoken); err != nil {
		return err
	}

	// 5. Send email
	return s.emailSender.SendVerificationEmail(req.Email, token)
}
