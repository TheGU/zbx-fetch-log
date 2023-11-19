package main

import (
	"fmt"

	"pattapongj/zbx-fetch-log/v5/zabbix"
)

func main() {

	zbxHost := getEnv("ZBX_HOST", "http://localhost:9080/api_jsonrpc.php")
	zbxUsername := getEnv("ZBX_USERNAME", "Admin")
	zbxPassword := getEnv("ZBX_PASSWORD", "zabbix")

	// Default approach - without session caching
	session, err := zabbix.NewSession(zbxHost, zbxUsername, zbxPassword)
	if err != nil {
		panic(err)
	}

	version, err := session.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Zabbix API v%s", version)
}
