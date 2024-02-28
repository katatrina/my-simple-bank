package applayer

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
)

var (
	InvalidCurrencyError = errors.New("invalid currency")
)

type MoneyTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (app *app) MakeMoneyTransfer(ctx *gin.Context, req MoneyTransferRequest, authenticatedOwner string) (db.TransferTxResult, error) {
	// Check if the "from account" exists and the currency is valid
	fromAccount, valid := app.checkValidCurrency(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return db.TransferTxResult{}, InvalidCurrencyError
	}

	if fromAccount.Owner != authenticatedOwner {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return db.TransferTxResult{}, err
	}

	_, valid = app.checkValidCurrency(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return db.TransferTxResult{}, InvalidCurrencyError
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := app.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return db.TransferTxResult{}, err
	}

	return result, nil
}

func (app *app) checkValidCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := app.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account with ID [%d] has currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return account, false
	}

	return account, true
}
