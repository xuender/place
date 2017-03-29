package main

import (
	"os"

	"./place"
	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "place"
	app.Usage = "自动移动文件分类存储."
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "xuender",
			Email: "xuender@139.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "调试模式",
			Hidden: false,
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "~/.place",
			Usage: "配置文件保存目录",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "move",
			Aliases: []string{"m"},
			Usage:   "移动文件、目录",
			Action: func(c *cli.Context) error {
				db := getDb(c)
				defer db.Close()
				actionInit(c)
				log.Info("移动操作")
				p := &place.Place{
					ConfigPath: c.GlobalString("config"),
					Db:         db,
					Preview:    false,
					History: place.History{
						Args: os.Args[1:],
					},
				}
				p.Run(c.Args())
				return nil
			},
		},
		{
			Name:    "preview",
			Aliases: []string{"p"},
			Usage:   "预览移动文件、目录",
			Action: func(c *cli.Context) error {
				db := getDb(c)
				defer db.Close()
				actionInit(c)
				log.Info("预览操作")
				p := &place.Place{
					ConfigPath: c.GlobalString("config"),
					Db:         db,
					Preview:    true,
				}
				p.Run(c.Args())
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "历史列表",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "clean, c",
					Usage:  "显示后，清除操作日志",
					Hidden: false,
				},
			},
			Action: func(c *cli.Context) error {
				db := getDb(c)
				defer db.Close()
				actionInit(c)
				place.ShowHistory(db)
				if c.Bool("clean") {
					place.CleanHistory(db)
				}
				return nil
			},
		},
		{
			Name:    "files",
			Aliases: []string{"f"},
			Usage:   "文件信息",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:   "clean, c",
					Usage:  "显示后，清除文件信息记录",
					Hidden: false,
				},
			},
			Action: func(c *cli.Context) error {
				db := getDb(c)
				defer db.Close()
				actionInit(c)
				place.ShowFiles(db)
				if c.Bool("clean") {
					place.CleanFiles(db)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func getDb(c *cli.Context) *leveldb.DB {
	db, err := leveldb.OpenFile(place.ToPath(c.GlobalString("config")+"/"+"db"), nil)
	if err != nil {
		panic("数据库创建失败")
	}
	return db
}
func actionInit(c *cli.Context) {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("调试状态: ", c.GlobalBool("debug"))
	log.Debug("配置文件目录: ", c.GlobalString("config"))
}
