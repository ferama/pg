package main

import (
	"fmt"

	"github.com/ferama/gopigi/cmd"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("conf.yaml")
	viper.SetConfigType("yaml")
	// viper.AddConfigPath("$HOME/.gopigi")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	cmd.Execute()
}
