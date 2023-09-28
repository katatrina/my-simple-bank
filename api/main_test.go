package api

import (
	db "github.com/katatrina/my-simple-bank/db/sqlc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/katatrina/my-simple-bank/util"
	"os"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
