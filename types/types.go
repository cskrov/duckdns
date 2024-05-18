package types

type Config struct {
	Token      string   `json:"token"`
	Domains    []string `json:"domains"`
	Log        bool     `json:"log"`
	VerboseLog bool     `json:"verbose_log"`
	Interval   int      `json:"interval"`
}
