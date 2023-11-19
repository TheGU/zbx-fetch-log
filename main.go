package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"pattapongj/zbx-fetch-log/v5/zabbix"

	"github.com/k0kubun/pp/v3"
)

func main() {

	FILENAME := "text.txt"

	zbxHost := getEnv("ZBX_URL", "http://localhost:9080/api_jsonrpc.php")
	zbxUsername := getEnv("ZBX_USERNAME", "Admin")
	zbxPassword := getEnv("ZBX_PASSWORD", "zabbix")

	// tzloc, err := time.LoadLocation("Asia/Bangkok")
	// if err != nil {
	// 	fmt.Println("Error loading location:", err)
	// 	return
	// }

	// Default approach - without session caching
	session, err := zabbix.NewSession(zbxHost, zbxUsername, zbxPassword)
	if err != nil {
		panic(err)
	}

	version, err := session.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Zabbix API v%s \n", version)

	hosts, err := session.GetHosts(zabbix.HostGetParams{})
	if err != nil {
		panic(err)
	}
	pp.Println(hosts[0])
	hostmap := make(map[string]zabbix.Host)
	for _, host := range hosts {
		hostmap[host.HostID] = host
	}
	// pp.Println(hostmap)

	items, err := session.GetItems(zabbix.ItemGetParams{})
	if err != nil {
		panic(err)
	}
	pp.Println(items[0])
	itemmap := make(map[int]zabbix.Item)
	for _, item := range items {
		itemmap[item.ItemID] = item
	}
	// pp.Println(itemmap)

	histories, err := session.GetHistories(zabbix.HistoryGetParams{})
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Unix(histories[0].Clock, 0).Format(time.RFC3339))
	pp.Println(histories[0])

	// Print history to file with host and item lookup
	f, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, history := range histories {
		_hostname := hostmap[strconv.Itoa(itemmap[history.ItemID].HostID)].Hostname
		text := fmt.Sprintf(
			"%s, HostID=\"%d\", Host=\"%s\", ItemID=\"%d\", Item=\"%s\", Value=\"%s\"\n",
			time.Unix(history.Clock, 0).Format(time.RFC3339),
			itemmap[history.ItemID].HostID,
			_hostname,
			history.ItemID,
			IReplace(itemmap[history.ItemID].ItemName, _hostname, ""),
			history.Value)
		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
	}

	// fmt.Println(time.Unix(histories[0].Clock, 0).In(tzloc).Format(time.RFC3339))

}
