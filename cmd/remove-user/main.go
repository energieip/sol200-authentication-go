package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/energieip/sol200-authentication-go/internal/core"

	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/romana/rlog"
)

func main() {
	var confFile string
	var username string
	var password string

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.StringVar(&confFile, "config", "", "Specify an alternate configuration file.")
	flag.StringVar(&confFile, "c", "", "Specify an alternate configuration file.")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&username, "u", "", "username.")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&password, "p", "", "password.")
	flag.Parse()

	conf, err := pkg.ReadServiceConfig(confFile)
	if err != nil {
		rlog.Error("Cannot parse configuration file " + err.Error())
		os.Exit(1)
	}
	os.Setenv("RLOG_LOG_LEVEL", conf.LogLevel)
	os.Setenv("RLOG_LOG_NOTIME", "yes")
	rlog.UpdateEnv()

	user := core.User{
		Username: username,
		Password: &password,
	}

	requestBody, err := json.Marshal(user)
	if err != nil {
		rlog.Error(err.Error())
		os.Exit(1)
	}

	url := "https://" + conf.InternalAPI.IP + ":" + conf.InternalAPI.Port + "/user"
	req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(requestBody))
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	resp, err := client.Do(req)

	if err != nil {
		rlog.Error(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	rlog.Info("User " + username + " successfully removed " + string(body))
}
