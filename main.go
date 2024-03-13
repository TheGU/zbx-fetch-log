package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	profile    string
	command    string
	outputFile string
	timeFrom   string
	key        = []byte("FQ@dhF#Lvo!ZtxA9ArnNdF!aeZZRdxiQ") // Global key variable

)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	setupFlagSet := flag.NewFlagSet("setup", flag.ExitOnError)
	runFlagSet := flag.NewFlagSet("run", flag.ExitOnError)

	setupFlagSet.StringVar(&profile, "profile", "", "Profile name")
	runFlagSet.StringVar(&profile, "profile", "", "The profile name")
	runFlagSet.StringVar(&outputFile, "output", "output.txt", "The output file")
	runFlagSet.StringVar(&timeFrom, "timeFrom", "5m", "The relative time from now")

	switch os.Args[1] {
	case "setup":
		setupFlagSet.Parse(os.Args[2:])
		setupProfile(profile, key)
	case "run":
		runFlagSet.Parse(os.Args[2:])
		runProfile(profile, key, outputFile, timeFrom, "")
	default:
		fmt.Println("Invalid command. Use setup or run.")
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: zbx-fetch-log [setup|run] --profile PROFILE_NAME [--output.txt OUTPUT_FILE] [--timeFrom TIME_FROM]

    --profile: The profile name.
    --output.txt: The output file. Default is "output.txt".
    --timeFrom: The relative time from now. Default is "5m".`)
}
