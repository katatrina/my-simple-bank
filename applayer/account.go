package applayer

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/util"
	"net/http"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (app *app) CreateAccount(ctx *gin.Context, req CreateAccountRequest) (db.Account, error) {
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := app.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return db.Account{}, err
	}

	return account, nil
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (app *app) GetAccount(ctx *gin.Context, req GetAccountRequest) (db.Account, error) {
	account, err := app.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, util.ErrorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		}

		return db.Account{}, err
	}

	return account, nil
}

type ListAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (app *app) ListAccounts(ctx *gin.Context, req ListAccountsRequest, owner string) ([]db.Account, error) {
	arg := db.ListAccountsParams{
		Owner:  owner,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := app.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return []db.Account{}, err
	}

	return accounts, nil
}
