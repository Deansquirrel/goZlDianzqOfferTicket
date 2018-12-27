package global

type SysConfig struct {
	TotalConfig    totalConfig    `toml:"Total"`
	PeiZhDbConfig  peiZhDbConfig  `toml:"PeiZhDb"`
	RabbitMqConfig rabbitMqConfig `toml:"RabbitMQ"`
	RedisConfig    redisConfig    `toml:"Redis"`
}

type totalConfig struct {
	MinTicketNum int    `toml:"minTicketNum"`
	SnoWorkerId  int    `toml:"snoWorkerId"`
	SnoServer    string `toml:"snoServer"`
	IsDebug      bool   `toml:"debug"`
}

type peiZhDbConfig struct {
	AppId    string `toml:"appid"`
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	DbName   string `toml:"dbName"`
	User     string `toml:"user"`
	PassWord string `toml:"password"`
}

type rabbitMqConfig struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	User     string `toml':"user"`
	Password string `toml:"password"`
}

type redisConfig struct {
	Server    string `toml:"server"`
	Password  string `toml:"password"`
	DbId1     int    `toml:"dbId1"`
	DbId2     int    `toml:"dbId2"`
	SessionDb int    `toml:"sessionDb"`
}
