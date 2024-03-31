package discov

type Service struct {
	// 整体的服务的名字,比如gateway_service
	Name      string      `json:"name"`
	Endpoints []*Endpoint `json:"endpoints"`
}

type Endpoint struct {
	// 每个节点的名字,比如gateway_service_machine1
	ServerName string `json:"server_name"`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
	Weight     int    `json:"weight"`
	Enable     bool   `json:"enable"`
}
