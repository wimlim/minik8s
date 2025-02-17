package kubeproxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiobj"
	"minik8s/pkg/apirequest"
	"minik8s/pkg/message"
	"os"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

type KubeProxy struct {
	ipvsManager   *IPVSManager
	subscriber    *message.Subscriber
	dnsSubscriber *message.Subscriber
}

func NewKubeProxy() *KubeProxy {
	ipvsManager := NewIPVSManager()
	if ipvsManager == nil {
		fmt.Println("Failed to create IPVS manager")
		return nil
	}
	subscriber := message.NewSubscriber()
	dnsbuscriber := message.NewSubscriber()

	return &KubeProxy{
		ipvsManager:   ipvsManager,
		subscriber:    subscriber,
		dnsSubscriber: dnsbuscriber,
	}
}

func podMatchesService(pod *apiobj.Pod, service *apiobj.Service) bool {
	labels := pod.MetaData.Labels
	for key, value := range service.Spec.Selector {
		if currentValue, ok := labels[key]; !ok || currentValue != value {
			return false
		}
	}
	return true
}

func (kp *KubeProxy) handleServiceAdd(msg message.Message) {
	fmt.Println("handleServiceAdd")
	var service apiobj.Service
	if err := json.Unmarshal([]byte(msg.Content), &service); err != nil {
		fmt.Println("Failed to unmarshal service:", err)
		return
	}

	pods, err := apirequest.GetAllPods()
	if err != nil {
		fmt.Println("Failed to get all pods:", err)
		return
	}

	var podIPs []string
	for _, pod := range pods {
		if podMatchesService(&pod, &service) {
			podIPs = append(podIPs, pod.Status.PodIP)
		}
	}

	kp.ipvsManager.AddService(service.Spec, podIPs)
}

func (kp *KubeProxy) handleServiceDelete(msg message.Message) {
	fmt.Println("handleServiceDelete")
	var service apiobj.Service
	if err := json.Unmarshal([]byte(msg.Content), &service); err != nil {
		fmt.Println("Failed to unmarshal service:", err)
		return
	}

	kp.ipvsManager.DeleteService(service.Spec)
}

func (kp *KubeProxy) handleServiceUpdate(msg message.Message) {
	fmt.Println("handleServiceUpdate")
	var svc apiobj.PodSvcMsg
	err := json.Unmarshal([]byte(msg.Content), &svc)
	if err != nil {
		fmt.Println("Failed to unmarshal service message:", err)
		return
	}
	// iterate svcports & podports
	switch msg.Name {
	case "Add":
		for i, svcport := range svc.SvcPorts {
			podport := svc.PodPorts[i]
			// add ipvs rule
			isvcport, _ := strconv.Atoi(svcport)
			ipodport, _ := strconv.Atoi(podport)
			kp.ipvsManager.AddRule(svc.SvcIp, uint16(isvcport), svc.PodIp, uint16(ipodport))
		}
	case "Delete":
		for i, svcport := range svc.SvcPorts {
			podport := svc.PodPorts[i]
			// delete ipvs rule
			isvcport, _ := strconv.Atoi(svcport)
			ipodport, _ := strconv.Atoi(podport)
			kp.ipvsManager.DeleteRule(svc.SvcIp, uint16(isvcport), svc.PodIp, uint16(ipodport))
		}
	}
}

func (kp *KubeProxy) handleDNSAdd(msg message.Message) {
	fmt.Println("handleDNSAdd")
	hostname := msg.Name
	nginxip := msg.Content
	file := "/etc/hosts"
	// open hosts file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Failed to open hosts file:", err)
		return
	}
	defer f.Close()
	// append nginxip hostname
	_, err = f.WriteString(nginxip + " " + hostname + "\n")
	if err != nil {
		fmt.Println("Failed to write to hosts file:", err)
		return
	}
}

func (kp *KubeProxy) handleDNSDelete(msg message.Message) {
	fmt.Println("handleDNSDelete" + msg.Name)
	hostname := msg.Name
	file := "/etc/hosts"
	// open hosts file
	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Failed to open hosts file:", err)
		return
	}
	defer f.Close()
	// delete nginxip hostname
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), hostname) {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	// rewrite
	f.Truncate(0)
	f.Seek(0, 0)
	for _, line := range lines {
		f.WriteString(line + "\n")
	}
}

func (kp *KubeProxy) Run() {
	defer kp.subscriber.Close()
	defer kp.ipvsManager.Close()
	defer kp.dnsSubscriber.Close()

	hostname, _ := os.Hostname()
	svc_que := fmt.Sprintf(message.ServiceQueue+"-%s", hostname)
	dns_que := fmt.Sprintf(message.DnsQueue+"-%s", hostname)
	go func() {
		kp.subscriber.Subscribe(svc_que, func(d amqp.Delivery) {
			var msg message.Message
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				fmt.Println("unmarshal message error:", err)
				return
			}

			switch msg.Type {
			case "Add":
				kp.handleServiceAdd(msg)
			case "Delete":
				kp.handleServiceDelete(msg)
			case "Update":
				kp.handleServiceUpdate(msg)
			}
		})
	}()

	go func() {
		kp.dnsSubscriber.Subscribe(dns_que, func(d amqp.Delivery) {
			var msg message.Message
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				fmt.Println("unmarshal message error:", err)
				return
			}

			switch msg.Type {
			case "Add":
				kp.handleDNSAdd(msg)
			case "Delete":
				kp.handleDNSDelete(msg)
			}
		})
	}()

	select {}
}
