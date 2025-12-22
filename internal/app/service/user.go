package service

import (
	"app/internal/app/repository"
	"app/internal/dto/request"
	resp "app/internal/dto/response"
	"app/pkg/config"
	"app/pkg/middleware"
	"app/pkg/toolkit"
	"context"
)

type SUser interface {
	Login(ctx context.Context, req request.Login) (*resp.RespLogin, error)
}

type sUser struct {
	rUser repository.RUser
	cfg   *config.Config
}

func NewSUser(rUser repository.RUser, cfg *config.Config) SUser {
	return &sUser{
		rUser: rUser,
		cfg:   cfg,
	}
}

func (s *sUser) Login(ctx context.Context, req request.Login) (*resp.RespLogin, error) {
	user, err := s.rUser.FindByUsernameOrEmail(ctx, req.Username, req.Username)
	if err != nil {
		return nil, err
	}

	if !toolkit.CheckPassword(user.Password, req.Password) {
		return nil, err
	}

	token, err := middleware.GenerateToken(user.ID, s.cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &resp.RespLogin{
		Token: token,
	}, nil
}
