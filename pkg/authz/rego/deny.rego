package authz

import data.deny

default deny = []

deny = deny.deny_user_endpoint
