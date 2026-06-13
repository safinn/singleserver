<p align="center">
  <img src="www/server-single-wordmark.svg" alt="single server" width="484">
</p>

<p><strong>Git push to production in &lt;5 seconds.</strong></p>

<p>One Linux box can run every app you've ever shipped. Single Server wires up Cloudflare, Tailscale, Docker, Kamal, and GitHub. From there, every <strong>git push</strong> goes live in under 5 seconds. No build queue, no per-app bills, no pile of YAML.</p>

<p>This README is for working on Single Server itself. For what it does and why, see <a href="https://singleserver.com">singleserver.com</a>. For the full manual, see the <a href="https://singleserver.com/docs/">docs</a>.</p>

## Quickstart

SSH into any fresh Linux box:

```sh
ssh root@<server_ip>
```

Run the interactive installer:

```sh
curl -fsSL https://singleserver.com/install.sh | sh
```

Point it at your repos:

```sh
singleserver add https://github.com/you/your-app
```

Every push to the configured branch deploys from then on. The [docs](https://singleserver.com/docs/) cover everything else.

## Repository layout

- `cmd/singleserverd` — entry point for the one Go binary, which runs as both the deploy daemon and the CLI.
- `internal/singleserver` — all of the product: config, deploys, provider connections, and commands.
- `www` — singleserver.com, including the docs page and `install.sh`. The site deploys as a Single Server app, and its `Dockerfile` builds the hosted binaries it serves at `/bin/`.
- `test/e2e-local-real` — the real-dependency end-to-end harness.
- `config/apps.example.yml` — an example app config.

## Build and test

The module targets Go 1.26 and has no service dependencies for development:

```sh
go build ./...
go test ./...
```

The unit tests are plain `go test` with no setup, no network, and no build tags. They also run inside the site's Docker build, so a deploy fails before it ships if a test fails.

`go run ./cmd/singleserverd help` prints the CLI. Most commands expect to run on a configured server, so for real behavior use the E2E harness or a throwaway VPS.

## End-to-end tests

The E2E harness spins up disposable Linux hosts in Docker Desktop and runs the real installer against real Tailscale, Cloudflare, and GitHub. It is serial, stateful, and needs provider credentials in a local `.env`.

Copy the template and fill it with real test credentials (the file is gitignored):

```sh
cp test/e2e-local-real/.env.example test/e2e-local-real/.env
chmod 600 test/e2e-local-real/.env
```

The `.env` holds credentials for three providers, all of which should be dedicated to testing, never production:

- **Cloudflare** — `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID`, plus a `TEST_ZONE` the tests create and delete temporary app domains in.
- **Tailscale** — an OAuth client (`TAILSCALE_OAUTH_CLIENT_ID`, `TAILSCALE_OAUTH_CLIENT_SECRET`, `TAILSCALE_TAG`) used to mint short-lived auth keys, or a `TAILSCALE_AUTHKEY` fallback.
- **GitHub** — a dedicated test App (`GITHUB_APP_ID`, `GITHUB_APP_SLUG`, `GITHUB_WEBHOOK_SECRET`, `GITHUB_APP_PRIVATE_KEY_PATH`) installed on a throwaway `GITHUB_TEST_REPO`. The harness repoints the App's webhook, so do not use the production App. Commits push with your `gh` login unless you set `GITHUB_PUSH_TOKEN`.

[test/e2e-local-real/README.md](test/e2e-local-real/README.md) explains how to obtain each, including a GitHub App bootstrap helper and the Tailscale tag policy. Then run:

```sh
test/e2e-local-real/run.sh
```

To keep secrets outside the repo, point the harness at a `.env` elsewhere:

```sh
E2E_ENV_FILE=~/secrets/singleserver-e2e.env test/e2e-local-real/run.sh
```

---

<p align="center">MIT License · <a href="https://singleserver.com">singleserver.com</a> · <a href="https://singleserver.com/docs/">Docs</a></p>
