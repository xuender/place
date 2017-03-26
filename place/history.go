package place
import (
	"time"
	"github.com/syndtr/goleveldb/leveldb"
	log "github.com/Sirupsen/logrus"
	// "time"
	"bytes"
	"encoding/gob"
)

type History struct {
	Timestamp time.Time
	Files map[string]string
}

func ShowHistory(db *leveldb.DB) {
	log.Info("显示操作历史")
	name := []byte("count")
	bs , _:= db.Get(name, nil)
	number := Bytes2Int(bs)
	for i := 0 ; i <= int(number); i++ {
		data, err := db.Get(Int2Bytes(int32(i)), nil)
		if err != nil {
			log.Error(err)
		} else {
			var h History
			buf := bytes.NewBuffer(data)
			dec := gob.NewDecoder(buf)
			dec.Decode(&h)
			log.Infof("第 %d 操作， 时间: %s", i, h.Timestamp.Format("2006-01-02 15:04:05"))
			for k, v := range h.Files {
				 log.Infof("文件 %s >>> %s", k, v)
			}
		}
	}
}

