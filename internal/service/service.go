package service

import (
	"os"

	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/api"
	"github.com/energieip/sol200-authentication-go/internal/database"
	"github.com/romana/rlog"
)

//CoreService content
type CoreService struct {
	db  database.Database
	api *api.API
}

//Initialize service
func (s *CoreService) Initialize(confFile string) error {
	conf, err := pkg.ReadServiceConfig(confFile)
	if err != nil {
		rlog.Error("Cannot parse configuration file " + err.Error())
		return err
	}
	os.Setenv("RLOG_LOG_LEVEL", conf.LogLevel)
	os.Setenv("RLOG_LOG_NOTIME", "yes")
	rlog.UpdateEnv()
	rlog.Info("Starting Authentication service")

	db, err := database.ConnectDatabase(conf.DB.ClientIP, conf.DB.ClientPort)
	if err != nil {
		rlog.Error("Cannot connect to database " + err.Error())
		return err
	}
	s.db = *db

	web := api.InitAPI(s.db, *conf)
	s.api = web

	rlog.Info("Authentication service started")
	return nil
}

//Stop service
func (s *CoreService) Stop() {
	rlog.Info("Stopping Authentication service")
	s.db.Close()
	rlog.Info("Authentication service stopped")
}

//Run service mainloop
func (s *CoreService) Run() error {
	for {
		select {}
	}
}
