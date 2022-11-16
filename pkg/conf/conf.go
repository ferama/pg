package conf

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

const ConfDir = ".pg"

const (
	SqlTextareaHeight = 12
	ColorBlur         = "#ffffff"
	ColorFocus        = "#77bbee"

	ColorHeader        = "#3a6"
	ColorError         = "#d66"
	ColorTitle         = "#44bb77"
	ColorTableRowFocus = "#33aa66"
	ItemMaxLen         = 50
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
