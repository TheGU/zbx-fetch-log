package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"pattapongj/zbx-fetch-log/v6/zabbix"

	"github.com/Unknwon/goconfig"
)

func runZabbixExport(zbxHost, zbxUsername, zbxPassword, outputFile, timeFrom string, timeTo string, profile string, cfg *goconfig.ConfigFile) {
	var zbxTimeTill time.Time
	var zbxTimeFrom time.Time
	var err error

	// Get call limit from config in int format
	// callLimit, _ := cfg.Int(profile, "Limit")
	exportType, _ := cfg.GetValue(profile, "ExportType")

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

	// Get hosts ==========================

	paramHost := zabbix.HostGetParams{
		GetParameters: zabbix.GetParameters{
			OutputFields: []string{"hostid", "host", "name", "status", "available"},
			// ResultLimit: 1000,
			SortField: []string{"hostid"},
			Filter:    map[string]interface{}{"status": 0},
			SortOrder: "DESC",
		},
		SelectHostGroups: []string{"groupid", "name"},
	}

	hostCount, err := session.GetHostCount(paramHost)
	if err != nil {
		panic(err)
	}

	hosts, err := session.GetHosts(paramHost)
	if err != nil {
		panic(err)
	}
	if len(hosts) != hostCount {
		panic(fmt.Sprintf("Warning: Host count mismatch: %d != %d\n", len(hosts), hostCount))
	}

	hostmap := make(map[string]zabbix.Host)
	groupmap := make(map[string]string)
	hostIDs := make([]string, len(hosts))
	for i, host := range hosts {
		hostmap[host.HostID] = host
		hostIDs[i] = host.HostID

		groupmap[host.HostID] = ""
		for _, group := range host.Groups {
			groupmap[host.HostID] += group.Name + ","
		}
		// fmt.Printf("Host: %s\n", host.Hostname)
		// pp.Println(host)
	}

	// Get items ==========================

	paramItem := zabbix.ItemGetParams{
		GetParameters: zabbix.GetParameters{
			OutputFields: []string{"hostid", "itemid", "key_", "prevvalue", "lastvalue", "value_type", "lastclock"},
			// ResultLimit:  callLimit,
			Filter: map[string]interface{}{"status": 0},
			// SortField:  []string{"itemid"},
			// SortOrder:  "ASC",
			TextSearch: make(map[string]string),
		},
		HostIDs: hostIDs,
	}

	itemCount, err := session.GetItemCount(paramItem)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Item count: %d\n", itemCount)

	items, err := session.GetItems(paramItem)
	if err != nil {
		panic(err)
	}

	// Print history to file with host and item lookup
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if exportType == "snapshot" {
		snapshotCount := 0

		// Sort items by LastClock
		sort.Slice(items, func(i, j int) bool {
			return items[i].LastClock < items[j].LastClock
		})

		for _, item := range items {
			// skip if LastClock not between timeFrom and timeTill
			if item.LastClock < zbxTimeFrom.Unix() || item.LastClock >= zbxTimeTill.Unix() {
				continue
			}

			snapshotCount += 1
			text := fmt.Sprintf(
				"%s, Host=\"%s\", Groups=\"%s\", Key=\"%s\", Value=\"%s\"\n",
				time.Unix(int64(item.LastClock), 0).Format(time.RFC3339),
				hostmap[item.HostID].Hostname,
				groupmap[item.HostID],
				item.Key,
				item.LastValue)

			if _, err = f.WriteString(text); err != nil {
				panic(err)
			}
		}

		fmt.Printf("Exported %d last snapshots to %s\n", snapshotCount, outputFile)

	} else if exportType == "history" {
		itemmap := make(map[string]zabbix.Item)
		for _, item := range items {
			itemmap[item.ItemID] = item
		}

		// Get history ==========================
		zbxTimeFromFloat := float64(zbxTimeFrom.Unix())
		zbxTimeTillFloat := float64(zbxTimeTill.Unix())

		histories, err := session.GetHistories(zabbix.HistoryGetParams{
			TimeFrom: zbxTimeFromFloat,
			TimeTill: zbxTimeTillFloat,
		})
		if err != nil {
			panic(err)
		}

		// // Sort histories by Clock
		// sort.Slice(histories, func(i, j int) bool {
		// 	return histories[i].Clock < histories[j].Clock
		// })

		for _, history := range histories {
			_hostname := hostmap[itemmap[history.ItemID].HostID].Hostname
			text := fmt.Sprintf(
				"%s, Host=\"%s\", Key=\"%s\", Value=\"%s\"\n",
				time.Unix(history.Clock, 0).Format(time.RFC3339),
				_hostname,
				itemmap[history.ItemID].Key,
				history.Value)

			if _, err = f.WriteString(text); err != nil {
				panic(err)
			}
		}

		fmt.Printf("Exported %d histories to %s\n", len(histories), outputFile)
	}
}
