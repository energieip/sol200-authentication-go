package main

import (
	"flag"
	"os"

	"github.com/energieip/sol200-authentication-go/internal/core"

	pkg "github.com/energieip/common-components-go/pkg/service"
	"github.com/energieip/sol200-authentication-go/internal/database"
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

	db, err := database.ConnectDatabase(conf.DB.ClientIP, conf.DB.ClientPort)
	if err != nil {
		rlog.Error("Cannot connect to database " + err.Error())
		os.Exit(1)
	}

	user := core.User{
		Username: username,
		Password: &password,
	}
	err = database.SaveUser(*db, user)
	if err != nil {
		rlog.Error("Error when registering user" + username)
		os.Exit(1)
	}
	rlog.Info("User " + username + " successfully added")
}
