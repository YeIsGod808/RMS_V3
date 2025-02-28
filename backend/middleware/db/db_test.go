package db

import (
	"testing"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
)

func TestDbConn(t *testing.T) {
	if err := config.Init(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	log.Info("log init success...")

	client := GetDB()
	t.Log("client:", client)
	if client == nil {
		t.Error("client is nil")
	}
}
