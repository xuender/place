package place

import (
	"time"
)

type History struct {
	Args      []string
	Timestamp time.Time
	Files     map[string]string
}
