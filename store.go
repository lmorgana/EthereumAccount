package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

type storageKey struct {
	file *os.File
	key  []byte
}

func (it *storageKey) Init(password_path, key_path string) error {
	file, err := os.OpenFile(password_path, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	it.file = file
	key_file, err := os.Open(key_path)
	if err != nil {
		defer it.Close()
		return err
	}
	reader := bufio.NewReader(key_file)
	line, err := reader.ReadBytes(0)
	it.key = line
	return nil
}

func (it *storageKey) Close() {
	it.file.Close()
}

//by stored passkey read hash of account password
func (it *storageKey) Read() (string, error) {
	reader := bufio.NewReader(it.file)
	ciphertext, err := reader.ReadBytes(0)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}
	block, err := aes.NewCipher(it.key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

//by stored passkey write hash of account password
func (it *storageKey) Store(s1 string) error {
	block, err := aes.NewCipher(it.key)
	if err != nil {
		return err
	}
	plaintext := []byte(s1)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	it.file.Write(ciphertext)
	return nil
}
