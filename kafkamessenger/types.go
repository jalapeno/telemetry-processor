package kafkamessenger

import (
	"encoding/json"
	"fmt"
	"log"
)

type Fields struct {
	AfName               string `json:"af/af_name"`
	AfRouteTargetAfName  string `json:"af/route_target/af_name"`
	AfRouteTargetType    string `json:"af/route_target/route_target_type"`
	AfRouteTargetValue   string `json:"af/route_target/route_target_value"`
	AfRouteTargetSafName string `json:"af/route_target/saf_name"`
	AfSafName            string `json:"af/saf_name"`
	InterfaceName        string `json:"interface/interface_name"`
	IsBigVrf             string `json:"is_big_vrf"`
	RouteDistinguisher   string `json:"route_distinguisher"`
	VrfNameXr            string `json:"vrf_name_xr"`
}

type Tags struct {
	Host         string `json:"host"`
	Path         string `json:"path"`
	Source       string `json:"source"`
	Subscription string `json:"subscription"`
	VrfName      string `json:"vrf_name"`
}

type Vrf struct {
	Fields    Fields `json:"fields"`
	Name      string `json:"name"`
	Tags      Tags   `json:"tags"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	jsonStr := `{
        "fields": {
            "af/af_name": "ipv6",
            "af/route_target/af_name": "ipv6",
            "af/route_target/route_target_type": "export",
            "af/route_target/route_target_value": "9:9",
            "af/route_target/saf_name": "unicast",
            "af/saf_name": "unicast",
            "interface/interface_name": "Loopback9",
            "is_big_vrf": "false",
            "route_distinguisher": "10.0.0.1:0",
            "vrf_name_xr": "red"
        },
        "name": "Cisco-IOS-XR-mpls-vpn-oper:l3vpn/vrfs/vrf",
        "tags": {
            "host": "telegraf",
            "path": "Cisco-IOS-XR-mpls-vpn-oper:l3vpn/vrfs/vrf",
            "source": "xrd01",
            "subscription": "base_metrics",
            "vrf_name": "red"
        },
        "timestamp": 1719867701
    }`

	var telemetry Telemetry

	err := json.Unmarshal([]byte(jsonStr), &telemetry)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", telemetry)
}
