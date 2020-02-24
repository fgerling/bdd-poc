#! /usr/bin/env sh
fulfilled="true"

if ! command -v SUSEConnect >/dev/null
then
	echo "SUSEConnect not found" 1>&2
	exit 2
fi

if ! (SUSEConnect --status | jq -er '.[] | select(.identifier=="caasp") | .subscription_status' | grep -q "ACTIVE")
then
	echo "SUSE CaaS Platform extension subscription is not ACTIVE" 1>&2
	fulfilled=false
fi
if ! (SUSEConnect --status | jq -er '.[] | select(.identifier=="sle-module-containers") | .status' | grep -q "Registered")
then
	echo "sle-module-containers is not 'Registered'"
	fulfilled=false
fi

if [ "$fulfilled" == "false" ]
then
	exit 1
fi
