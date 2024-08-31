package config

type Flow int

const (
	HOME Flow = iota
	CREATE
	UPDATE
	DELETE
)
