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

## Authentication
For authentication you just need to get a [personal access token](https://docs.withcoherence.com) from coherence.
cocli will expect the token to be set as `COHERENCE_ACCESS_TOKEN`
