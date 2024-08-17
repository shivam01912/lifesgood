package model

type Credentials struct {
	Username string
	Cost     int `default:9`
	Password string
}
