package delivery

import (
	"fmt"
	"github.com/maxzhovtyj/adtelligent-hw/internal/services"
	"github.com/maxzhovtyj/adtelligent-hw/internal/storage"
	"github.com/maxzhovtyj/adtelligent-hw/pkg/db/postgres"
	"testing"
)

var (
	testHandler = initTestEnv()
)

func initTestEnv() *Handler {
	conn, err := postgres.NewConn()
	if err != nil {
		panic(err)
	}
	appStorage := storage.New(conn)
	appServices := services.New(appStorage)

	return New(appServices)
}

func BenchmarkHandler_sourceCampaigns(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req := acquire()
		defer release(req)

		req.ID = 500

		campaigns, err := testHandler.sourceCampaigns(req)
		if err != nil {
			fmt.Println(err)
		}

		_ = campaigns
	}
}
