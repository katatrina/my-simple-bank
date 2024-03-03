package api

import (
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/applayer"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
)

func (server *HTTPServer) renewAccessToken(ctx *gin.Context) {
	var req applayer.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	rsp, err := server.app.RenewAccessToken(ctx, req)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
