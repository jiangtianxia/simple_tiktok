package store

import "github.com/olivere/elastic/v7"

var (
	client *elastic.Client
)

func GetES() *elastic.Client {
	return client
}

func InitES(conn string, user string, password string) (err error) {
	client, err = elastic.NewClient(
		elastic.SetURL(conn),
		elastic.SetBasicAuth(user, password))
	return
}
