# fx-coffee-bot

First of all, it makes sense to start with https://devcenter.heroku.com/articles/getting-started-with-go


## Prerequisites

Install Heroku CLI as described here:
    https://devcenter.heroku.com/articles/heroku-cli

The full reference for CLI commands is available here:
    https://devcenter.heroku.com/articles/heroku-cli-commands

### Login into your Heroku account
```
	$ heroku login
	$ heroku login -i
```

### Create new Heroku app

Unless your application was already created (see next section), please follow the steps to create new app on Heroku:
    https://devcenter.heroku.com/articles/creating-apps

### Link your working folder with your Heroku app
```
    $ heroku git:remote -a fx-coffee-bot	*set git remote heroku to https://git.heroku.com/fx-coffee-bot.git*
    $ heroku info
```
Basically, remote Heroku app must be configured as remote git repo. Therefore, an alternative and straighforward way to do that is:
```
    $ git remote add heroku git@heroku.com:fx-coffee-bot.git
```

### GitHub integration with Heroku [optional]

With GitHub integration, every push to target branch in git repo on GitHub will trigger automatic deployment to the app on Heroku as described here:
    https://devcenter.heroku.com/articles/github-integration


### Speak to BotFather
Setup Bot commands:
```
	/setcommands
	coffee - Place request for a hot beverage
	collect - Collect all requests to purchase coffee
	help - description of the process interaction with bot
```

## HOWTO(s)

### How to build locally
Checkout source code:
```
	git clone https://github.com/alelebe/fx-coffee-bot
```

*Create new file '.env'in working directory with the following content:*
```property file
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
```
	$ go install
	$ heroku local
	$ heroku local -e .env.test
	$ hrekou ci:debug
```

### How to stop the app on Heroku
```
	$ heroku ps:scale web=0
	$ heroku ps 	*//List the dynos for an app*
```

### How to start the app on Heroku
```
	$ git push heroku master
	$ heroku ps:scale web=1
	$ heroku ps 	*//List the dynos for an app
```

### How to check logs on Heroku
```
	$ heroku logs --tail
```

### How to change ENV variable on Heroku
```
	$ heroku config:set GIN_MODE=release
```

### How to debug the app locally
Follow steps to install debugger for Golang:
	https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code
Install Delve:
```
	$ xcode-select --install
	$ go get -u github.com/derekparker/delve/cmd/dlv
```

### How to clear build cache for your app on Heroku
	https://help.heroku.com/18PI5RSY/how-do-i-clear-the-build-cache
```
	$ heroku plugins:install heroku-repo
	$ heroku repo:purge_cache -a fx-coffee-bot
	$ git commit --allow-empty -m "Purge cache"
	$ git push heroku master
```
