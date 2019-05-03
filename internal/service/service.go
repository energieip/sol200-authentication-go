package service

import (
	"os"

	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/api"
	"github.com/energieip/sol200-authentication-go/internal/database"
	"github.com/energieip/sol200-authentication-go/internal/network"
	"github.com/romana/rlog"
)

//CoreService content
type CoreService struct {
	server network.ServerNetwork //Remote server
	db     database.Database
	api    *api.API
}

//Initialize service
func (s *CoreService) Initialize(confFile string) error {
	clientID := "AuthAPI"
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

	serverNet, err := network.CreateServerNetwork()
	if err != nil {
		rlog.Error("Cannot connect to broker " + conf.NetworkBroker.IP + " error: " + err.Error())
		return err
	}
	s.server = *serverNet

	err = s.server.LocalConnection(*conf, clientID)
	if err != nil {
		rlog.Error("Cannot connect to drivers broker " + conf.NetworkBroker.IP + " error: " + err.Error())
		return err
	}
	web := api.InitAPI(s.db, *conf)
	s.api = web

	rlog.Info("Authentication service started")
	return nil
}

//Stop service
func (s *CoreService) Stop() {
	rlog.Info("Stopping Authentication service")
	s.server.Disconnect()
	s.db.Close()
	rlog.Info("Authentication service stopped")
}

//Run service mainloop
func (s *CoreService) Run() error {
	for {
		select {
		case apiEvents := <-s.api.EventsToBackend:
			for eventType, event := range apiEvents {
				rlog.Info("get API event", eventType, event)
				switch eventType {

				}
			}
		}
	}
}
