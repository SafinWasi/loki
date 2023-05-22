package client

type Client struct {
	Client_id     string
	Client_secret string
}

type ClientDict map[string]Client

func (cd ClientDict) AddClient(client_id string, client_secret string, alias string) {
	var newClient Client
	newClient.Client_id = client_id
	newClient.Client_secret = client_secret
	cd[alias] = newClient
}

func (cd ClientDict) DeleteClient(alias string) {
	delete(cd, alias)
}
