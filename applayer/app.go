package applayer

import (
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
)

type App interface {
	CreateUser(ctx *gin.Context, req CreateUserRequest) (db.User, error)
	LoginUser(ctx *gin.Context, req LoginUserRequest) (db.User, error)
	CreateAccount(ctx *gin.Context, req CreateAccountRequest) (db.Account, error)
	GetAccount(ctx *gin.Context, req GetAccountRequest) (db.Account, error)
	ListAccounts(ctx *gin.Context, req ListAccountsRequest, owner string) ([]db.Account, error)
	MakeMoneyTransfer(ctx *gin.Context, req MoneyTransferRequest, authenticatedOwner string) (db.TransferTxResult, error)
}

type app struct {
	store db.Store
}

func New(store db.Store) App {
	return &app{
		store: store,
	}
}
