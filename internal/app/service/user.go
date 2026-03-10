package service

import (
	"app/internal/app/repository"
	"app/internal/dto/request"
	resp "app/internal/dto/response"
	"app/internal/model"
	"app/pkg/config"
	rds "app/pkg/database/redis"
	"app/pkg/middleware"
	"app/pkg/toolkit"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SUser interface {
	Login(ctx context.Context, req request.Login) (*resp.RespLogin, string, time.Duration, error)
	RefreshToken(ctx context.Context, cookieValue string) (*resp.RespLogin, string, time.Duration, error)
	Logout(ctx context.Context, cookieValue string, jti string, userID int) error
}

type sUser struct {
	rUser repository.RUser
	rAuth repository.RAuth
	cfg   *config.Config
	rds   rds.Redis
}

func NewSUser(rUser repository.RUser, rAuth repository.RAuth, cfg *config.Config, rds rds.Redis) SUser {
	return &sUser{
		rUser: rUser,
		rAuth: rAuth,
		cfg:   cfg,
		rds:   rds,
	}
}

/*-------------------------- Main Function --------------------------*/

func (s *sUser) Login(ctx context.Context, req request.Login) (*resp.RespLogin, string, time.Duration, error) {
	lockKey := fmt.Sprintf("login:lock:%s", req.Username)
	failKey := fmt.Sprintf("login:fail:%s", req.Username)

	if err := s.checkLock(ctx, lockKey, req.Username); err != nil {
		return nil, "", 0, err
	}

	user, err := s.rUser.FindByUsernameOrEmail(ctx, req.Username, req.Username)
	if err != nil {
		return nil, "", 0, err
	}

	if !toolkit.CheckPassword(user.Password, req.Password) {
		if err := s.checkFail(ctx, failKey, lockKey, req.Username); err != nil {
			return nil, "", 0, err
		}

		return nil, "", 0, fmt.Errorf("invalid password or username")
	}

	// generate short-lived access token
	jti := uuid.New().String()
	accessTTL := 15 * time.Minute

	token, err := middleware.GenerateToken(user.ID, jti, s.cfg.PrivateKey, accessTTL)
	if err != nil {
		return nil, "", 0, err
	}

	session := resp.SessionRecord{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Email:    user.Email,
	}

	if err := s.saveSession(ctx, jti, &session, accessTTL); err != nil {
		return nil, "", 0, err
	}

	refreshToken, refreshTTL, err := s.createRefreshToken(ctx, user.ID, jti, req.RememberMe)
	if err != nil {
		return nil, "", 0, err
	}

	return &resp.RespLogin{
		AccessToken: token,
		ExpiresIn:   int64(accessTTL.Seconds()),
	}, refreshToken, refreshTTL, nil
}

func (s *sUser) RefreshToken(ctx context.Context, cookieValue string) (*resp.RespLogin, string, time.Duration, error) {
	tokenID, secret, err := splitRefreshCookie(cookieValue)
	if err != nil {
		return nil, "", 0, err
	}

	tokenRecord, err := s.rAuth.FindByTokenID(ctx, tokenID)
	if err != nil {
		return nil, "", 0, fmt.Errorf("invalid refresh token")
	}

	if tokenRecord.Revoked || time.Now().After(tokenRecord.ExpiresAt) {
		return nil, "", 0, fmt.Errorf("refresh token expired or revoked")
	}

	// validate secret against stored hash
	if !toolkit.CheckPassword(tokenRecord.TokenHash, secret) {
		return nil, "", 0, fmt.Errorf("invalid refresh token")
	}

	// load user
	user, err := s.rUser.FindByID(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, "", 0, fmt.Errorf("user not found")
	}

	// rotate refresh token and issue new access token
	newJTI := uuid.New().String()
	accessTTL := 15 * time.Minute

	accessToken, err := middleware.GenerateToken(user.ID, newJTI, s.cfg.PrivateKey, accessTTL)
	if err != nil {
		return nil, "", 0, err
	}

	session := resp.SessionRecord{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Email:    user.Email,
	}

	if err := s.saveSession(ctx, newJTI, &session, accessTTL); err != nil {
		return nil, "", 0, err
	}

	newRefreshToken, refreshTTL, err := s.rotateRefreshToken(ctx, &tokenRecord, newJTI)
	if err != nil {
		return nil, "", 0, err
	}

	return &resp.RespLogin{
		AccessToken: accessToken,
		ExpiresIn:   int64(accessTTL.Seconds()),
	}, newRefreshToken, refreshTTL, nil
}

func (s *sUser) Logout(ctx context.Context, cookieValue string, jti string, userID int) error {
	// best-effort refresh token revocation
	if cookieValue != "" {
		if tokenID, _, err := splitRefreshCookie(cookieValue); err == nil {
			if tokenRecord, err := s.rAuth.FindByTokenID(ctx, tokenID); err == nil {
				_ = s.rAuth.RevokeToken(ctx, tokenRecord.ID, nil)
			}
		}
	}

	// optionally revoke all tokens for user
	if err := s.rAuth.RevokeAllByUserID(ctx, userID); err != nil {
		return err
	}

	// delete redis session
	if jti != "" {
		key := fmt.Sprintf("session:%s", jti)
		_ = s.rds.Delete(ctx, key)
	}

	return nil
}

/*-------------------------- Helper Function --------------------------*/

func (s *sUser) checkLock(ctx context.Context, key, username string) error {
	value, err := s.rds.Exists(ctx, key)
	if err != nil {
		return err
	}

	if value {
		return fmt.Errorf("user %s is locked, wait 5 minutes", username)
	}
	return nil
}

func (s *sUser) checkFail(ctx context.Context, failKey, lockKey, username string) error {
	value, err := s.rds.Incr(ctx, failKey)
	if err != nil {
		return err
	}

	if value > 3 {
		if err := s.rds.Expire(ctx, lockKey, time.Minute*5); err != nil {
			return err
		}

		if err := s.rds.Delete(ctx, failKey); err != nil {
			return err
		}

		return fmt.Errorf("user %s is locked", username)
	}

	return nil
}

func (s *sUser) saveSession(ctx context.Context, jti string, session *resp.SessionRecord, ttl time.Duration) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("session:%s", jti)
	if err := s.rds.SetWithDuration(ctx, key, string(sessionJSON), ttl); err != nil {
		return err
	}

	return nil
}

