package conf

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

func GetAvailableConnections() []string {
	var conf Conf

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	res := make([]string, 0)
	for _, i := range conf.Connections {
		res = append(res, i.Name)
	}
	return res
}

func GetURL(connString string) (string, error) {
	var conf Conf

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	for _, i := range conf.Connections {
		if i.Name == connString {
			return i.Url, nil
		}
	}
	return "", errors.New("conn string not found in config")
}
