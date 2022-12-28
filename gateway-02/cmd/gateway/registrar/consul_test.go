package registrar

import (
	"fmt"
	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
	"net/http"
	"testing"
)

const (
	consulAddr  = "127.0.0.1:8500"
	localIp     = "127.0.0.1"
	localPort   = 9000
	srvName     = "srv1"
	servicePort = 9100
)

/*
测试通过接口向consul注册服务
*/
func TestRegister(t *testing.T) {
	//先启动一个HTTP服务（即要注册到consul的服务）
	done := make(chan bool, 1)
	go func() {
		http.HandleFunc("/", Handler) //只有健康检查处理
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", localIp, localPort), nil)
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
		done <- true
	}()

	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = uuid.NewString()
	registration.Name = srvName
	registration.Port = servicePort
	registration.Tags = []string{"test", "srv1"}
	registration.Address = localIp

	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, localPort)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check
	//more config
	//registration.Connect
	//registration.Kind
	//registration.Meta
	//registration.Namespace
	//registration.Proxy
	//registration.Weights

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)

	<-done
	fmt.Println("exit")
}

// 从consul中发现服务
func TestFind(t *testing.T) {
	// 创建consul连接
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return
	}

	// 获取指定service
	services, _ := client.Agent().Services()
	for _, service := range services {
		fmt.Printf("service address is:%s:%d service name is:%s\n", service.Address, service.Port, service.Service)
	}
	//service, _, err := client.Agent().Service("serviceId", nil)
	//if err != nil{
	//	fmt.Println("consul get service error : ", err)
	//	return
	//}
	//fmt.Printf("address is:%s:%d\n", service.Address, service.Port)

	//只获取健康的service
	hs, _, err := client.Health().Service(srvName, "", true, nil)
	if err != nil {
		fmt.Println("consul get health service error : ", err)
		return
	}

	fmt.Println("serviceHealthy address is:", hs[0].Node.Address)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("health check api"))
	fmt.Println("response health check api")
	if err != nil {
		_ = fmt.Errorf("err %v", err)
		return
	}
}
