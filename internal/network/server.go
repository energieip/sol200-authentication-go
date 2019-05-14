package network

import (
	"time"

	genericNetwork "github.com/energieip/common-components-go/pkg/network"
	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/romana/rlog"
)

//ServerNetwork network object
type ServerNetwork struct {
	Iface genericNetwork.NetworkInterface
}

//CreateServerNetwork create network server object
func CreateServerNetwork() (*ServerNetwork, error) {
	serverBroker, err := genericNetwork.NewNetwork(genericNetwork.MQTT)
	if err != nil {
		return nil, err
	}
	serverNet := ServerNetwork{
		Iface: serverBroker,
	}
	return &serverNet, nil

}

//LocalConnection connect service to server broker
func (net ServerNetwork) LocalConnection(conf pkg.ServiceConfig, clientID string) error {
	cbkServer := make(map[string]func(genericNetwork.Client, genericNetwork.Message))

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
