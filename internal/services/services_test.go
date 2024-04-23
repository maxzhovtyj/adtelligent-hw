package services

import (
	"fmt"
	"github.com/maxzhovtyj/adtelligent-hw/internal/storage"
	"github.com/maxzhovtyj/adtelligent-hw/pkg/db/postgres"
	"testing"
)

var (
	testServices = initTestEnv()
)

func initTestEnv() Services {
	conn, err := postgres.NewConn()
	if err != nil {
		panic(err)
	}
	appStorage := storage.New(conn)

	return New(appStorage)
}

func TestServices_Generate(t *testing.T) {
	err := testServices.Generate()
	if err != nil {
		t.Error(err)
	}
}

func TestServices_GetUnlinkedCampaigns(t *testing.T) {
	campaigns, err := testServices.GetUnlinkedCampaigns()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(campaigns))
}

func TestServices_GetEntitiesNames(t *testing.T) {
	names, err := testServices.GetEntitiesNames()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(names))
}

func TestServices_MostDemandedSources(t *testing.T) {
	sources, err := testServices.MostDemandedSources(5)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(sources))
	for _, src := range sources {
		fmt.Println(src.ID, src.Count)
	}
}

func BenchmarkHandler_sourceCampaigns(b *testing.B) {
	req := Acquire()
	defer Release(req)
	req.ID = 500

	for i := 0; i < b.N; i++ {
		campaigns, err := testServices.GetSourceCampaigns(req)
		if err != nil {
			fmt.Println(err)
		}

		_ = campaigns
	}
}
