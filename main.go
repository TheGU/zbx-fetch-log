package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"flag"

	"pattapongj/zbx-fetch-log/v6/zabbix"

	"github.com/k0kubun/pp/v3"
)

func main() {

	zbxFilename := flag.String("filename", "text.txt", "The name of the output file")
	zbxHost := flag.String("api", "http://localhost:9080/api_jsonrpc.php", "The Zabbix host URL")
	zbxUsername := flag.String("username", "Admin", "The Zabbix username")
	zbxPassword := flag.String("password", "zabbix", "The Zabbix password")
	timeFrom := flag.String("timeFrom", "1m", "The relative time from now")
	flag.Parse()

	// Time variable to set time from and time to in history.get
	zbxTimeTill := time.Now()
	zbxTimeFrom, err := relativeToAbsoluteTime(*timeFrom)
	if err != nil {
		panic(err)
	}
	pp.Println("time from", zbxTimeFrom)
	pp.Println("time till", zbxTimeTill)
	// zbxTimeTillFormatted := zbxTimeTill.Format("2006-01-02 15:04:05")

	// tzloc, err := time.LoadLocation("Asia/Bangkok")
	// if err != nil {
	// 	fmt.Println("Error loading location:", err)
	// 	return
	// }

	// Default approach - without session caching
	session, err := zabbix.NewSession(*zbxHost, *zbxUsername, *zbxPassword)
	if err != nil {
		panic(err)
	}

	version, err := session.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Zabbix API v%s : %s\n", version, *zbxHost)

	hosts, err := session.GetHosts(zabbix.HostGetParams{})
	if err != nil {
		panic(err)
	}
	// pp.Println(hosts)

	hostmap := make(map[string]zabbix.Host)
	for _, host := range hosts {
		hostmap[host.HostID] = host
	}
	// pp.Println(hostmap)

	// pp.Println(zabbix.ItemGetParams{
	// 	GetParameters: zabbix.GetParameters{OutputFields: zabbix.SelectExtendedOutput},
	// })

	items, err := session.GetItems(zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{OutputFields: zabbix.SelectExtendedOutput},
	})
	if err != nil {
		panic(err)
	}
	// pp.Println(items[0])

	itemmap := make(map[int]zabbix.Item)
	for _, item := range items {
		itemmap[item.ItemID] = item
	}
	// pp.Println(itemmap)

	zbxTimeFromFloat := float64(zbxTimeFrom.Unix())
	zbxTimeTillFloat := float64(zbxTimeTill.Unix())

	histories, err := session.GetHistories(zabbix.HistoryGetParams{
		TimeFrom: zbxTimeFromFloat,
		TimeTill: zbxTimeTillFloat,
	})
	if err != nil {
		panic(err)
	}
	// fmt.Println(time.Unix(histories[0].Clock, 0).Format(time.RFC3339))
	// pp.Println(histories[0])

	// Print history to file with host and item lookup
	f, err := os.OpenFile(*zbxFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// for _, history := range histories {
	// 	_hostname := hostmap[strconv.Itoa(itemmap[history.ItemID].HostID)].Hostname
	// 	text := fmt.Sprintf(
	// 		"%s, HostID=\"%d\", Host=\"%s\", ItemID=\"%d\", Item=\"%s\", Value=\"%s\"\n",
	// 		time.Unix(history.Clock, 0).Format(time.RFC3339),
	// 		itemmap[history.ItemID].HostID,
	// 		_hostname,
	// 		history.ItemID,
	// 		IReplace(itemmap[history.ItemID].ItemName, _hostname, ""),
	// 		history.Value)

	// 	if _, err = f.WriteString(text); err != nil {
	// 		panic(err)
	// 	}
	// }

	// fmt.Println(time.Unix(histories[0].Clock, 0).In(tzloc).Format(time.RFC3339))

	// Sort histories by Clock
	sort.Slice(histories, func(i, j int) bool {
		return histories[i].Clock < histories[j].Clock
	})

	for _, history := range histories {
		_hostname := hostmap[strconv.Itoa(itemmap[history.ItemID].HostID)].Hostname
		// text := fmt.Sprintf(
		// 	"%s, HostID=\"%d\", Host=\"%s\", ItemID=\"%d\", Item=\"%s\", Value=\"%s\"\n",
		// 	time.Unix(history.Clock, 0).Format(time.RFC3339),
		// 	itemmap[history.ItemID].HostID,
		// 	_hostname,
		// 	history.ItemID,
		// 	IReplace(itemmap[history.ItemID].ItemName, _hostname, ""),
		// 	history.Value)
		text := fmt.Sprintf(
			"%s, HostID=\"%d\", Host=\"%s\", ItemID=\"%d\", Key=\"%s\", Value=\"%s\"\n",
			time.Unix(history.Clock, 0).Format(time.RFC3339),
			itemmap[history.ItemID].HostID,
			_hostname,
			history.ItemID,
			itemmap[history.ItemID].Key,
			history.Value)

		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
	}

}
