package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var data DataBase

func ResetDataBase(c *cli.Context) {
	data.ResetDataBase()
}

func main() {
	data := InitializeDataBase()
	data.LoadData("Data/data.json")

	cliApp := cli.NewApp()
	cliApp.Name = "Groups-Management"
	cliApp.Version = "0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:   "reset",
			Usage:  "Removed the previous data.json file",
			Action: ResetDataBase,
		},
		{
			Name:      "createUser",
			Usage:     "Creates a new user in the data base",
			ArgsUsage: "Username",
			Action:    data.ResetDataBase,
		},
	}
	cliApp.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "debug, d",
			Value: 0,
			Usage: "debug-level: 1 for terse, 5 for maximal",
		},
	}
	cliApp.Before = func(c *cli.Context) error {
		//log.SetDebugVisible(c.Int("debug"))
		return nil
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	data.SaveData("Data/data.json")
}
