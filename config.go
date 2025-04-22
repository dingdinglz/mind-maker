package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Model ModelConfig `json:"model"`
	Mcps  []McpConfig `json:"mcps"`
}

type ModelConfig struct {
	Model   string `json:"model"`
	ApiKey  string `json:"apikey"`
	BaseURL string `json:"baseurl"`
}

type McpConfig struct {
	Command string   `json:"command"`
	Env     []string `json:"env"`
	Args    []string `json:"args"`
}

var ActivateConfig Config

func LoadConfig() {
	if !FileExist("config.json") {
		panic("未发现配置文件")
	}
	res, _ := os.ReadFile("config.json")
	json.Unmarshal(res, &ActivateConfig)
}
