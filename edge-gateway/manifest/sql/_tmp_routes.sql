USE shopa_edge_gateway;
SELECT route_code,method,path_pattern,auth_required,status FROM edge_proxy_route ORDER BY route_code;
