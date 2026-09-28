package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"

	utils "clickfurnance/utils"
)

func main() {
	app := cli.NewApp()
	app.Name = "clickfurnance"
	app.Usage = "Generate data for ClickHouse"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "furnance-config.yml",
			Usage: "Config to set up clickfurnance",
		},
		cli.BoolFlag{
			Name:  "run",
			Usage: "Run clickfurnance",
		},
	}

	app.Action = func(c *cli.Context) error {
		run := c.GlobalBool("run")
		if !run {
			fmt.Println("Use run flag to launch clickfurnance")
			return nil
		}
		configName := c.GlobalString("config")
		if !utils.CheckFileExists(configName) {
			fmt.Printf("File %s does not exist\n", configName)
			return nil
		}
		fmt.Printf("Initialized config %s\n", configName)
		return nil
	}

	app.Run(os.Args)
}
