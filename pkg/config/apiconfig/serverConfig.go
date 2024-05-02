package apiconfig

import "fmt"
const(
	ServerDefaultListenIp = "0.0.0.0"
	ServerDefaultPort = 8080
	HttpScheme = "http://"
	ServerLocalIP = "127.0.0.1"
	ServerMasterIP = "10.119.13.134"
)

func GetServerLocalUrl()(string){
	ServerLocaltURL := fmt.Sprintf(HttpScheme + ServerLocalIP +":%d",ServerDefaultPort)
	return ServerLocaltURL;
}
