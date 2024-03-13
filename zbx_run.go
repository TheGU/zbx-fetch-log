package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"pattapongj/zbx-fetch-log/v6/zabbix"
)

func runZabbixExport(zbxHost, zbxUsername, zbxPassword, outputFile, timeFrom string, timeTo string) {
	var zbxTimeTill time.Time
	var zbxTimeFrom time.Time
	var err error

	// Time variable to set time from and time to in history.get
	if timeTo != "" {
		zbxTimeTill, err = relativeToAbsoluteTime(timeTo)
		if err != nil {
			panic(err)
		}
	} else {
		zbxTimeTill = time.Now()
	}

	if timeFrom != "" {
		zbxTimeFrom, err = relativeToAbsoluteTime(timeFrom)
		if err != nil {
			panic(err)
		}
	} else {
		zbxTimeFrom, err = relativeToAbsoluteTime("5m")
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("time from", zbxTimeFrom)
	fmt.Println("time till", zbxTimeTill)

	// Default approach - without session caching
	session, err := zabbix.NewSession(zbxHost, zbxUsername, zbxPassword)
	if err != nil {
		panic(err)
	}

	version, err := session.GetVersion()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Zabbix API v%s : %s\n", version, zbxHost)

	hosts, err := session.GetHosts(zabbix.HostGetParams{})
	if err != nil {
		panic(err)
	}

	hostmap := make(map[string]zabbix.Host)
	for _, host := range hosts {
		hostmap[host.HostID] = host
	}

	items, err := session.GetItems(zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{OutputFields: zabbix.SelectExtendedOutput},
	})
	if err != nil {
		panic(err)
	}

	itemmap := make(map[int]zabbix.Item)
	for _, item := range items {
		itemmap[item.ItemID] = item
	}

	zbxTimeFromFloat := float64(zbxTimeFrom.Unix())
	zbxTimeTillFloat := float64(zbxTimeTill.Unix())

	histories, err := session.GetHistories(zabbix.HistoryGetParams{
		TimeFrom: zbxTimeFromFloat,
		TimeTill: zbxTimeTillFloat,
	})
	if err != nil {
		panic(err)
	}

	// Print history to file with host and item lookup
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Sort histories by Clock
	sort.Slice(histories, func(i, j int) bool {
		return histories[i].Clock < histories[j].Clock
	})

	for _, history := range histories {
		_hostname := hostmap[strconv.Itoa(itemmap[history.ItemID].HostID)].Hostname
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

	fmt.Printf("Exported %d histories to %s\n", len(histories), outputFile)
}
