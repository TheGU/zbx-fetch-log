package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	profile  string
	output   string
	timeFrom string
	allLog   bool
	key      = []byte("FQ@dhF#Lvo!ZtxA9ArnNdF!aeZZRdxiQ") // Global key variable

)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	setupFlagSet := flag.NewFlagSet("setup", flag.ExitOnError)
	runFlagSet := flag.NewFlagSet("run", flag.ExitOnError)

	setupFlagSet.StringVar(&profile, "profile", "", "Profile name")
	runFlagSet.BoolVar(&allLog, "allLog", false, "Fetch all log from zabbix server. Default is false. only log start with system, vm, vfs")
	runFlagSet.StringVar(&profile, "profile", "", "The profile name")
	runFlagSet.StringVar(&output, "output", "output.txt", "The output format: stdout, tcp, udp, syslog or filename. Default is file output.txt")
	runFlagSet.StringVar(&timeFrom, "timeFrom", "5m", "The relative time from now")

	switch os.Args[1] {
	case "setup":
		setupFlagSet.Parse(os.Args[2:])
		setupProfile(profile, key)
	case "run":
		runFlagSet.Parse(os.Args[2:])
		runProfile(profile, key, output, timeFrom, "", allLog)
	default:
		fmt.Println("Invalid command. Use setup or run.")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: zbx-fetch-log [setup|run] --profile PROFILE_NAME [--output.txt OUTPUT_FILE] [--timeFrom TIME_FROM] [--allLog]

    --profile: The profile name.
    --output output.txt: The output file. Default is "output.txt".
    --timeFrom: The relative time from now. Default is "5m".
    --allLog: Fetch all log from zabbix server. Default is false. only log start with system, vm, vfs

	`)
}
