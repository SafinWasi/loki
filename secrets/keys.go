package secrets

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/99designs/keyring"
)

var kr keyring.Keyring

func Initialize(debug bool) error {
	keyring.Debug = debug
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	credsDir := pwd + string(os.PathSeparator) + ".credentials"
	_, err = os.ReadDir(credsDir)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Credentials directory does not exist. Creating...")
		os.Mkdir(credsDir, os.ModePerm)
	}
	kr, err = keyring.Open(
		keyring.Config{
			AllowedBackends:  []keyring.BackendType{keyring.FileBackend},
			ServiceName:      "loki",
			FilePasswordFunc: keyring.FixedStringPrompt("abdcef"),
			FileDir:          credsDir})
	return err
}

func Set(key string, data []byte) error {
	return kr.Set(keyring.Item{Key: key, Data: data})
}

func Get(key string) ([]byte, error) {
	val, err := kr.Get(key)
	if err != nil {
		return nil, err
	}
	return val.Data, nil
}

func GetKeys() ([]string, error) {
	return kr.Keys()
}

func RemoveKey(key string) error {
	return kr.Remove(key)
}

func RemoveKeyring() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	credsDir := pwd + string(os.PathSeparator) + ".credentials"
	_, err = os.ReadDir(credsDir)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Credentials directory does not exist.")
		return nil
	}
	d, err := os.Open(credsDir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(credsDir, name))
		if err != nil {
			return err
		}
	}
	err = os.Remove(credsDir)
	return err
}
