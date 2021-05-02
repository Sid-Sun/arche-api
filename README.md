# Arche's REST API

## Authentication
### Sign-Up:

Method: `POST`

Path: `/v1/signup`

Body:
```json
{
    "email": "jane@example.com",
    "password": "&now:we@pluto"
}
```



### Login:

Method: `POST`

Path: `/v1/login`

Body:
```json
{
    "email": "jane@example.com",
    "password": "&now:we@pluto"
}
```

## ~~Session Management~~
### Refresh:

Method: `POST`

Path: `/v1/session/refresh`

Body:
```json
{
    "email": "jane@example.com",
    "password": "&now:we@pluto"
}
```

## Folders
### Create:
Method: `POST`

Path: `/v1/folders/create`

Body:
```json
{
    "name": "koala"
}
```

### Get:
Method: `GET`

Path: `/v1/folders/get`


### Delete:
Method: `DELETE`

Path: `/v1/folders/delete`

Body:
```json
{
    "folder_id": 1
}
```

## Notes
### GetAll:
Method: `GET`

Path: `/v1/notes/getall`

### Create:
Method: `POST`

Path: `/v1/notes/create`

Body:
```json
{
    "name": "hello",
    "data": "I am a butterfly, flying through the sky",
    "folder_id": 1
}
```

### Update:
Method: `PUT`

Path: `/v1/notes/update`

Body:
```json
{
    "folder_id": 1,
    "note_id": 7,
    "name": "squirrel",
    "data": "I am a squirrel"
}
```

### Delete:
Method: `DELETE`

Path: `/v1/notes/delete`

Body:
```json
{
    "note_id": 12
}
```