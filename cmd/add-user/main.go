package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/energieip/sol200-authentication-go/internal/core"

	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/romana/rlog"
)

type arrayString []string

func (i *arrayString) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayString) Set(value string) error {
	if value == "" {
		return nil
	}

	var list []string
	for _, in := range strings.Split(value, ",") {
		list = append(list, in)
	}

	*i = arrayString(list)
	return nil
}

func (i *arrayString) Get() interface{} { return []string(*i) }

type arrayInt []int

func (i *arrayInt) Set(val string) error {
	if val == "" {
		return nil
	}

	var list []int
	for _, in := range strings.Split(val, ",") {
		i, err := strconv.Atoi(in)
		if err != nil {
			return err
		}

		list = append(list, i)
	}

	*i = arrayInt(list)
	return nil
}

func (i *arrayInt) Get() interface{} { return []int(*i) }

func (i *arrayInt) String() string {
	var list []string
	for _, in := range *i {
		list = append(list, strconv.Itoa(in))
	}
	return strings.Join(list, ",")
}

func main() {
	var confFile string
	var username string
	var password string
	var priviledge string
	var teams arrayString
	var services arrayString
	var groups arrayInt

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.StringVar(&confFile, "config", "", "Specify an alternate configuration file.")
	flag.StringVar(&confFile, "c", "", "Specify an alternate configuration file.")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&username, "u", "", "username.")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&password, "p", "", "password.")
	flag.StringVar(&priviledge, "priviledge", "user", "Priviledge")
	flag.StringVar(&priviledge, "a", "user", "Priviledge")
	flag.Var(&groups, "groups", "Groups list comma separated")
	flag.Var(&groups, "g", "Groups list comma separated")
	flag.Var(&teams, "teams", "Teams list comma separated")
	flag.Var(&teams, "t", "Teams list comma separated")
	flag.Var(&services, "services", "Services list comma separated")
	flag.Var(&services, "s", "Services list comma separated")
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
		Username:     username,
		Password:     &password,
		Priviledge:   priviledge,
		Teams:        teams,
		AccessGroups: groups,
		Services:     services,
	}

	requestBody, err := json.Marshal(user)
	if err != nil {
		rlog.Error(err.Error())
		os.Exit(1)
	}

	url := "https://" + conf.InternalAPI.IP + ":" + conf.InternalAPI.Port + "/user"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
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

	rlog.Info("User " + username + " successfully added " + string(body))
}