// createRefreshToken creates a new refresh token for a login.
func (s *sUser) createRefreshToken(ctx context.Context, userID int, jti string, rememberMe bool) (string, time.Duration, error) {
	var refreshTTL time.Duration
	if rememberMe {
		refreshTTL = 30 * 24 * time.Hour
	} else {
		refreshTTL = 7 * 24 * time.Hour
	}

	tokenID := uuid.New().String()
	secret := uuid.New().String()
	hash, err := toolkit.HashPassword(secret)
	if err != nil {
		return "", 0, err
	}

	rt := &model.RefreshToken{
		UserID:    userID,
		TokenID:   tokenID,
		TokenHash: hash,
		JTI:       jti,
		ExpiresAt: time.Now().Add(refreshTTL),
	}

	if err := s.rAuth.CreateRefreshToken(ctx, rt); err != nil {
		return "", 0, err
	}

	return buildRefreshCookieValue(tokenID, secret), refreshTTL, nil
}

// rotateRefreshToken revokes the old token and issues a new one.
func (s *sUser) rotateRefreshToken(ctx context.Context, old *model.RefreshToken, newJTI string) (string, time.Duration, error) {
	remaining := time.Until(old.ExpiresAt)
	if remaining <= 0 {
		remaining = 7 * 24 * time.Hour
	}

	tokenID := uuid.New().String()
	secret := uuid.New().String()
	hash, err := toolkit.HashPassword(secret)
	if err != nil {
		return "", 0, err
	}

	newToken := &model.RefreshToken{
		UserID:    old.UserID,
		TokenID:   tokenID,
		TokenHash: hash,
		JTI:       newJTI,
		ExpiresAt: time.Now().Add(remaining),
	}

	if err := s.rAuth.CreateRefreshToken(ctx, newToken); err != nil {
		return "", 0, err
	}

	if err := s.rAuth.RevokeToken(ctx, old.ID, &newToken.ID); err != nil {
		return "", 0, err
	}

	return buildRefreshCookieValue(tokenID, secret), remaining, nil
}

func buildRefreshCookieValue(tokenID, secret string) string {
	return tokenID + ":" + secret
}

func splitRefreshCookie(value string) (string, string, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid refresh token format")
	}
	return parts[0], parts[1], nil
}

