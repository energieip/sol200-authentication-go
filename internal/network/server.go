package network

import (
	"encoding/json"
	"time"

	"github.com/energieip/common-components-go/pkg/dserver"
	genericNetwork "github.com/energieip/common-components-go/pkg/network"
	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/core"
	"github.com/romana/rlog"
)

//ServerNetwork network object
type ServerNetwork struct {
	Iface  genericNetwork.NetworkInterface
	Events chan map[string]dserver.ServerConfig
}

//CreateServerNetwork create network server object
func CreateServerNetwork() (*ServerNetwork, error) {
	serverBroker, err := genericNetwork.NewNetwork(genericNetwork.MQTT)
	if err != nil {
		return nil, err
	}
	serverNet := ServerNetwork{
		Iface:  serverBroker,
		Events: make(chan map[string]dserver.ServerConfig),
	}
	return &serverNet, nil

}

//LocalConnection connect service to server broker
func (net ServerNetwork) LocalConnection(conf pkg.ServiceConfig, clientID string) error {
	cbkServer := make(map[string]func(genericNetwork.Client, genericNetwork.Message))
	cbkServer["/read/server/+/setup/hello"] = net.onHello

	confServer := genericNetwork.NetworkConfig{
		IP:         conf.NetworkBroker.IP,
		Port:       conf.NetworkBroker.Port,
		ClientName: clientID,
		Callbacks:  cbkServer,
		LogLevel:   conf.LogLevel,
		User:       conf.NetworkBroker.Login,
		Password:   conf.NetworkBroker.Password,
		CaPath:     conf.NetworkBroker.CaPath,
	}

	for {
		rlog.Info("Try to connect to " + conf.NetworkBroker.IP)
		err := net.Iface.Initialize(confServer)
		if err == nil {
			rlog.Info(clientID + " connected to server broker " + conf.NetworkBroker.IP)
			return err
		}
		timer := time.NewTicker(time.Second)
		rlog.Error("Cannot connect to broker " + conf.NetworkBroker.IP + " error: " + err.Error())
		rlog.Error("Try to reconnect " + conf.NetworkBroker.IP + " in 1s")

		select {
		case <-timer.C:
			continue
		}
	}
}

//Disconnect from server
func (net ServerNetwork) Disconnect() {
	net.Iface.Disconnect()
}

func (net ServerNetwork) onHello(client genericNetwork.Client, msg genericNetwork.Message) {
	payload := msg.Payload()
	rlog.Info(msg.Topic() + " : " + string(payload))
	var serverConfig dserver.ServerConfig
	err := json.Unmarshal(payload, &serverConfig)
	if err != nil {
		rlog.Error("Cannot parse config ", err.Error())
		return
	}

	event := make(map[string]dserver.ServerConfig)
	event[core.HelloEvent] = serverConfig
	net.Events <- event
}

//SendData to server
func (net ServerNetwork) SendData(topic, content string) error {
	err := net.Iface.SendCommand(topic, content)
	if err != nil {
		rlog.Error("Cannot send : " + content + " on: " + topic + " Error: " + err.Error())
	} else {
		rlog.Info("Sent : " + content + " on: " + topic)
	}
	return err
}
