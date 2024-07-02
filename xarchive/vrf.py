import json
import re

def parse_telemetry_string(telemetry_string):
    # Define the keys to extract and their corresponding desired keys in the output
    keys_to_extract = {
        "source": "source",
        "vrf_name": "vrf_name",
        "route_distinguisher": "route_distinguisher",
        "interface/interface_name": "interface_name",
        "af/af_name": "af_name",
        "af/route_target/route_target_type": "route_target_type",
        "af/route_target/route_target_value": "route_target_value"
    }

    # Use regular expressions to find all key-value pairs in the string
    kv_pairs = re.findall(r'(\S+?)=(?:"(.*?)"|(.*?))(?:,|\s|$)', telemetry_string)

    # Initialize the JSON dictionary
    json_dict = {}
    
    for key, value1, value2 in kv_pairs:
        # Determine which value to use
        value = value1 if value1 else value2
        # Check if the key is in the keys_to_extract set
        if key in keys_to_extract:
            json_dict[keys_to_extract[key]] = value

    # Combine 'source' and 'vrf_name' to create the '_key' field
    if "source" in json_dict and "vrf_name" in json_dict:
        json_dict["_key"] = f"{json_dict['source']}_vrf_{json_dict['vrf_name']}"

    return json_dict

# Example usage
telemetry_string = "Cisco-IOS-XR-mpls-vpn-oper:l3vpn/vrfs/vrf,host=telegraf,path=Cisco-IOS-XR-mpls-vpn-oper:l3vpn/vrfs/vrf,source=xrd01,subscription=base_metrics,vrf_name=red vrf_name_xr=\"red\",route_distinguisher=\"10.0.0.1:0\",interface/interface_name=\"Loopback9\",af/af_name=\"ipv6\",af/saf_name=\"unicast\",af/route_target/route_target_type=\"export\",af/route_target/route_target_value=\"9:9\",af/route_target/af_name=\"ipv6\",af/route_target/saf_name=\"unicast\",is_big_vrf=\"false\" 1719859230833000000"
parsed_json = parse_telemetry_string(telemetry_string)
print(json.dumps(parsed_json, indent=4))
