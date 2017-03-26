package main

import (
	"os"
	"github.com/urfave/cli"
	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"./place"
)

func main() {
	app := cli.NewApp()
	app.Name = "place"
	app.Usage = "自动移动文件分类存储."
	app.Version = "0.0.1"
	app.Authors = []cli.Author {
		cli.Author{
			Name: "xuender",
			Email: "xuender@139.com",
		},
	}

	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "debug, d",
			Usage: "调试模式",
			Hidden: false,
		},
		cli.StringFlag{
			Name: "config, c",
			Value: "~/.place",
			Usage: "配置文件保存目录",
		},
		cli.BoolFlag{
			Name: "history, i",
			Usage: "显示操作历史",
			Hidden: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 && !c.Bool("history") {
			return cli.ShowAppHelp(c)
		}
		db, err := leveldb.OpenFile(place.ToPath(c.String("config") + "/" + "db"), nil)
		defer db.Close()
		if c.Bool("history") {
			place.ShowHistory(db)
			return nil
		}
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		log.Debug("调试状态: ", c.Bool("debug"))
		log.Debug("配置文件目录: ", c.String("config"))
		log.Debug("输入文件: ", c.Args())

		if err != nil {
			panic("数据库创建失败")
		}
		p := &place.Place{
			ConfigPath: c.String("config"),
			Db: db,
		}
		p.Run(c.Args())
		return nil
	}
	app.Run(os.Args)
}
