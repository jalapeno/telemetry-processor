package kafkamessenger

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
