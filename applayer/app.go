package applayer

import (
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/katatrina/my-simple-bank/token"
	"github.com/katatrina/my-simple-bank/util"
)

type App interface {
	CreateUser(ctx *gin.Context, req CreateUserRequest) (CreateUserResponse, error)
	LoginUser(ctx *gin.Context, req LoginUserRequest) (LoginUserResponse, error)
	CreateAccount(ctx *gin.Context, req CreateAccountRequest) (db.Account, error)
	GetAccount(ctx *gin.Context, req GetAccountRequest) (db.Account, error)
	ListAccounts(ctx *gin.Context, req ListAccountsRequest, owner string) ([]db.Account, error)
	MakeMoneyTransfer(ctx *gin.Context, req MoneyTransferRequest, authenticatedOwner string) (db.TransferTxResult, error)
	RenewAccessToken(ctx *gin.Context, req RenewAccessTokenRequest) (RenewAccessTokenResponse, error)
}

type app struct {
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func New(store db.Store, tokenMaker token.Maker, config util.Config) App {
	return &app{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}
