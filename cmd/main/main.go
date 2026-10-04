package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"

	config "clickfurnace/config"
	file "clickfurnace/file"
)

func main() {
	app := cli.NewApp()
	app.Name = "clickfurnace"
	app.Usage = "Generate data for ClickHouse"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "furnace-config.yml",
			Usage: "config to set up clickfurnace",
		},
		cli.BoolFlag{
			Name:  "run",
			Usage: "Run clickfurnace",
		},
	}

	app.Action = func(c *cli.Context) error {
		run := c.GlobalBool("run")
		if !run {
			fmt.Println("Use run flag to launch clickfurnace")
			return nil
		}
		configigName := c.GlobalString("config")
		if !file.CheckFileExists(configigName) {
			fmt.Printf("File %s does not exist\n", configigName)
			return nil
		}
		_, err := config.New(configigName)
		if err != nil {
			fmt.Printf("Some error: %v\n", err)
		}
		return nil
	}

	app.Run(os.Args)
}
