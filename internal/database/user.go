package database

import (
	"github.com/energieip/sol200-authentication-go/internal/core"
	"github.com/energieip/sol200-authentication-go/internal/tools"
)

//SaveUser register user in database
func SaveUser(db Database, usr core.User) error {
	criteria := make(map[string]interface{})
	criteria["Username"] = usr.Username
	if usr.Password != nil {
		//offuscate password
		hashPass, _ := tools.HashAndSalt(*usr.Password)
		usr.Password = &hashPass
	}
	return SaveOnUpdateObject(db, usr, ConfigDB, UsersTable, criteria)
}

//GetUser retrive user from the database
func GetUser(db Database, username string) *core.User {
	criteria := make(map[string]interface{})
	criteria["Username"] = username
	stored, err := db.GetRecord(ConfigDB, UsersTable, criteria)
	if err != nil || stored == nil {
		return nil
	}
	user, err := core.ToUser(stored)
	if err != nil {
		return nil
	}
	return user
}

func RemoveUser(db Database, usr core.User) error {
	criteria := make(map[string]interface{})
	criteria["Username"] = usr.Username
	return db.DeleteRecord(ConfigDB, UsersTable, criteria)
}
