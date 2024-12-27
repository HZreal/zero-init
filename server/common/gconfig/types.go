package gconfig

type GlobalConfig struct {
	// 基础路径
	BasePath string

	Auth struct {
		AccessSecret         string
		AccessExpireSP       int64
		AccessExpireCustomer int64
	}

	RedisConf struct {
		Host     string
		Type     string
		Password string
		Tls      bool
	}

	Etcd struct {
		Hosts []string
	}

	// RPC服务列表
	RpcServerList struct {
		Overall    string
		Auth       string
		OperateLog string
		Maintain   string
		Operation  string
		Resource   string
		Route      string
		Query      string
		Common     string
	}

	// mydb1 数据库链接字符串
	MyDB1Link string

	// mydb2 数据库链接字符串
	MyDB2Link string

	// 日志客户端配置
	LogConf struct {
		CliConf struct {
			Enable        bool
			ServerSocket  string
			ClientSocket  string
			DefaultLevel  int
			CommandProces string
		}
		ServerConf struct {
			Enable       bool
			DefaultLevel int
			ServerInfo   string
			TPType       string
		}
		ConsoleConf struct {
			Enable bool
		}
		SyslogConf struct {
			Enable bool
		}
	}
}
