# Single Server

Single Server is a tiny deploy daemon for running many small apps on one server.

It receives GitHub App `push` webhooks, checks a central allowlist, fetches the exact pushed SHA, and runs Kamal on the host.

## Minimal config

```yaml
apps:
  - dvassallo/sillyface-games
  - dvassallo/fullsend
```

The repo name drives the defaults:

```text
app name:  repository name
checkout:  /srv/repos/<app>
deploy:    kamal setup -q on first deploy, kamal redeploy -q after that
branch:    repository default branch from the webhook payload
```

Use an object only when an app needs overrides:

```yaml
apps:
  - repo: dvassallo/fullsend
    branch: master
    healthcheck: https://fullsend.game/up
```

## Host secrets

Secrets live on the server, not in app repositories.

```text
/etc/singleserver/apps.yml
/etc/singleserver/singleserver.env
/etc/singleserver/github-app.json
/etc/singleserver/github-app.private-key.pem
```

Required environment:

```sh
SINGLESERVER_CONFIG=/etc/singleserver/apps.yml
SINGLESERVER_PORT=8787
SINGLESERVER_PUBLIC_URL=https://hooks.singleserver.com
SINGLESERVER_WEBHOOK_SECRET=...
SINGLESERVER_GITHUB_TOKEN=...
```

With `SINGLESERVER_GITHUB_TOKEN` set, the daemon creates or updates matching repo webhooks for every app in `apps.yml` on startup and once per minute.

## GitHub App

The GitHub App needs:

- repository contents: read
- commit statuses: write
- event subscription: push

Install it with access to all repositories, then let `apps.yml` be the deployment allowlist.

The daemon includes a one-time setup page:

```text
https://hooks.singleserver.com/setup/github-app?token=<setup-token>
```

That page creates the GitHub App from a manifest, exchanges GitHub's callback code, and stores the app credentials under `/etc/singleserver`.
