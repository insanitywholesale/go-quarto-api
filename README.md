# adise1941

adise ergasia

# usage
how to use with curl

## register user
run:
```bash
curl -X POST\
	-H "Content-Type: application/json"\
	-d '{"username": "inherently", "password": "verybigsecret"}'\
	localhost:8000/user/register
```

returns:
```json
{"username":"inherently","user_id":"G8boeMc7g"}
```

or if you have `jq` installed you can extract the `user_id` value using:
```bash
curl -X POST\
	-H "Content-Type: application/json"\
	-d '{"username": "inherently", "password": "verybigsecret"}'\
	localhost:8000/user/register | jq '.user_id'
```

## create game
run:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"user_id": "G8boeMc7g"}' localhost:8000/game/new
```

returns:
```json
{"game_id":"NvFtm757g","players":[{"username":"","user_id":"G8boeMc7g"}],"activity_status":true,"game_state":{}}
```
