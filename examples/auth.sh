#!/usr/bin/env bash

# Trigger authentication (device code flow)
cocli auth login
# => Please login at https://dev-mkiob4vl.us.auth0.com/activate?user_code=ABCD-1234

# Print refresh token (can be used to authenticate in CI or anywhere browser login is not possible)
cocli auth print_refresh_token
# => COCLI_REFRESH_TOKEN='SbS8w3Fofsh9JFp3mT2qfHEzp33YAuTkznwm8k2J-5ib7'
