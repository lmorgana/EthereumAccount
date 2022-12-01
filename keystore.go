package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type myKeystore struct {
	ks      *keystore.KeyStore
	ks_path string
}

func (it *myKeystore) Init(ks_path string) {
	it.ks_path = ks_path
	it.ks = keystore.NewKeyStore(ks_path, keystore.StandardScryptN, keystore.StandardScryptP)
}

func (it *myKeystore) isAccountExist() bool {
	existAccount := it.ks.Accounts()
	return len(existAccount) > 0
}

//get first account
func (it *myKeystore) getAccount() *accounts.Account {
	existAcccount := it.ks.Accounts()
	if len(existAcccount) > 0 {
		return &existAcccount[0]
	}
	return nil
}

//create new account and lock it's private key until user login
func (it *myKeystore) createAccount(password string) (*accounts.Account, error) {
	acc, err := it.ks.NewAccount(password)
	if err != nil {
		return nil, err
	}
	it.ks.Lock(acc.Address)
	return &acc, nil
}

//unlock user's private key by passphrase
func (it *myKeystore) loginAccount(acc accounts.Account, password string) error {
	err := it.ks.Unlock(acc, password)
	return err
}

func (it *myKeystore) testUnlock() {
	acc := it.getAccount()
	byte, err := it.ks.SignHash(*acc, []byte("JftTmWtFb8fUvr6bR4xneJYaLcynX52s"))
	if err != nil {
		fmt.Println("somebad", err.Error())
	}
	fmt.Println("we can sign ", byte)
}
