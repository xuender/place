package place

import (
	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func ShowFiles(db *leveldb.DB) {
	log.Info("显示记录的文件")
	iter := db.NewIterator(util.BytesPrefix([]byte("f-")), nil)
	for iter.Next() {
		log.Infof("sha256: %x file: %s", iter.Key()[2:], string(iter.Value()))
	}
	iter.Release()
}

func CleanFiles(db *leveldb.DB) {
	log.Info("清空文件记录")
	iter := db.NewIterator(util.BytesPrefix([]byte("f-")), nil)
	for iter.Next() {
		log.Debug("删除记录: ", string(iter.Value()))
		db.Delete(iter.Key(), nil)
	}
	iter.Release()
}
