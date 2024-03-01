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

	userResponse, err := server.app.CreateUser(ctx, req)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (server *HTTPServer) loginUser(ctx *gin.Context) {
	var req applayer.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	loginUserResponse, err := server.app.LoginUser(ctx, req)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse)
}
