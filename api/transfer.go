package api

import (
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/applayer"
	"github.com/katatrina/my-simple-bank/token"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
)

func (server *HTTPServer) makeTransfer(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req applayer.MoneyTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	transfer, err := server.app.MakeMoneyTransfer(ctx, req, authPayload.Subject)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}
