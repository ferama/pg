package conf

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

const (
	SqlTextareaHeight = 8
	ColorBlur         = "#ffffff"
	ColorFocus        = "#77bbee"
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

func GetDBConnURL(connName string) (string, error) {
	var conf Conf

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	for _, i := range conf.Connections {
		if i.Name == connName {
			return i.Url, nil
		}
	}
	return "", errors.New("conn string not found in config")
}
