package serverlessconfig

const (
	ServerDefaultPort     = 8081
	HttpScheme            = "http://"
	ServerLocalIP         = "127.0.0.1"
	ServerMasterIP        = "10.119.13.134"
)

func GetServerMasterIP() string {
	return ServerLocalIP
	// return ServerMasterIP
}
