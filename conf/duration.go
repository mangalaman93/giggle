package conf

import (
	"fmt"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalJSON(b []byte) error {
	var err error
	sd := string(b[1 : len(b)-1])
	d.Duration, err = time.ParseDuration(sd)
	return err
}

func (d duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
