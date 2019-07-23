package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var data DataBase

func ResetDataBase(c *cli.Context) {
	data.ResetDataBase()
}

func CreateUser(c *cli.Context) {
	if c.NArg() != 1 {
		_ = errors.New("please give the following argument: UserID")
		return
	}
	userid := c.Args().First()
	_ = CreateNewUser(UserID(userid))
}

func main() {
	data = InitializeDataBase()
	data.LoadData("Data/data.json")
	fmt.Println("Starting\t", data)

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
			ArgsUsage: "UserID",
			Action:    CreateUser,
		},
	}
	cliApp.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "debug, d",
			Value: 0,
			Usage: "debug-level: 1 for terse, 5 for maximal",
		},
	}
	cliApp.After = func(c *cli.Context) error {
		return nil
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	data.SaveData("Data/data.json")
	fmt.Println("Ending\t", data, "\n\n")
}
