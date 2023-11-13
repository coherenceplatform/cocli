#!/usr/bin/env bash

# Commit sha
commit_sha="aecf6d23c94ffa81fe6ea60fdb120c30841572e4"

# Branch name for a branch that does not exist yet
branch_name="my-feature-123"

# Get app id
app_id=$(cocli apps list | jq '.[] | select(.title=="Next App AWS") | .id')

# Create feature
cocli features create $branch_name --app_id $app_id --commit_sha $commit_sha
