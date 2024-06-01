package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"luidium.com/luidium-cli/cmd"
)

func main() {
	viper.SetConfigName("luidium")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		viper.Set("config", "NOT_CONFIGURED")
	} else {
		viper.Set("config", "CONFIGURED")
	}

	tokenViper := viper.New()
	tokenViper.SetConfigName("cli_token")
	tokenViper.SetConfigType("toml")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}
	configDir := filepath.Join(homeDir, ".luidium")
	tokenViper.AddConfigPath(configDir)
	err = tokenViper.ReadInConfig()
	if err != nil {
		viper.Set("cli_token", "NOT_AUTHENTICATED")
	} else {
		viper.Set("cli_token", tokenViper.Get("cli_token"))
	}

	cmd.Execute()
}
