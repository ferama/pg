package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ferama/pg/cmd"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")

	usr, _ := user.Current()
	homeConf := filepath.Join(usr.HomeDir, ".pg/")
	viper.AddConfigPath(".")
	viper.AddConfigPath(homeConf)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println(fmt.Errorf("fatal error config file: %w", err))
		os.Exit(1)
	}
}

func main() {
	cmd.Execute()
}
