package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/Unknwon/goconfig"
	"golang.org/x/term"
)

func setupProfile(profile string, key []byte) {
	fmt.Print("Enter Zabbix URL: ")
	var zbxHost string
	fmt.Scanln(&zbxHost)

	fmt.Print("Enter Zabbix username: ")
	var zbxUsername string
	fmt.Scanln(&zbxUsername)

	fmt.Print("Enter Zabbix password: ")
	zbxPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()

	fmt.Print("Confirm Zabbix password: ")
	zbxPasswordConfirm, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()

	if string(zbxPassword) != string(zbxPasswordConfirm) {
		fmt.Println("Passwords do not match.")
		os.Exit(1)
	}

	// Encrypt password
	ciphertext, err := encrypt(key, zbxPassword)
	if err != nil {
		fmt.Println("Failed to encrypt password:", err)
		os.Exit(1)
	}

	// Save to INI file
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg, err = goconfig.LoadFromData([]byte{})
			if err != nil {
				fmt.Println("Failed to create new config file:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Failed to load config file:", err)
			os.Exit(1)
		}
	}

	cfg.SetValue(profile, "zbxHost", zbxHost)
	cfg.SetValue(profile, "zbxUsername", zbxUsername)
	cfg.SetValue(profile, "zbxPassword", base64.StdEncoding.EncodeToString(ciphertext))
	cfg.SetValue(profile, "ExportType", "snapshot")
	goconfig.SaveConfigFile(cfg, "config.ini")
	fmt.Println("Saved profile to config.ini")
}

func runProfile(profile string, key []byte, outputFile string, timeFrom string, timeTo string) {
	// Load from INI file
	cfg, _ := goconfig.LoadConfigFile("config.ini")
	zbxHost, _ := cfg.GetValue(profile, "zbxHost")
	zbxUsername, _ := cfg.GetValue(profile, "zbxUsername")
	zbxPasswordBase64, _ := cfg.GetValue(profile, "zbxPassword")

	// Decrypt password
	zbxPasswordCiphertext, _ := base64.StdEncoding.DecodeString(zbxPasswordBase64)
	zbxPassword, err := decrypt(key, zbxPasswordCiphertext)
	if err != nil {
		fmt.Println("Failed to decrypt password:", err)
		os.Exit(1)
	}

	runZabbixExport(zbxHost, zbxUsername, string(zbxPassword), outputFile, timeFrom, timeTo, profile, cfg)
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Add padding
	padding := aes.BlockSize - len(text)%aes.BlockSize
	padtext := append(text, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, aes.BlockSize+len(padtext))
	iv := ciphertext[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padtext)

	return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	padding := ciphertext[len(ciphertext)-1]
	return ciphertext[:len(ciphertext)-int(padding)], nil
}
