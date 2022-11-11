package conf

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const ConfDir = ".pg"

const (
	SqlTextareaHeight = 8
	ColorBlur         = "#ffffff"
	ColorFocus        = "#77bbee"

	ColorHeader        = "#33aa66"
	ColorTitle         = "#44bb77"
	ColorTableRowFocus = "#33aa66"
	ItemMaxLen         = 40
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

	user := viper.GetString("user-override")
	password := viper.GetString("password-override")
	database := viper.GetString("database-override")

	for _, i := range conf.Connections {
		if i.Name == connName {
			if user == "" {
				user = i.User
			}
			if password == "" {
				password = i.Password
			}
			if database == "" {
				database = i.Database
			}
			url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
				user,
				password,
				i.Host,
				i.Port,
				database,
			)
			return url, nil
		}
	}
	return "", errors.New("conn string not found in config")
}
