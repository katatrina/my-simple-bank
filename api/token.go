package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	// Verify the refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	sessionID, err := uuid.Parse(refreshPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	// Find the corresponding session of the refresh token in the database
	session, err := server.store.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, newErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	if session.IsBlocked {
		err = errors.New("blocked session")
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	if session.Username != refreshPayload.Subject {
		err = errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	if req.RefreshToken != session.RefreshToken {
		err = errors.New("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err = errors.New("expired session")
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Subject, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}
	ctx.JSON(http.StatusOK, rsp)
}
