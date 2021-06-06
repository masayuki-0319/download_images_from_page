package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// GetHTMLConfigList ...
type GetHTMLConfigList struct {
	ImgURLPattern	string
	StoreDir		string
}

// Config ...
var Config GetHTMLConfigList

func init() {
	cfg, err := ini.Load("../config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = GetHTMLConfigList{
		ImgURLPattern:	cfg.Section("base").Key("ImgURLPattern").String(),
		StoreDir:		cfg.Section("base").Key("StoreDir").String(),
	}
}
