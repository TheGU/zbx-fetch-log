package main

import (
	"fmt"
	"log"
	"strconv"

	"os"
	"path/filepath"
	"sort"
	"time"

	"pattapongj/zbx-fetch-log/v6/zabbix"

	"github.com/Unknwon/goconfig"
	gsyslog "github.com/hashicorp/go-syslog"
)

func runZabbixExport(
	zbxHost, zbxUsername, zbxPassword,
	output, timeFrom string, timeTo string,
	profile string, cfg *goconfig.ConfigFile, allLog bool) {
	var zbxTimeTill time.Time
	var zbxTimeFrom time.Time
	var err error

	// Get call limit from config in int format
	// callLimit, _ := cfg.Int(profile, "Limit")
	exportType := cfg.MustValue(profile, "ExportType", "snapshots")
	callHostBatch := cfg.MustInt(profile, "CallHostBatch", 100)

	// Time variable to set time from and time to in history.get
	if timeTo != "" {
		zbxTimeTill, err = relativeToAbsoluteTime(timeTo)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		zbxTimeTill = time.Now()
	}

	if timeFrom != "" {
		zbxTimeFrom, err = relativeToAbsoluteTime(timeFrom)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		zbxTimeFrom, err = relativeToAbsoluteTime("5m")
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("time from", zbxTimeFrom)
	fmt.Println("time till", zbxTimeTill)

	// Default approach - without session caching
	session, err := zabbix.NewSession(zbxHost, zbxUsername, zbxPassword)
	if err != nil {
		log.Fatal(err)
	}

	version, err := session.GetVersion()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to Zabbix API v%s : %s\n", version, zbxHost)

	// Get hosts ==========================
	paramHost := zabbix.HostGetParams{
		GetParameters: zabbix.GetParameters{
			OutputFields: []string{"hostid", "host", "name", "status", "available"},
			// ResultLimit: 1000,
			SortField: []string{"hostid"},
			Filter:    map[string]interface{}{"status": 0},
			SortOrder: "ASC",
		},
		SelectHostGroups: []string{"groupid", "name"},
	}

	hostCount, err := session.GetHostCount(paramHost)
	if err != nil {
		log.Fatal(err)
	}

	hosts, err := session.GetHosts(paramHost)
	if err != nil {
		log.Fatal(err)
	}
	if len(hosts) != hostCount {
		log.Fatalf("Warning: Host count mismatch: %d != %d\n", len(hosts), hostCount)
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

	// loop through all hosts in batches of callHostBatch
	for i := 0; i < hostCount; i += callHostBatch {

		hostBatch := min(i+callHostBatch, hostCount)

		paramItem := zabbix.ItemGetParams{
			GetParameters: zabbix.GetParameters{
				OutputFields: []string{"hostid", "itemid", "key_", "prevvalue", "lastvalue", "value_type", "lastclock"},
				// ResultLimit:  callLimit,
				Filter: map[string]interface{}{"status": 0},
				// SortField:  []string{"itemid"},
				// SortOrder:  "ASC",
				TextSearch: make(map[string]interface{}),
			},
			HostIDs: hostIDs[i:hostBatch],
		}
		if !allLog {
			paramItem.TextSearch["key_"] = []string{"system", "vm", "vfs"}
			paramItem.SearchByAny = true
			paramItem.TextSearchByStart = true
		}

		itemCount, err := session.GetItemCount(paramItem)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Host Batch: %d-%d", i, hostBatch)
		fmt.Printf(", Item total: %d\n", itemCount)

		items, err := session.GetItems(paramItem)
		if err != nil {
			log.Fatal(err)
		}

		// Print history to file with host and item lookup
		programName := filepath.Base(os.Args[0])
		log.SetFlags(0)
		log.SetPrefix("")

		if output == "stdout" {
			log.SetOutput(os.Stdout)
			fmt.Println("Logging to stdout")
		} else if output == "tcp" || output == "udp" {
			syslogServer := cfg.MustValue(profile, "SyslogServer", "localhost:514")
			syslogFacility := cfg.MustValue(profile, "SyslogFacility", "LOCAL7")
			fmt.Printf("SyslogServer: %s\n", syslogServer)
			fmt.Printf("SyslogFacility: %s\n", syslogFacility)

			sysLog, err := gsyslog.DialLogger(output, syslogServer, gsyslog.LOG_INFO, syslogFacility, programName)
			if err != nil {
				log.Fatal(err)
			}
			defer sysLog.Close()

			log.SetOutput(sysLog)
			fmt.Printf("Logging to syslog server %s://%s INFO.%s\n", output, syslogServer, syslogFacility)

		} else if output == "syslog" {
			syslogFacility := cfg.MustValue(profile, "SyslogFacility", "LOCAL7")
			fmt.Printf("SyslogFacility: %s\n", syslogFacility)

			sysLog, err := gsyslog.NewLogger(gsyslog.LOG_INFO, syslogFacility, programName)
			if err != nil {
				log.Fatal(err)
			}
			defer sysLog.Close()
			log.SetOutput(sysLog)
			fmt.Printf("Logging to syslog %s INFO.%s\n", programName, syslogFacility)
		} else {
			f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			log.SetOutput(f)
			fmt.Printf("Logging to file %s\n", output)
		}

		if exportType == "snapshot" {
			snapshotCount := 0
			snapshotSkip := 0

			// Sort items by LastClock
			sort.Slice(items, func(i, j int) bool {
				return items[i].LastClock < items[j].LastClock
			})

			for _, item := range items {
				// skip if LastClock not between timeFrom and timeTill
				if item.LastClock < zbxTimeFrom.Unix() || item.LastClock >= zbxTimeTill.Unix() {
					snapshotSkip += 1
					continue
				}
				snapshotCount += 1

				// Attempt to parse LastValue and PrevValue as floats
				lastValueFloat, errLast := strconv.ParseFloat(item.LastValue, 64)
				prevValueFloat, errPrev := strconv.ParseFloat(item.PrevValue, 64)

				// Prepare formatted strings for LastValue and PrevValue
				var formattedLastValue, formattedPrevValue string
				if errLast == nil {
					formattedLastValue = fmt.Sprintf("%.4g", lastValueFloat)
				} else {
					formattedLastValue = item.LastValue // Keep original string if not a float
				}
				if errPrev == nil {
					formattedPrevValue = fmt.Sprintf("%.4g", prevValueFloat)
				} else {
					formattedPrevValue = item.PrevValue // Keep original string if not a float
				}

				text := fmt.Sprintf(
					"Time=\"%s\", HostName=\"%s\", Host=\"%s\", Groups=\"%s\", Key=\"%s\", Value=\"%s\", PrevValue=\"%s\"\n",
					time.Unix(int64(item.LastClock), 0).Format(time.RFC3339),
					hostmap[item.HostID].Hostname,
					hostmap[item.HostID].DisplayName,
					groupmap[item.HostID],
					item.Key,
					formattedLastValue,
					formattedPrevValue,
				)

				log.Print(text)
			}

			fmt.Printf("Exported total=%d, write=%d, skip=%d (outside selected time)\n", snapshotCount+snapshotSkip, snapshotCount, snapshotSkip)

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
				log.Fatal(err)
			}

			// // Sort histories by Clock
			// sort.Slice(histories, func(i, j int) bool {
			// 	return histories[i].Clock < histories[j].Clock
			// })

			for _, history := range histories {
				text := fmt.Sprintf(
					"Time=\"%s\", Host=\"%s\", Groups=\"%s\", Key=\"%s\", Value=\"%s\"\n",
					time.Unix(history.Clock, 0).Format(time.RFC3339),
					hostmap[itemmap[history.ItemID].HostID].Hostname,
					groupmap[itemmap[history.ItemID].HostID],
					itemmap[history.ItemID].Key,
					history.Value)

				log.Print(text)
			}

			fmt.Printf("Exported %d histories\n", len(histories))
		}
	}
}
