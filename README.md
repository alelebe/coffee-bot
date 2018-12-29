# fx-coffee-bot

First of all, it makes sense to start with https://devcenter.heroku.com/articles/getting-started-with-go

## Prerequisit

Install Heroku CLI, login and link it with your application on Heroku.
https://devcenter.heroku.com/articles/getting-started-with-go#set-up


### GitHub integration with Heroku

With GitHub integration, every push to target branch in git repo on GitHub will trigger automatic deployment to the app on Heroku as described https://devcenter.heroku.com/articles/github-integration


## How to build locally

git clone https://github.com/alelebe/fx-coffee-bot

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
 

## How to run locally
```
	$ go install
	$ heroku local
	$ heroku local -e .env.test
	$ hrekou ci:debug
```

## How to stop the app on Heroku
```
	$ heroku ps:scale web=0
```

## How to start the app on Heroku
```
	$ git push heroku master
	$ heroku ps:scale web=1
	$ heroku open
```

## How to check logs on Heroku
```
	$ heroku logs --tail
```

## How to change ENV variable on Heroku
```
	$ heroku config:set GIN_MODE=release
```

## How to debug the app locally
    Follow steps to install debugger for Golang: https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code
	Install Delve:
```
	$ xcode-select --install
	$ go get -u github.com/derekparker/delve/cmd/dlv
```
