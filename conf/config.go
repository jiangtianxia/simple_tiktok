package conf

type DefaultConf struct {
	MachineId      uint
	BucketRate     int64
	BucketCapacity int64
	COSAddr        string
	COSSecretId    string
	COSSecretKey   string
	Md5Salt        string
	JwtKey         string
	JwtExpire      int
	HashSalt       string
	UploadBase     string
	UploadAddr     string
}

type ServerConf struct {
	HTTPPort        uint
	GRPCPort        uint
	Logstash        bool
	LogstashConn    string
	LogEnableFile   bool
	LogPath         string
	LogFile         string
	LogRotationTime int64
	LogMaxAge       int64
	AllowOrigins    []string
	OpenCORS        bool
}

type DbConf struct {
	Name   string // default tiktok
	Driver string // 数据库类型 mysql,postgres,sqlite,mssql
	Source string // 连接字符串
}

type RedisConf struct {
	Name          string // default tiktok
	RedisHost     string
	RedisPassword string
	RedisPort     int
	RedisDB       int
}

var (
	globalConf *DefaultConf
)

func GetGlobalConf() *DefaultConf {
	return globalConf
}

func SetGlobalConf(c *DefaultConf) {
	globalConf = c
}
