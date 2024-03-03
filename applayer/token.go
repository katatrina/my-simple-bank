package applayer

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
	"time"
)

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (app *app) RenewAccessToken(ctx *gin.Context, req RenewAccessTokenRequest) (RenewAccessTokenResponse, error) {
	refreshPayload, err := app.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	session, err := app.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return RenewAccessTokenResponse{}, err
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	if session.IsBlocked {
		err = errors.New("session is blocked")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	if session.Username != refreshPayload.Subject {
		err = errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	if session.RefreshToken != req.RefreshToken {
		err = errors.New("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	if time.Now().After(session.ExpiresAt) {
		err = errors.New("expired session")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	accessToken, accessPayload, err := app.tokenMaker.CreateToken(
		refreshPayload.Subject,
		app.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return RenewAccessTokenResponse{}, err
	}

	rsp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}

	return rsp, nil
}
