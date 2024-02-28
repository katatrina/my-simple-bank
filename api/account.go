package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/applayer"
	"github.com/katatrina/my-simple-bank/token"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
)

func (server *HTTPServer) createAccount(ctx *gin.Context) {
	// What happens if ctx.MustGet panics?
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req applayer.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	if req.Owner != authPayload.Subject {
		err := errors.New("account owner must match authenticated user")
		ctx.JSON(http.StatusUnprocessableEntity, util.ErrorResponse(err))
		return
	}

	account, err := server.app.CreateAccount(ctx, req)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *HTTPServer) getAccount(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req applayer.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := server.app.GetAccount(ctx, req)
	if err != nil {
		return
	}

	if account.Owner != authPayload.Subject {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *HTTPServer) listAccounts(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req applayer.ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	accounts, err := server.app.ListAccounts(ctx, req, authPayload.Subject)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
