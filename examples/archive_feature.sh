#!/usr/bin/env bash

# Branch name that feature is based on
branch_name="my-feature-123"

# Get app id
app_id=$(cocli apps list | jq '.[] | select(.title=="Next App AWS") | .id')

# Get feature id
feature_id=$(cocli features list --app_id $app_id | jq ".data[] | select(.branch_name==\"$branch_name\") | .id")

# Create feature
cocli features archive --feature_id $feature_id
