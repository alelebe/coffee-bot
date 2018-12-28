# fx-coffee-bot

## How to build

git clone https://github.com/alelebe/fx-coffee-bot

*Create new file '.env'in working directory with the following content:
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
 
