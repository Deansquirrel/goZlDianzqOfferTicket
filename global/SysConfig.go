package global

type SysConfig struct {
	TotalConfig   totalConfig   `toml:"Total"`
	PeiZhDbConfig peiZhDbConfig `toml:"PeiZhDb"`
}

type totalConfig struct {
	MinTicketNum int    `toml:"minTicketNum"`
	AppId        string `toml:"appid"`
	JPeiZh       string `toml:"jpeizh"`
	Port         int    `toml:"port"`
	IsDebug      bool   `toml:"debug"`
}

type peiZhDbConfig struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	DbName   string `toml:"dbName"`
	User     string `toml:"user"`
	PassWord string `toml:"password"`
}
