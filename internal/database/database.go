package database

import (
	"github.com/energieip/common-components-go/pkg/database"
	"github.com/energieip/sol200-authentication-go/internal/core"

	"github.com/romana/rlog"
)

const (
	ConfigDB = "config"

	UsersTable = "users"
	TeamsTable = "teams"
)

type databaseError struct {
	s string
}

func (e *databaseError) Error() string {
	return e.s
}

// NewError raise an error
func NewError(text string) error {
	return &databaseError{text}
}

type Database = database.DatabaseInterface

//ConnectDatabase plug datbase
func ConnectDatabase(ip, port string) (*Database, error) {
	db, err := database.NewDatabase(database.RETHINKDB)
	if err != nil {
		rlog.Error("database err " + err.Error())
		return nil, err
	}

	confDb := database.DatabaseConfig{
		IP:   ip,
		Port: port,
	}
	err = db.Initialize(confDb)
	if err != nil {
		rlog.Error("Cannot connect to database " + err.Error())
		return nil, err
	}

	for _, dbName := range []string{ConfigDB} {
		err = db.CreateDB(dbName)
		if err != nil {
			rlog.Warn("Create DB ", err.Error())
		}

		tableCfg := make(map[string]interface{})
		if dbName == ConfigDB {
			tableCfg[TeamsTable] = core.Team{}
			tableCfg[UsersTable] = core.User{}
		}
		for tableName, objs := range tableCfg {
			err = db.CreateTable(dbName, tableName, &objs)
			if err != nil {
				rlog.Warn("Create table ", err.Error())
			}
		}
	}

	return &db, nil
}

//GetObjectID return id
func GetObjectID(db Database, dbName, tbName string, criteria map[string]interface{}) string {
	stored, err := db.GetRecord(dbName, tbName, criteria)
	if err == nil && stored != nil {
		m := stored.(map[string]interface{})
		id, ok := m["id"]
		if ok {
			return id.(string)
		}
	}
	return ""
}

//SaveOnUpdateObject in database
func SaveOnUpdateObject(db Database, obj interface{}, dbName, tbName string, criteria map[string]interface{}) error {
	var err error
	dbID := GetObjectID(db, dbName, tbName, criteria)
	if dbID == "" {
		_, err = db.InsertRecord(dbName, tbName, obj)
	} else {
		err = db.UpdateRecord(dbName, tbName, dbID, obj)
	}
	return err
}
