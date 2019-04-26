package core

import "encoding/json"

//User
type User struct {
	Username    string   `json:"username"`
	Password    *string  `json:"password,omitempty"`
	Priviledges []string `json:"priviledges"`
	Teams       []string `json:"teams"`
	AccessGroup []int    `json:"accessGroups"`
}

// ToJSON dump User struct
func (u User) ToJSON() (string, error) {
	inrec, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(inrec[:]), err
}

//ToUser convert map interface to User object
func ToUser(val interface{}) (*User, error) {
	var m User
	inrec, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &m)
	return &m, err
}

//UserAuthorization
type UserAuthorization struct {
	Priviledges []string `json:"priviledges"`
	AccessGroup []int    `json:"accessGroups"`
}

// ToJSON dump User struct
func (u UserAuthorization) ToJSON() (string, error) {
	inrec, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(inrec[:]), err
}

//ToUserAuthorization convert map interface to UserAuthorization object
func ToUserAuthorization(val interface{}) (*UserAuthorization, error) {
	var m UserAuthorization
	inrec, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &m)
	return &m, err
}
