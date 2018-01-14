package main

import (
	"fmt"

	"github.com/go-ini/ini"
)

type SerestiService struct {
	name              string
	group             string
	enabled           bool
	version           int
	method            string
	sh                string
	input_parameters  []string
	output_parameters []string
}

type SerestiGlobal struct {
	listen            string
	server_name       string
	serve_files       bool
	serve_files_path  string
	serve_files_uri   string
	serve_api_version int
	verbose_error     bool
	debug             bool
	index_path        string
}

type SerestiConfig struct {
	global   SerestiGlobal
	services []SerestiService
}

func GetKey(cfg *ini.File, section string, key string) *ini.Key {
	v, err := cfg.Section(section).GetKey(key)
	if err != nil {
		k, _ := cfg.Section(section).NewKey(key, CONFIG_DEFAULT[key].(string))
		return k
	}
	return v
}

func GetKey_S(section *ini.Section, key string) *ini.Key {
	v, err := section.GetKey(key)
	if err != nil {
		panic(fmt.Errorf("[SERVICE ERROR]: %v", err))
	}
	return v
}

func ParseConfig(path string) (*SerestiConfig, error) {
	config := new(SerestiConfig)
	cfg, err := ini.Load(path)
	if err != nil {
		return config, err
	} else {
		config.global.listen = GetKey(cfg, "global", "listen").String()
		config.global.server_name = GetKey(cfg, "global", "server_name").String()
		config.global.serve_files = GetKey(cfg, "global", "serve_files").MustBool()
		config.global.serve_files_path = GetKey(cfg, "global", "serve_files_path").String()
		config.global.serve_files_uri = GetKey(cfg, "global", "serve_files_uri").String()
		config.global.serve_api_version = GetKey(cfg, "global", "serve_api_version").MustInt()
		config.global.verbose_error = GetKey(cfg, "global", "verbose_error").MustBool()
		config.global.debug = GetKey(cfg, "global", "debug").MustBool()
		config.global.index_path = GetKey(cfg, "global", "index_path").String()
		sections := cfg.Sections()
		for _, section := range sections {
			if section.Name() != "global" && section.Name() != "DEFAULT" {
				name := GetKey_S(section, "name").String()
				group := GetKey_S(section, "group").String()
				enabled := GetKey_S(section, "enabled").MustBool()
				version := GetKey_S(section, "version").MustInt()
				method := GetKey_S(section, "method").String()
				sh := GetKey_S(section, "sh").String()
				input_parameters := GetKey_S(section, "input_parameters").Strings("|")
				output_parameters := GetKey_S(section, "output_parameters").Strings("|")
				config.services = append(config.services, SerestiService{
					name:              name,
					group:             group,
					enabled:           enabled,
					version:           version,
					method:            method,
					sh:                sh,
					input_parameters:  input_parameters,
					output_parameters: output_parameters,
				})
			}
		}
	}
	return config, nil
}
