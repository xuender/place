package place

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/h2non/filetype.v1"
)

type Place struct {
	ConfigPath string
	Config     Config
	Db         *leveldb.DB
	Preview    bool
	Number     int
	History    History
}

func (p *Place) Run(files []string) {
	p.loadConfig()
	p.getNumber()
	p.History.Files = make(map[string]string)
	log.Debug("Run:", files)
	for _, v := range files {
		p.run(v)
	}
	if !p.Preview {
		p.saveHistory()
	}
}

func (p *Place) getNumber() {
	name := []byte("count")
	bs, _ := p.Db.Get(name, nil)
	p.Number = Bytes2Int(bs)
	p.Number += 1
	log.Debugf("操作次数: %d", p.Number)
	p.Db.Put(name, Int2Bytes(p.Number), nil)
}

func (p *Place) saveHistory() {
	log.Debugf("保存历史: %v", p.History.Files)
	p.History.Timestamp = time.Now()
	bf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(bf)
	err := enc.Encode(p.History)
	if err != nil {
		panic("操作记录编码失败")
	}
	p.Db.Put([]byte(fmt.Sprintf("no-%04d", p.Number)), bf.Bytes(), nil)
}

func (p *Place) loadConfig() {
	p.Config = Config{}
	p.Config.Load(p.ConfigPath)
}

func (p *Place) run(file string) {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		log.Error("文件未找到: ", file)
		return
	}
	log.Debugf("%s 目录: %t", file, info.IsDir())
	if info.IsDir() {
		log.Info("处理目录: ", file)
	} else {
		log.Debug("处理文件: ", file)
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
			log.Debugf("扩展文件名: %s type: %s, subType: %s", path.Ext(file), kind.MIME.Type, kind.MIME.Subtype)
			hash := sha256.New()
			hash.Write(buf)
			log.Debugf("%s sha256: %x", file, hash.Sum(nil))
			old := p.find(hash.Sum(nil))
			if old == "" {
				newFile := p.move(kind.MIME.Type, kind.MIME.Subtype, file)
				p.Db.Put(BytesPrefix("f-", hash.Sum(nil)), []byte(newFile), nil)
			} else {
				log.Warnf("文件 %s 重复, 原文件 %s", file, old)
			}
		}
	}
}

func (p *Place) find(hash []byte) string {
	file, err := p.Db.Get(BytesPrefix("f-", hash), nil)
	if err != nil {
		return ""
	}
	return string(file)
}

func (p *Place) move(mime string, subtype string, file string) string {
	log.Debugf("mime: %s, subtype: %s, 搬移文件: %s", mime, subtype, file)
	ext := path.Ext(file)
	for _, path := range p.Config.Paths {
		if path.Mime == mime && path.Subtype == subtype {
			if path.Ext != "" || path.Ext == ext {
				dir := ToPath(path.Dir)
				p.History.Files[file] = dir
				log.Infof("文件: %s >>> %s", file, dir)
				if p.Preview {
					log.Debug("演示执行...")
					return file
				} else {
					log.Debug("真实执行...")
					// TODO 文件移动
					// TODO 移动成功 后的文件
					return file
				}
			}
		}
	}
	return file
}
