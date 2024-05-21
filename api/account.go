package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// createAccount creates a new account.
//
//		 	@Router  		/accounts [post]
//			@Summary		Create a new account
//			@Description	Create by json account
//	     @Security		bearerToken
//			@Tags			accounts
//			@Accept			json
//			@Produce		json
//			@Param			account	body		createAccountRequest	true "Account info"
//			@Success		201
//			@Failure		400
//			@Failure		404
//			@Failure		500
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	authorizedPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)

	arg := db.CreateAccountParams{
		Owner:    authorizedPayload.Subject,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, newErrorResponse(pqErr))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, newErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	authorizedPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)
	if account.Owner != authorizedPayload.Subject {
		err = errors.New("accounts doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, newErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newErrorResponse(err))
		return
	}

	authorizedPayload := ctx.MustGet(authorizationPayloadKey).(*jwt.RegisteredClaims)

	arg := db.ListAccountsByOwnerParams{
		Owner:  authorizedPayload.Subject,
		Limit:  req.PageSize,
		Offset: req.PageSize * (req.PageID - 1),
	}

	accounts, err := server.store.ListAccountsByOwner(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, newErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// TODO: update account balance and delete account api
