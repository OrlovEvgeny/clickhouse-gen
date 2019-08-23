package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

const (
	COLOR_RESET = "\u001b[0m"
	COLOR_RED   = "\033[31m"
	COLOR_BLUE  = "\u001b[34m"
	COLOR_GREEN = "\u001b[32m"
	COLOR_CYAN  = "\u001b[36m"
)

var setting settingModel

var (
	SettingFile string
	Path        string
	Pack        string
	Table       string
	ClickHouse  = &setting.DBClickhouse
)

//settingModel
type settingModel struct {
	DBClickhouse struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		BaseName string `yaml:"basename"`
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"clickhouse"`
}

//buildConnect
func buildConnect() string {
	var (
		host     = ClickHouse.Host
		port     = ClickHouse.Port
		username = ClickHouse.UserName
		password = ClickHouse.Password
		database = ClickHouse.BaseName
	)

	return fmt.Sprintf("tcp://%s:%d?debug=%v&username=%s&password=%s&database=%s",
		host,
		port,
		false,
		username,
		password,
		database)
}

//Read and parse config file
func LoadSettings(configFile string) (error error) {
	error = nil
	filename, err := filepath.Abs(configFile)
	if err != nil {
		fmt.Println("Fail find settings")
		return err
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Fail open settings")
		return err
	}
	err = yaml.Unmarshal(yamlFile, &setting)
	if err != nil {
		log.Printf("[%s] error: %s parse from file %s\n", "clickhouse-gen", err, filename)
		return err
	}
	fmt.Printf("%vload settings √%v\n", COLOR_GREEN, COLOR_RESET)
	return error
}
