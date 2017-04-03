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
			Subtype: "png, bmp",
		},
		&Path{
			Dir:     "~/图片/照片",
			Mime:    "image",
			Subtype: "jpeg, tiff, x-canon-cr2",
			Subdir:  "yyyy/mm",
		},
		&Path{
			Dir:     "~/图片/动画",
			Mime:    "image",
			Subtype: "gif, webp",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:     "~/图片/动画",
			Mime:    "application",
			Subtype: "x-shockwave-flash",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:    "~/视频",
			Mime:   "video",
			Subdir: "yyyy",
		},
		&Path{
			Dir:  "~/音乐",
			Mime: "audio",
		},
		&Path{
			Dir:     "~/文档/字体",
			Mime:    "application",
			Subtype: "font-",
		},
		&Path{
			Dir:     "~/文档/设计",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".xmind",
		},
		&Path{
			Dir:     "~/文档/设计",
			Mime:    "image",
			Subtype: "vnd.adobe.photoshop, image/x-icon, vnd.ms-photo",
		},
		&Path{
			Dir:     "~/文档/办公",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".pptx, .docx",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:     "~/文档/办公",
			Mime:    "application",
			Subtype: "rtf",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:    "~/文档/办公",
			Ext:    ".doc .ppt",
			Subdir: "yyyy",
		},
		&Path{
			Dir:     "~/文档/图书",
			Mime:    "application",
			Subtype: "pdf, epub+zip",
		},
		&Path{
			Dir: "~/文档/图书",
			Ext: ".mobi",
		},
		&Path{
			Dir:     "~/文档/归档",
			Mime:    "application",
			Subtype: "x-tar, x-rar-compressed, gzip, x-bzip2, x-7z-compressed, x-xz, x-unix-archive, x-compress, x-lzip",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:     "~/文档/归档",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".zip",
			Subdir:  "yyyy",
		},
		&Path{
			Dir:     "~/文档/应用",
			Mime:    "application",
			Subtype: "x-msdownload, x-deb, x-google-chrome-extension, vnd.ms-cab-compressed",
		},
		&Path{
			Dir:     "~/文档/手机",
			Mime:    "application",
			Subtype: "zip",
			Ext:     ".apk",
		},
	}
}
