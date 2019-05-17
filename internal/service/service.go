package service

import (
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/energieip/sol200-authentication-go/internal/core"

	"github.com/energieip/common-components-go/pkg/duser"
	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/api"
	"github.com/energieip/sol200-authentication-go/internal/database"
	"github.com/energieip/sol200-authentication-go/internal/network"
	"github.com/romana/rlog"
)

//CoreService content
type CoreService struct {
	server      network.ServerNetwork //Remote server
	db          database.Database
	api         *api.API
	internalApi *api.InternalAPI
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

	internal := api.InitInternalAPI(s.db, *conf)
	s.internalApi = internal

	rlog.Info("Authentication service started")
	return nil
}

func (s *CoreService) createUser(evt interface{}) {
	user, err := core.ToUser(evt)
	if err != nil || user == nil {
		rlog.Error("could not parse event")
		return
	}

	err = database.SaveUser(s.db, *user)
	if err == nil {
		rlog.Info("User " + user.Username + " successfully added")
		token := user.Username + *user.Password
		hasher := sha256.New()
		hasher.Write([]byte(token))
		access := duser.UserAccess{
			UserHash:     hex.EncodeToString(hasher.Sum(nil)),
			Priviledge:   user.Priviledge,
			AccessGroups: user.AccessGroups,
			Services:     user.Services,
		}
		dump, _ := access.ToJSON()
		s.server.SendData(core.CreateUserEvent, dump)

	} else {
		rlog.Error("Error during database register: " + err.Error())
	}
}

func (s *CoreService) removeUser(evt interface{}) {
	user, err := core.ToUser(evt)
	if err != nil || user == nil {
		rlog.Error("could not parse event")
		return
	}

	if user.Password == nil {
		rlog.Error("could not find password")
		return
	}

	err = database.RemoveUser(s.db, *user)
	if err == nil {
		rlog.Info("User " + user.Username + " successfully removed")
		token := user.Username + *user.Password
		hasher := sha256.New()
		hasher.Write([]byte(token))
		access := duser.UserAccess{
			UserHash: hex.EncodeToString(hasher.Sum(nil)),
		}
		dump, _ := access.ToJSON()
		s.server.SendData(core.RemoveUserEvent, dump)

	} else {
		rlog.Error("Error during database register: " + err.Error())
	}
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

		case internalApiEvents := <-s.internalApi.EventsToBackend:
			for eventType, event := range internalApiEvents {
				rlog.Info("get internal API event", eventType, event)
				switch eventType {
				case core.CreateUserEvent:
					s.createUser(event)
				case core.RemoveUserEvent:
					s.removeUser(event)

				}
			}
		}
	}
}
