package rpc

type (
	User struct {
		Id         uint   `json:"id"`
		Fd         uint64 `json:"fd"`
		Nickname   string `json:"nickname"`
		Phone      string `json:"phone"`
		Type       string `json:"type"`
		System     string `json:"system"`
		AuthGroups []int  `json:"auth_groups"`
	}
	ServerParams struct {
		RequestTime      uint64  `json:"request_time"`
		RequestTimeMilli float64 `json:"request_time_milli"`
		ServerProtocol   string  `json:"server_protocol"`
		RemoteAddr       string  `json:"remote_addr"`
		RemoteIp         string  `json:"remote_ip"`
		GatewayId        string  `json:"gateway_id"`
	}
	RpcRequestData struct {
		FromId       uint              `json:"fromId"`
		FromInfo     User              `json:"fromInfo"`
		Timestamp    int64             `json:"timestamp"`
		Type         string            `json:"type"`
		RequestType  string            `json:"request_type"`
		Path         string            `json:"path"`
		Params       map[string]any    `json:"params"`
		Headers      map[string]string `json:"headers"`
		ServerParams ServerParams      `json:"serverParams"`
	}
)
