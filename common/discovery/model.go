package discovery

import "encoding/json"

// EndpointInfo 服务发现信息,这个就是存在etcd中的kv中的value
type EndpointInfo struct {
	IP       string                 `json:"ip"`
	Port     string                 `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	e := &EndpointInfo{}
	err := json.Unmarshal(data, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *EndpointInfo) Marshal() string {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(data)
}
