package place

import (
	"time"

	"bytes"
	"encoding/gob"

	log "github.com/Sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

type History struct {
	Args      []string
	Timestamp time.Time
	Files     map[string]string
}

func ShowHistory(db *leveldb.DB) {
	log.Info("显示操作历史")
	name := []byte("count")
	bs, _ := db.Get(name, nil)
	number := Bytes2Int(bs)
	for i := number; i > 0; i-- {
		data, err := db.Get(Int2Bytes(i), nil)
		if err != nil {
			log.Error(err)
		} else {
			var h History
			buf := bytes.NewBuffer(data)
			dec := gob.NewDecoder(buf)
			dec.Decode(&h)
			log.Infof("%03d [ %s ] 参数: %s", i, h.Timestamp.Format("2006-01-02 15:04:05"), h.Args)
			for k, v := range h.Files {
				log.Infof("    %s >>> %s", k, v)
			}
		}
	}
}
