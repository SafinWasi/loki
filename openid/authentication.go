package openid

func Authenticate(flow string, alias string) error {
	if flow == "code" {
		code()
	} else {
		token()
	}
	return nil
}

func code() error {
	return nil
}

func token() error {
	return nil
}
