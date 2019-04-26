package core

import "encoding/json"

type Team struct {
	Name        string `json:"name"` //team name
	AccessGroup []int  `json:"accessGroups"`
}

// ToJSON dump User struct
func (u Team) ToJSON() (string, error) {
	inrec, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(inrec[:]), err
}

//ToTeam convert map interface to Team object
func ToTeam(val interface{}) (*Team, error) {
	var m Team
	inrec, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(inrec, &m)
	return &m, err
}
