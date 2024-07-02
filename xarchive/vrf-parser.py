import json

def parse_telemetry_string(telemetry_string):
    # Define the keys to extract and their corresponding desired keys in the output
    keys_to_extract = {
        "source": "source",
        "vrf_name": "vrf_name",
        "route_distinguisher": "route_distinguisher",
        "interface/interface_name": "interface_name",
        "af/af_name": "af_name",
        "af/route_target/route_target_type": "route_target_type",
        "af/route_target/route_target_value": "route_target_value",
        #"_key": "source" + "_" + "vrf_name"
    }

    # Split the string by commas and spaces to get individual key-value pairs
    kv_pairs = telemetry_string.split(',')
    
    # Initialize the JSON dictionary
    json_dict = {}
    
    for pair in kv_pairs:
        # Further split each pair by '=' to separate keys and values
        if '=' in pair:
            key, value = pair.split('=', 1)
            # Remove quotes from values
            value = value.strip('"')
            # Check if the key is in the keys_to_extract set
            if key in keys_to_extract:
                json_dict[keys_to_extract[key]] = value
    
    return json_dict

# Example usage
telemetry_string = "Cisco-IOS-XR-mpls-vpn-oper:l3vpn/vrfs/vrf,host=telegraf,path=Cisco-IOS-XR-mpls-vpn-oper:l3vpn/vrfs/vrf,source=xrd01,subscription=base_metrics,vrf_name=red vrf_name_xr=\"red\",route_distinguisher=\"10.0.0.1:0\",interface/interface_name=\"Loopback9\",af/af_name=\"ipv6\",af/saf_name=\"unicast\",af/route_target/route_target_type=\"export\",af/route_target/route_target_value=\"9:9\",af/route_target/af_name=\"ipv6\",af/route_target/saf_name=\"unicast\",is_big_vrf=\"false\" 1719859230833000000"
parsed_json = parse_telemetry_string(telemetry_string)
print(json.dumps(parsed_json, indent=4))

# # convert json string to dict
# msgdict = json.loads(parsed_json)
# vrf = msgdict['fields']['vrf']
# name = msgdict['fields']['source']
# print(name, vrf)

# # generate DB ID and Key
# key = name + "_" + vrf
# id = "srv6_local_sids/" + key
# msgdict['_key'] = key
# print("id: ", id)
