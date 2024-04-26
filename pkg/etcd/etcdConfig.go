package etcd
import(
	"time"
)
type EtcdConfig struct{
	EtcdEndpoints []string
	EtcdDialTimeout time.Duration
}
func EtcdDefaultConfig() *EtcdConfig{
	return &EtcdConfig{
		EtcdEndpoints: []string{"localhost:2379"},
		EtcdDialTimeout:  5 * time.Second,
	}
}