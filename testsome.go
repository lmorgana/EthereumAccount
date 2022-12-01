package main

import "fmt"

func main() {
	password_path := "./pass/key.txt"
	hardcored_key_path := "./hardcored_key"

	var storage storageKey

	err := storage.Init(password_path, hardcored_key_path)
	if err != nil {
		return
	}
	defer storage.Close()

	//err = storage.Store("myPassword")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	line, err := storage.Read()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(line)
}
