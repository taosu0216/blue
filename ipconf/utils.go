package ipconf

import "blue/ipconf/domain"

// TODO: EndPort
func top5EndPort(ips []*domain.Endport) []*domain.Endport {
	if len(ips) <= 5 {
		return ips
	} else {
		return ips[:5]
	}
}

func packResp(endPoint []*domain.Endport) Response {
	return Response{
		Message: "ok",
		Code:    0,
		Data:    endPoint,
	}
}
