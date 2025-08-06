package storage

import (
	"encoding/json"
	"errors"

	"os"
	"sojebsikder/utils"
)

type Account struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type Store struct {
	Accounts []Account `json:"accounts"`
}

const storageFileName = "accounts.enc"

func Load(passphrase []byte) (*Store, error) {
	if _, err := os.Stat(storageFileName); os.IsNotExist(err) {
		return &Store{Accounts: []Account{}}, nil
	}

	data, err := os.ReadFile(storageFileName)
	if err != nil {
		return nil, err
	}

	plain, err := utils.DecryptJSON(data, passphrase)
	if err != nil {
		return nil, errors.New("failed to decrypt store: " + err.Error())
	}

	var st Store
	if err := json.Unmarshal(plain, &st); err != nil {
		return nil, err
	}
	return &st, nil
}

func Save(st *Store, passphrase []byte) error {
	plain, err := json.Marshal(st)
	if err != nil {
		return err
	}

	enc, err := utils.EncryptJSON(plain, passphrase)
	if err != nil {
		return err
	}

	tmp := storageFileName + ".tmp"
	if err := os.WriteFile(tmp, enc, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, storageFileName)
}
