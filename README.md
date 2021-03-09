# JingBot

A friendly greeting bot in memory of our beloved friend Jing <3. This app will post a :remote-sleepy-morning: emoji
 according to the cron schedule every day. Note that the cron schedule is defined in UTC time, so time zone conversion
 is necessary.

You need to set the following environment variables for the app to function:
```
JINGBOT_API_KEY=<slack bot API key>
JINGBOT_USER_ID=<slack bot user ID>
JINGBOT_ENVIRONMENT=<e.g dev or prod>
JINGBOT_TIME=<cron schedule string>
N_RANDOM_MINUTES=<minutes to wait>
PORT=8080
```

And deploy wherever you want!

## Running Locally

Make sure you have the latest version of docker installed, as well as golang 1.16 for more local testing.

1. run `make local-docker-image` to create the docker image
2. Setup your environment variables in a file called `.env`
3. run `make docker-compose-run` to run the docker container

See the `makefile` for other commands.



## API Usage

The app exposes the endpoints:

- `api/v1/hello` used in slack to respond to the `/jing` command with a `:remote-sleepy-morning:` 
- `api/v1/schedule` the current schedule at which the app posts a good morning message
- `api/v1/channels` the channels where the bot is currently active
    - It is highly recommended that you blocklist this to the public

You can see `hello` and `schedule` publicly here:
- https://jingbot.dev.services.atlassian.com/api/v1/hello
- https://jingbot.dev.services.atlassian.com/api/v1/schedule