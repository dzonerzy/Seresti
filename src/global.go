package main

import "log"

var VERSION string = "0.1.0"
var APP_NAME string = "seresti"
var LOGGER *log.Logger
var VAR_PREFIX string = "SERESTI_"
var SERVER_POSTFIX = "powered by Seresti"

func LOG(config *SerestiConfig, v ...interface{}) {
	if config.global.debug {
		LOGGER.Println(v...)
	}
}

func ERR(config *SerestiConfig, v ...interface{}) {
	if config.global.debug {
		LOGGER.Fatalln(v...)
	}
}

var CONFIG_DEFAULT = map[string]interface{}{
	"listen":            ":8080",
	"server_name":       "Go",
	"serve_files":       "no",
	"serve_files_path":  "",
	"serve_files_uri":   "",
	"serve_api_version": "1",
	"verbose_error":     "yes",
	"debug":             "yes",
	"index_path":        "",
}
