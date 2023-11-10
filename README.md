# cocli

## Installation

1. Go to [releases](https://github.com/coherenceplatform/cocli/releases) and download the binary for your os/architecture (e.g. linux/arm64)

2. Move the file into a location that is in your `$PATH` (e.g. /usr/local/bin). We recommend renaming the file `cocli`. e.g.
```bash
mv ~/Downloads/cocli-linux-arm64 /usr/local/bin/cocli
```

3. You'll likely need to update permissions to make the file executable:
```bash
chmod +x /usr/local/bin/cocli
```

4. Test that cocli works:
```bash
cocli -h
```

N.B. If you get an error trying to use the cli on a mac (along the lines of e.g. "...can’t be opened because apple cannot check it for malicious software.") then [follow these steps](https://support.apple.com/guide/mac-help/apple-cant-check-app-for-malicious-software-mchleab3a043/mac) to grant an exception for cocli.

## Development

By default, cocli is configured for production use. During development you should `export COHERENCE_ENVIRONMENT=review`.
To change your target environment update the `CoherenceDomain` in the `devConfig`. It can be found in `pkg/cocli/cocli.go`.

If you are targeting a workspace, the cli will need to use a slightly different api url path. This will be handled automatically by cocli.

### Versioning

To update the cocli version it needs to be set in 2 places:
- update the `cliVersion` constant in `pkg/cocli/cocli.go`
- update the version in `cli_version.txt`

In addition to the above, the cli api version should be updated in the control-plane metadata response.