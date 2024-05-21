package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
)

type makeMoneyTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) makeMoneyTransfer(ctx *gin.Context) {
	var req makeMoneyTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	// Validate currency
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authorizedPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)
	if fromAccount.Owner != authorizedPayload.Subject {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferMoneyParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferMoneyTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, newErrorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account with ID %d has currency %s, not the expected currency %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, newErrorResponse(err))
		return account, false
	}

	return account, true
}
