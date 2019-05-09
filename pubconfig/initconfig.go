package pubconfig

type DbConfig struct {
	Listen string
	Redis  string
	Mysql  string
	Name   string
	Pass   string
	DbName string
}

type HallConfig struct {
	Listen   string
	RemoteDb string
	RemoteGM string
}

type GameConfig struct {
	Listen     string
	RemoteDb   string
	RemoteHall string
}

type LoginConfig struct {
	Listen     string
	RemoteHall string
}
