# go-backend-project
It support only Add item and get item using DynamoDB.

```
API TO get users https://i5imc6al18.execute-api.us-east-2.amazonaws.com/deplopment/eventFunction
```

**Sample Payload**
```
API:- POST http://localhost:6000/event
Body {
    "ID": 7,
    "name": "Test name ",
    "description": "This is a test",
    "status": "Idle",
    "schedule": {
        "start_time": "10:20",
        "stop_time": "12:20"
    },
    "User": "xyz@gmail.com"
}
```
**Sample Response**
```
API:- GET http://localhost:6000/event/xyz@gmail.com
Response {
    "ID": 7,
    "name": "Test name ",
    "description": "This is a test",
    "status": "Idle",
    "schedule": {
        "start_time": "10:20",
        "stop_time": "12:20"
    },
    "User": "xyz@gmail.com"
}
{
    "ID": 4,
    "name": "Test name 1 ",
    "description": "This is a test 0",
    "status": "Idle",
    "schedule": {
        "start_time": "01:20",
        "stop_time": "09:20"
    },
    "User": "xyz@gmail.com"
}
```
