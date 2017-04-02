package place

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
)

type Path struct {
	Dir     string `json:"dir"`
	Mime    string `json:"mime,omitempty"`
	Subtype string `json:"subtype,omitempty"`
	Ext     string `json:"ext,omitempty"`
	Subdir  string `json:"subdir,omitempty"`
}

type Config struct {
	Paths []*Path `json:"paths"`
}

func (c *Config) Load(configPath string) {
	dir := ToPath(configPath)

	file := path.Join(dir, "config.json")
	log.Debug("加载配置: ", file)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Debugf("%s 不存在，创建", file)
		c.def()
		file, _ := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0777)
		defer file.Close()
		enc := json.NewEncoder(file)
		enc.SetIndent("  ", "  ")
		if enc.Encode(*c) != nil {
			panic("配置文件保存失败")
		}
	} else {
		log.Debugf("%s 存在，开始加载", file)
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic("配置文件读取失败")
		}
		json.Unmarshal(bytes, c)
	}
	log.Debugf("配置内容: %s", c)
}

func (c *Config) def() {
	log.Debug("创建默认设置")
	c.Paths = []*Path{
		&Path{
			Dir:     "~/图片/截图",
			Mime:    "image",
			Subtype: "png",
		},
		&Path{
			Dir:     "~/图片/照片",
			Mime:    "image",
			Subtype: "jpeg",
			Subdir:  "yyyy/mm",
		},
		&Path{
			Dir:     "~/文档/归档",
			Mime:    "application",
			Subtype: "x-7z-compressed",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:     "~/文档/设计",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".xmind",
		},
		&Path{
			Dir:    "~/文档/办公",
			Ext:    ".ppt",
			Subdir: "yyyy",
		},
		&Path{
			Dir:     "~/文档/办公",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".pptx",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:     "~/文档/办公",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".docx",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:    "~/文档/办公",
			Ext:    ".doc",
			Subdir: "yyyy",
		},
		&Path{
			Dir:     "~/文档/图书",
			Mime:    "application",
			Subtype: "pdf",
		},
	}
}
