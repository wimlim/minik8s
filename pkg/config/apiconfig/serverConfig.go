package apiconfig
import(
	
)
const(
	ServerDefaultListenIp = "0.0.0.0"
	ServerDefaultPort = 8080
	ServerLocaltURL = "http://127.0.0.1:8080"
)

func getServerLocalUrl()(string){
	return ServerLocaltURL;
}
