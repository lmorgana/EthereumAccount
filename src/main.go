package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"os"
)

func isValidPass(pass string) bool {
	return len(pass) >= 8
}

//get new password from stdin
func askForPassword() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("There is not an wallet, create new one.")

	fmt.Println("Enter password: ")
	pass, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	for !isValidPass(pass) {
		fmt.Println("Password not valid: must contain at least 8 characters.\nTry again: ")
		pass, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
	}
	return pass, nil
}

func getHash(s1 string) string {
	h := sha256.New()
	h.Write([]byte(s1))
	return string(h.Sum(nil))
}

func makeNewWallet(keystore *myKeystore, storage *storageKey) error {
	password, err := askForPassword()
	if err != nil {
		return err
	}
	h_pass := getHash(password)
	acc, err := keystore.createAccount(h_pass + password)
	if err = storage.Store(h_pass); err != nil {
		return err
	}
	fmt.Printf("New password was created with public key - %s\n", acc.Address)
	return nil
}

func loginWallet(acc accounts.Account, keystore *myKeystore, storage *storageKey) error {
	h_true_pass, err := storage.Read()
	if err != nil {
		return err
	}
	fmt.Println("Have an acount please enter password:")

	reader := bufio.NewReader(os.Stdin)
	pass, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	h_pass := getHash(pass)
	i := 0
	for i = 0; (h_pass) != h_true_pass && i < 4; i++ {
		fmt.Printf("Wrong answer %d try out of 5, password incorect, please, try again:\n", i+1)
		pass, err = reader.ReadString('\n')
		h_pass = getHash(pass)
		if err != nil {
			return err
		}
	}
	if i == 4 {
		return errors.New("You have exceeded the number of attempts")
	} else {
		err = keystore.loginAccount(acc, h_pass+pass)
		if err != nil {
			return err
		} else {
			fmt.Println(acc.Address, "was unlocked for test")
		}
	}
	return nil
}

func main() {
	path_pr_keys := "./src/wallets_keys/"
	password_path := "./src/pass/key.txt"
	hardcored_key_path := "./src/hardcored_key"

	var storage storageKey
	var keyStore myKeystore

	keyStore.Init(path_pr_keys)
	err := storage.Init(password_path, hardcored_key_path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer storage.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if keyStore.isAccountExist() {
		acc := keyStore.getAccount()
		err = loginWallet(*acc, &keyStore, &storage)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if keyStore.testCanWeMakeSign() {
			fmt.Println("we can make sign")
		} else {
			fmt.Println("we can't make sign")
		}
	} else {
		err = makeNewWallet(&keyStore, &storage)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
