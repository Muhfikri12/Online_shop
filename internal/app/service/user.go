package service

import (
	"app/internal/app/repository"
	"app/internal/dto/request"
	resp "app/internal/dto/response"
	"app/pkg/config"
	rds "app/pkg/database/redis"
	"app/pkg/middleware"
	"app/pkg/toolkit"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SUser interface {
	Login(ctx context.Context, req request.Login) (*resp.RespLogin, error)
}

type sUser struct {
	rUser repository.RUser
	cfg   *config.Config
	rds   rds.Redis
}

func NewSUser(rUser repository.RUser, cfg *config.Config, rds rds.Redis) SUser {
	return &sUser{
		rUser: rUser,
		cfg:   cfg,
		rds:   rds,
	}
}

/*-------------------------- Main Function --------------------------*/

func (s *sUser) Login(ctx context.Context, req request.Login) (*resp.RespLogin, error) {

	// check lock login
	lockKey := fmt.Sprintf("login:lock:%s", req.Username)
	failKey := fmt.Sprintf("login:fail:%s", req.Username)

	// check lock
	err := s.checkLock(ctx, lockKey, req.Username)
	if err != nil {
		return nil, err
	}

	// find user by username or email
	user, err := s.rUser.FindByUsernameOrEmail(ctx, req.Username, req.Username)
	if err != nil {
		return nil, err
	}

	// check password
	if !toolkit.CheckPassword(user.Password, req.Password) {
		// lock user
		err := s.checkFail(ctx, failKey, lockKey, req.Username)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("invalid password or username")
	}

	// generate jti
	jti := uuid.New().String()

	// generate token
	token, err := middleware.GenerateToken(user.ID, jti, s.cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	// create session
	session := resp.SessionRecord{
		ID:       user.ID,
		UUID:     user.UUID,
		Username: user.Username,
		Email:    user.Email,
	}

	// save token
	err = s.saveToken(ctx, jti, &session)
	if err != nil {
		return nil, err
	}

	// return response
	return &resp.RespLogin{
		Token: token,
	}, nil
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

	// check fail login
	if value > 3 {

		// set lock user 15 min
		err = s.rds.Expire(ctx, lockKey, time.Minute*5)
		if err != nil {
			return err
		}

		err = s.rds.Delete(ctx, failKey)
		if err != nil {
			return err
		}

		return fmt.Errorf("user %s is locked", username)
	}

	return nil
}

func (s *sUser) saveToken(ctx context.Context, jti string, session *resp.SessionRecord) error {
	// Marshal user to json
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("session:%s", jti)

	// save token to redis
	err = s.rds.SetWithDuration(ctx, key, string(sessionJson), time.Minute*10)
	if err != nil {
		return err
	}

	return nil
}
