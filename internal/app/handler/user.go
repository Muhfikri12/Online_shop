package handler

import (
	"app/internal/app/service"
	"app/internal/dto/request"
	"app/pkg/config"
	"app/pkg/toolkit"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HUser interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type hUser struct {
	sUser service.SUser
	cfg   *config.Config
}

func NewHUser(sUser service.SUser, cfg *config.Config) HUser {
	return &hUser{
		sUser: sUser,
		cfg:   cfg,
	}
}

const refreshCookieName = "refresh_token"

// Login godoc
// @Summary      Login
// @Description  Authenticate with username and password. Returns access token (JWT) and sets refresh token in HttpOnly cookie. Use remember_me for longer refresh token validity (30 days vs 7 days).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.Login  true  "Login credentials"
// @Success      200   {object}  toolkit.Response  "data contains access_token, expires_in"
// @Failure      400   {object}  toolkit.Response
// @Router       /login [post]
func (h *hUser) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		toolkit.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	respLogin, refreshToken, refreshTTL, err := h.sUser.Login(c.Request.Context(), req)
	if err != nil {
		toolkit.ResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	h.setRefreshCookie(c, refreshToken, refreshTTL)
	toolkit.ResponseOK(c, respLogin)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Issue a new access token using the refresh token from HttpOnly cookie (sent automatically by the browser). Performs token rotation (old refresh token is revoked).
// @Tags         auth
// @Produce      json
// @Success      200  {object}  toolkit.Response  "data contains new access_token, expires_in"
// @Failure      401  {object}  toolkit.Response
// @Router       /refresh-token [post]
func (h *hUser) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		toolkit.ResponseError(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	respLogin, newRefreshToken, refreshTTL, err := h.sUser.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		toolkit.ResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.setRefreshCookie(c, newRefreshToken, refreshTTL)
	toolkit.ResponseOK(c, respLogin)
}

// Logout godoc
// @Summary      Logout
// @Description  Invalidate current session: revoke refresh tokens and clear Redis session. Requires valid access token.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  toolkit.Response  "data contains message"
// @Failure      401  {object}  toolkit.Response
// @Failure      500  {object}  toolkit.Response
// @Router       /logout [post]
func (h *hUser) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie(refreshCookieName)
	jti := c.GetString("jti")
	userID := c.GetInt("user_id")

	if err := h.sUser.Logout(c.Request.Context(), refreshToken, jti, userID); err != nil {
		toolkit.ResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.clearRefreshCookie(c)
	toolkit.ResponseOK(c, gin.H{"message": "logged out"})
}

func (h *hUser) setRefreshCookie(c *gin.Context, token string, ttl time.Duration) {
	maxAge := int(ttl.Seconds())
	secure := h.cfg.Environment != "local" && h.cfg.Environment != "development"

	c.SetCookie(
		refreshCookieName,
		token,
		maxAge,
		"/",
		"",
		secure,
		true, // HttpOnly
	)
}

func (h *hUser) clearRefreshCookie(c *gin.Context) {
	secure := h.cfg.Environment != "local" && h.cfg.Environment != "development"

	c.SetCookie(
		refreshCookieName,
		"",
		-1,
		"/",
		"",
		secure,
		true,
	)
}

