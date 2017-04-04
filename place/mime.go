package place

import (
	"bytes"
	"os/exec"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/h2non/filetype.v1"
	"gopkg.in/h2non/filetype.v1/types"
)

func Mime(bs []byte, file string) types.Type {
	kind, err := filetype.Match(bs)
	if err != nil || kind == filetype.Unknown || kind.MIME.Subtype == "zip" {
		cmd := exec.Command("file", "-i", file)
		cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			s := strings.Split(out.String(), " ")
			if len(s) > 1 {
				log.Debug("新类型: ", s[1])
				return types.NewType(path.Ext(file), s[1][:len(s[1])-1])
			}
		}
	}
	return kind
}
