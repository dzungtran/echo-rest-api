package authz

import data.utils
import future.keywords.in

# By default, deny requests.
default allow = false

default usr_role = "guest"

default org_id = ""

default req_permission = "not_found"

roles_chart_graph[role_name] = edges {
	data.roles_chart[role_name]
	edges := {neighbor | data.roles_chart[neighbor].owner == role_name}
}

roles_chart_permissions[role_name] = access {
	data.roles_chart[role_name]
	reachable := graph.reachable(roles_chart_graph, {role_name})
	access := {item | reachable[k]; item := data.roles_chart[k].access[_]}
}

org_id = id {
	id = format_int(input.org.id, 10)
}

org_id = id {
	id = format_int(input.project.org_id, 10)
}

org_id = id {
	id = format_int(input.payload.org_id, 10)
}

usr_role = input.user.org_role[org_id]

req_permission = access {
	access = data.endpoints_acl[input.endpoint][input.method]
}

# Alway allow Super Admin
allow {
	utils.is_super_admin
}

# Start check ACL
allow {
	req_permission in roles_chart_permissions[usr_role]
	count(deny) == 0
}

# Check additional permissions with resource from input
allow {
	input.resource_id in input.resource_perms[req_permission]
	count(deny) == 0
}

# Some endpoint does not require role
allow {
	no_need_role_check[_] = {
		"endpoint": input.endpoint,
		"method": input.method,
	}

	count(deny) == 0
}
