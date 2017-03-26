package place

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"gopkg.in/h2non/filetype.v1"
)

type Place struct {
	ConfigPath	string
	Config Config
}

func (p *Place) Run(files []string) {
	p.loadConfig()
	log.Debug("Run:", files)
	for _, v := range files {
		p.run(v)
	}
}

func (p *Place) loadConfig() {
	p.Config = Config{}
	p.Config.Load(p.ConfigPath)
}

func (p *Place) run(file string) {
	info, err := os.Stat(file);
	if os.IsNotExist(err) {
		log.Error("文件未找到: ", file)
		return;
	}

	log.Debugf("%s 目录: %t", file, info.IsDir())
	if info.IsDir() {
		log.Info("处理目录: ", file)
	} else {
		log.Info("处理文件: ", file)
		buf, _ := ioutil.ReadFile(file)
		head := buf[:261]
		if filetype.IsImage(head) {
			log.Debugf("文件 %s 是图片", file)
		} else {
			log.Debugf("文件 %s 不是图片", file)
		}
		kind, _ := filetype.Match(head)
		if kind == filetype.Unknown {
			log.Debug("未知的文件类型")
		} else {
			log.Debugf("扩展文件名: %s", kind.Extension)
			p.move(kind.MIME.Type, kind.MIME.Subtype, file)
		}
	}
}

func (p *Place) move(mime string, subtype string, file string) {
	log.Debugf("mime: %s, subtype: %s, 搬移文件: %s", mime, subtype, file)
	for _, path := range p.Config.Paths {
		if path.Mime == mime && path.Subtype == subtype {
			dir := ToPath(path.Dir)
			log.Debugf("文件: %s >>> %s", file, dir)
			return;
		}
	}
}

