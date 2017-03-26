package place

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"gopkg.in/h2non/filetype.v1"
	"crypto/sha256"
	"github.com/syndtr/goleveldb/leveldb"
	"bytes"
	"encoding/binary"
	"time"
	"encoding/gob"
)

type Place struct {
	ConfigPath	string
	Config Config
	Db *leveldb.DB
	Number int32
	Files map[string]string
}

func (p *Place) Run(files []string) {
	p.loadConfig()
	p.getNumber()
	p.Files = make(map[string]string)
	log.Debug("Run:", files)
	for _, v := range files {
		p.run(v)
	}
	p.saveHistory()
}

func (p *Place) getNumber() {
	name := []byte("count")
	bs , _:= p.Db.Get(name, nil)
	p.Number = Bytes2Int(bs)
	p.Number += 1
	log.Debugf("操作次数: %d", p.Number)
	p.Db.Put(name, Int2Bytes(p.Number), nil)
}

func (p *Place) saveHistory() {
	log.Debugf("保存历史: %v", p.Files)
	h := History{
		Timestamp: time.Now(),
		Files: p.Files,
	}
	buf := bytes .NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, p.Number)
	bf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(bf)
	err := enc.Encode(h)
	if err != nil {
		panic("操作记录编码失败")
	}
	p.Db.Put(buf.Bytes(), bf.Bytes(), nil)

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
			hash := sha256.New()
			hash.Write(buf)
			log.Debugf("%s sha256: %x", file, hash.Sum(nil))
			p.move(kind.MIME.Type, kind.MIME.Subtype, file)
		}
	}
}

func (p *Place) move(mime string, subtype string, file string) {
	log.Debugf("mime: %s, subtype: %s, 搬移文件: %s", mime, subtype, file)
	for _, path := range p.Config.Paths {
		if path.Mime == mime && path.Subtype == subtype {
			dir := ToPath(path.Dir)
			p.Files[file] = dir
			log.Debugf("文件: %s >>> %s", file, dir)
			return;
		}
	}
}
