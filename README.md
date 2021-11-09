# adise1941

adise ergasia

# usage
how to use this

## register user
using curl

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
