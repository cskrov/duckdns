package log

import "time"

type Timestamp int64

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Unix(int64(t), 0).UTC().Format("2006-01-02T15:04:05Z") + "\""), nil
}
