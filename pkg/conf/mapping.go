package conf

type Item struct {
	Name string

	User     string
	Password string
	Host     string
	Port     int
	Database string
}
type Conf struct {
	Connections []Item
}
