package place

import (
	"bytes"
	"encoding/gob"

	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func ShowHistory(db *leveldb.DB) {
	log.Info("显示操作历史")
	iter := db.NewIterator(util.BytesPrefix([]byte("no-")), nil)
	for iter.Next() {
		var h History
		buf := bytes.NewBuffer(iter.Value())
		dec := gob.NewDecoder(buf)
		dec.Decode(&h)
		log.Infof("%s [ %s ] 参数: %s", string(iter.Key()[3:]), h.Timestamp.Format("2006-01-02 15:04:05"), h.Args)
		for k, v := range h.Files {
			log.Infof("    %s >>> %s", k, v)
		}
	}
	iter.Release()
}

func CleanHistory(db *leveldb.DB) {
	log.Info("清空操作历史")
	iter := db.NewIterator(util.BytesPrefix([]byte("no-")), nil)
	for iter.Next() {
		log.Debug("删除记录: ", string(iter.Key()))
		db.Delete(iter.Key(), nil)
	}
	iter.Release()
}
