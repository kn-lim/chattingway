# chattingway

![Go](https://img.shields.io/github/go-mod/go-version/kn-lim/chattingway)
[![Go Reference](https://pkg.go.dev/badge/github.com/kn-lim/chattingway/v2.svg)](https://pkg.go.dev/github.com/kn-lim/chattingway/v2)
![GitHub Workflow Status - Release](https://img.shields.io/github/actions/workflow/status/kn-lim/chattingway/release.yaml)
![GitHub Workflow Status - Test](https://img.shields.io/github/actions/workflow/status/kn-lim/chattingway/test.yaml?label=tests)
[![codecov](https://codecov.io/gh/kn-lim/chattingway/branch/main/graph/badge.svg)](https://codecov.io/gh/kn-lim/chattingway)
![License](https://img.shields.io/github/license/kn-lim/chattingway)

Go module holds the shared code for my chat bots:

- [dreamingway-bot](https://github.com/kn-lim/dreamingway-bot) (Discord)
- ~~[slackingway-bot](https://github.com/kn-lim/slackingway-bot) (Slack)~~

This module contains packages that are for general use, interact with my homelab or for gaming-related purposes. Everything in this module is generic and each bot contains the platform-specific code.

## Packages

- `aws` - Helpers for interacting with AWS services used by the chat bots.
- `cloudflare` - Helpers for managing Cloudflare DNS records used by the chat bots.
- `counter` - Stores named counters for each guild in a DynamoDB table.
- `gamble` - Chance-based games such as coin flips and dice rolls.
- `healthcheck` - Simple liveness commands for verifying that a bot is responsive.
- `mcstatus` - Reports the status of a Minecraft Java server via the mcstatus.io API.
- `projectzomboid` - Orchestrates the lifecycle of a Project Zomboid game server.
- `rcon` - Executes commands against a server over the RCON protocol.
