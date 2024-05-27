package jobserver


	"fmt"
	"io/ioutil`"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func 

func main() {
	// 读取私钥
	key, err := ioutil.ReadFile("path/to/private/key")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// 解析私钥
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// 创建SSH客户端配置
	config := &ssh.ClientConfig{
		User: "username",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// 连接到SSH服务器
	client, err := ssh.Dial("tcp", "hostname:22", config)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer client.Close()

	// 创建session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput("ls -l")
	if err != nil {
		log.Fatalf("failed to run command: %v", err)
	}
	fmt.Println(string(output))
}
