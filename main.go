package main

import (
	"fmt"
	"os"

	"pattapongj/zbx-fetch-log/v5/zabbix"
)

func zabbixLogin(z *zabbix.Context, zbxHost, zbxUsername, zbxPassword string) {

	if err := z.Login(zbxHost, zbxUsername, zbxPassword); err != nil {
		fmt.Println("Login error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Login: success")
	}
}

func zabbixLogout(z *zabbix.Context) {

	if err := z.Logout(); err != nil {
		fmt.Println("Logout error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Logout: success")
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {

	var z zabbix.Context

	/* Get variables from environment to login to Zabbix server */
	zbxHost := getEnv("ZABBIX_HOST", "http://localhost:9080/api_jsonrpc.php")
	zbxUsername := getEnv("ZABBIX_USERNAME", "Admin")
	zbxPassword := getEnv("ZABBIX_PASSWORD", "zabbix")
	if zbxHost == "" || zbxUsername == "" || zbxPassword == "" {
		fmt.Println("Login error: make sure environment variables `ZABBIX_HOST`, `ZABBIX_USERNAME` and `ZABBIX_PASSWORD` are defined")
		os.Exit(1)
	}

	/* Login to Zabbix server */
	zabbixLogin(&z, zbxHost, zbxUsername, zbxPassword)
	defer zabbixLogout(&z)

	/* Get all hosts */
	hObjects, _, err := z.HostGet(zabbix.HostGetParams{
		GetParameters: zabbix.GetParameters{
			Output: zabbix.SelectExtendedOutput,
		},
	})
	if err != nil {
		fmt.Println("Hosts get error:", err)
		return
	}

	/* Print names of retrieved hosts */
	fmt.Println("Hosts list:")
	for _, h := range hObjects {
		fmt.Println("-", h.Host)
	}
}
