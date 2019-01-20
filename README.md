# fx-coffee-bot

First of all, it makes sense to start with:
    [<https://devcenter.heroku.com/articles/getting-started-with-go]>

## Prerequisites

Install Heroku CLI as described here:
    [<https://devcenter.heroku.com/articles/heroku-cli]>

The full reference for CLI commands is available here:
    [<https://devcenter.heroku.com/articles/heroku-cli-commands]>

### Login into your Heroku account

```bash
heroku login
heroku login -i
```

### Create new Heroku app

Unless your application was already created (see next section), please follow the steps to create new app on Heroku:
    [<https://devcenter.heroku.com/articles/creating-apps]>

### Link your working folder with your Heroku app

```bash
heroku git:remote -a fx-coffee-bot   *set git remote heroku to [https://git.heroku.com/fx-coffee-bot.git]*
heroku info
```

Basically, remote Heroku app must be configured as remote git repo. Therefore, an alternative and straighforward way to do that is:

```bash
git remote add heroku git@heroku.com:fx-coffee-bot.git
```

### GitHub integration with Heroku [optional]

With GitHub integration, every push to target branch in git repo on GitHub will trigger automatic deployment to the app on Heroku as described here:
    [<https://devcenter.heroku.com/articles/github-integration]>

### Speak to BotFather

Setup Bot commands:

```command
/setcommands
coffee - Place request for a coffee
collect - Collect requests and buy beverages
help - How to speak to the bot
watch - Add/remove yourself from the watchers
```

## HOWTO(s)

### How to build locally

Checkout source code:

```bash
git clone https://github.com/alelebe/fx-coffee-bot
```

*Create new file '.env'in working directory with the following content:*

```prop
PORT=8080
ENV=local
GIN_MODE=debug
DEBUG_BOT=0
BASE_URL=https://fx-coffee-bot.herokuapp.com/
NEWS_TOKEN=<...>
COFFEE_TOKEN=<...>
```

where:

- PORT is automatically defined by Heroku to listen on in your HTTP server
- ENV : {local | heroku}
- GIN_MODE : {debug | release }
- DEBUG_BOT : {0 | 1}
- BASE_URL : heroku app URL (not used in local mode
- ...TOKEN : telegram Bot token

### How to run locally

```bash
go install
heroku local
heroku local -e .env.test
hrekou ci:debug
```

### How to stop the app on Heroku

```bash
heroku ps:scale web=0
heroku ps   *//List the dynos for an app*
```

### How to start the app on Heroku

```bash
git push heroku master
heroku ps:scale web=1
heroku ps   *//List the dynos for an app
```

### How to check logs on Heroku

```bash
heroku logs --tail
```

### How to change ENV variable on Heroku

```bash
heroku config:set GIN_MODE=release
```

### How to debug the app locally

Follow steps to install debugger for Golang:[<https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code]>

#### Install Delve

```bash
xcode-select --install
go get -u github.com/derekparker/delve/cmd/dlv
```

Create .vscode/launch.json with standard configuration:

```json
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "env": {
                "ENV":"local",
                "DEBUG_BOT":"1",
                "NEWS_TOKEN":"<your token>",
                "COFFEE_TOKEN":"<your token>"
            },
            "args": []
        }
```

### How to clear build cache for your app on Heroku

[<https://help.heroku.com/18PI5RSY/how-do-i-clear-the-build-cache]>

```bash
heroku plugins:install heroku-repo
heroku repo:purge_cache -a fx-coffee-bot
git commit --allow-empty -m "Purge cache"
git push heroku master
```

## Memcache

Based on Heroku article: [<https://devcenter.heroku.com/articles/gin-memcache]>

### Verify Heroku account

Please verify your account to install this add-on plan (please enter a credit card) For more information, see
 ▸    [<https://devcenter.heroku.com/categories/billing]>
 ▸    Verify now at [<https://heroku.com/verify]>

### Add MemCachier Addon to your application

```bash
heroku addons:create memcachier:dev
govendor fetch github.com/memcachier/mc
```

### Install MemCachier locally

```bash
$ brew install memcached
==> memcached
```

To have launchd start memcached now and restart at login:

```bash
brew services start memcached
```

Or, if you don't want/need a background service you can just run:

```bash
/usr/local/opt/memcached/bin/memcached
```

## Logging Addon: Timber.io

[<https://elements.heroku.com/addons/timber-logging]>

### Add Timber.io to heroku app

```bash
heroku addons:create timber-logging:free
```
