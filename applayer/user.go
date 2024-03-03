package applayer

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (app *app) CreateUser(ctx *gin.Context, req CreateUserRequest) (CreateUserResponse, error) {
	passwordHashed, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return CreateUserResponse{}, err
	}

	arg := db.CreatUserParams{
		Username:       req.Username,
		HashedPassword: passwordHashed,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := app.store.CreatUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return CreateUserResponse{}, err
	}

	rsp := NewUserResponse(user)

	return rsp, nil
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID          `json:"session_id"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  CreateUserResponse `json:"user"`
}

func (app *app) LoginUser(ctx *gin.Context, req LoginUserRequest) (LoginUserResponse, error) {
	user, err := app.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		}

		return LoginUserResponse{}, err
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return LoginUserResponse{}, err
	}

	accessToken, accessPayload, err := app.tokenMaker.CreateToken(
		user.Username,
		app.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return LoginUserResponse{}, err
	}

	refreshToken, refreshPayload, err := app.tokenMaker.CreateToken(
		user.Username,
		app.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return LoginUserResponse{}, err
	}

	session, err := app.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.GetHeader("User-Agent"),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return LoginUserResponse{}, err
	}

	rsp := LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt.Time,
		User:                  NewUserResponse(user),
	}

	return rsp, nil
}
