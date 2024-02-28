package api

import (
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/applayer"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
)

func (server *HTTPServer) createUser(ctx *gin.Context) {
	var req applayer.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	user, err := server.app.CreateUser(ctx, req)
	if err != nil {
		return
	}

	rsp := applayer.NewUserResponse(user)

	ctx.JSON(http.StatusOK, rsp)
}

func (server *HTTPServer) loginUser(ctx *gin.Context) {
	var req applayer.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	user, err := server.app.LoginUser(ctx, req)
	if err != nil {
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	rsp := applayer.LoginUserResponse{
		AccessToken: accessToken,
		User:        applayer.NewUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
