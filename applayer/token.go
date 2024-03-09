package applayer

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
	"time"
)

type IssueNewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type IssueNewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (app *app) IssueNewAccessToken(ctx *gin.Context, req IssueNewAccessTokenRequest) (rsp IssueNewAccessTokenResponse, err error) {
	// Verify refresh token
	refreshPayload, err := app.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil { // The refresh token is invalid or expired
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return rsp, err
	}

	// Find the refresh token in the "sessions" table in the database
	session, err := app.store.GetSession(ctx, refreshPayload.ID)
	if err != nil { // The refresh token is not found
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return rsp, err
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err)) // Other errors
		return rsp, err
	}

	// Check if the refresh token is blocked
	if session.IsBlocked {
		err = errors.New("session is blocked")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return rsp, err
	}

	// Ensure the session's username is the same as the refresh token's subject
	if session.Username != refreshPayload.Subject {
		err = errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return rsp, err
	}

	// Ensure the session's refresh token is the same as the request's refresh token
	if session.RefreshToken != req.RefreshToken {
		err = errors.New("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return rsp, err
	}

	// Create a new access token
	accessToken, accessPayload, err := app.tokenMaker.CreateToken(
		refreshPayload.Subject,
		app.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return rsp, err
	}

	rsp = IssueNewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}

	return rsp, nil
}
