package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/battlesrv/go-gsstat/minecraft"
	"github.com/battlesrv/go-gsstat/steam"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	app := cli.NewApp()
	app.Author = "Konstantin Kruglov"
	app.Email = "kruglovk@gmail.com"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "steam",
			Aliases: []string{"s"},
			Action:  steamStat,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "info",
					Usage: "get information about targer",
				},
				cli.BoolFlag{
					Name:  "players",
					Usage: "get information about players",
				},
				cli.BoolFlag{
					Name:  "rules",
					Usage: "get information about rules",
				},
				cli.StringFlag{
					Name:  "addr",
					Usage: "host:port of game server",
				},
			},
		},
		{
			Name:    "minecraft",
			Aliases: []string{"m"},
			Action:  minecraftStat,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "addr",
					Usage: "host:port of game server",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func steamStat(c *cli.Context) {
	addr := c.String("addr")

	if c.Bool("info") {
		if steamInfo, err := steam.GetInfo(addr, time.Second*5); err != nil {
			log.Fatalln(err)
		} else {
			outputJSON(steamInfo)
		}
	}
	if c.Bool("players") {
		if steamPlayers, err := steam.GetPlayers(addr, time.Second*5); err != nil {
			log.Fatalln(err)
		} else {
			outputJSON(steamPlayers)
		}
	}
	if c.Bool("rules") {
		if steamRules, err := steam.GetRules(addr, time.Second*5); err != nil {
			log.Fatalln(err)
		} else {
			outputJSON(steamRules)
		}
	}
}

func minecraftStat(c *cli.Context) {
	if mstats, err := minecraft.GetStats(c.String("addr"), time.Second*5); err != nil {
		log.Fatalln(err)
	} else {
		outputJSON(mstats)
	}
}

func outputJSON(i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
