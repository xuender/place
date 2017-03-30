package place

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"gopkg.in/h2non/filetype.v1"
)

type Place struct {
	ConfigPath string
	Config     Config
	Db         *leveldb.DB
	Preview    bool
	Number     int
	History    History
	TmpFile    map[string][]byte
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
		bs, _ := ioutil.ReadFile(file)
		head := bs[:261]
		if filetype.IsImage(head) {
			log.Debugf("文件 %s 是图片", file)
		} else {
			log.Debugf("文件 %s 不是图片", file)
		}
		hash := sha256.New()
		hash.Write(bs)
		log.Debugf("%s sha256: %x", file, hash.Sum(nil))
		kind, _ := filetype.Match(head)
		if kind == filetype.Unknown {
			log.Debug("未知的文件类型")
		} else {
			log.Debugf("扩展文件名: %s type: %s, subType: %s", path.Ext(file), kind.MIME.Type, kind.MIME.Subtype)
			old := p.find(hash.Sum(nil))
			if old == "" {
				newFile, err := p.move(kind.MIME.Type, kind.MIME.Subtype, file)
				if err == nil {
					if !p.Preview {
						p.Db.Put(BytesPrefix("f-", hash.Sum(nil)), []byte(newFile), nil)
					}
				} else {
					log.Errorf("%s: %s", file, err)
				}
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

func (p *Place) move(mime string, subtype string, file string) (string, error) {
	log.Debugf("mime: %s, subtype: %s, 搬移文件: %s", mime, subtype, file)
	ext := path.Ext(file)
	for _, ap := range p.Config.Paths {
		if ap.Mime == mime && ap.Subtype == subtype {
			if ap.Ext == "" || ap.Ext == ext {
				dir := ToPath(ap.Dir)
				p.History.Files[file] = dir
				newFile := path.Join(dir, path.Base(file))
				info, err := os.Stat(newFile)
				if os.IsNotExist(err) {
					if p.Preview {
						log.Infof("预览: %s >>> %s", file, newFile)
					} else {
						log.Infof("移动: %s >>> %s", file, newFile)
						os.Rename(file, newFile)
					}
					return newFile, nil
				}
				if info.IsDir() {
					return file, errors.New("同名目录已经存在")
				} else {
					return file, errors.New("同名文件已经存在")
				}
			}
		}
	}
	return file, errors.New("无匹配目录")
}

func (p *Place) Scan() {
	log.Info("开始目录扫描")
	p.loadConfig()

	p.TmpFile = make(map[string][]byte)
	iter := p.Db.NewIterator(util.BytesPrefix([]byte("f-")), nil)
	for iter.Next() {
		file := string(iter.Value())
		p.TmpFile[file] = iter.Key()
		log.Debug("缓存文件: ", file)
	}
	iter.Release()

	for _, ap := range p.Config.Paths {
		log.Debug("扫描目录: ", ap.Dir)
		p.scanning(ToPath(ap.Dir))
	}

	for file, key := range p.TmpFile {
		log.Debug("删除缓存: ", file)
		p.Db.Delete(key, nil)
	}
}

func (p *Place) scanning(dir string) {
	err := filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
		if filename == dir {
			return nil
		}
		if fi.IsDir() {
			p.scanning(filename)
		} else {
			log.Debug("扫描文件: ", filename)
			_, ok := p.TmpFile[filename]
			if ok {
				delete(p.TmpFile, filename)
			} else {
				bs, _ := ioutil.ReadFile(filename)
				hash := sha256.New()
				hash.Write(bs)
				key := BytesPrefix("f-", hash.Sum(nil))
				old, err := p.Db.Get(key, nil)
				if err == nil {
					log.Warnf("文件重复: %s = %s", filename, old)
				} else {
					log.Info("新增文件: ", filename)
					p.Db.Put(key, []byte(filename), nil)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Error("目录扫描错误: ", err)
	}
}
