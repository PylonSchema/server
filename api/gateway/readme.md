# Socket Protocol

### Base form
```
{
    "op": opcode (int),
    "d" : payload (object)
}
```
opcode
- 0 : MessageHeartbeat, ping-pong check websocket is alive  
- 1 : MessageAuthentication, authorize user and define socket client   
- 2 : MessageData, communication with server (like message, custom payload etc..)   
- 8 : MessageEvent, event message (like change of user or user check noitification)
- 9 : MessageError, error with several reason (like eauthentication error, websocket write length error etc.)
- 10: MessageClose, websocket close event

### Heartbeat - OPCODE 0
Heartbeat check for websocket manage
1. Sever -> Client : heartbeat check
```
{
    "op": opcodeHeartbeat,
    "d": nil
}
```
2. Client -> Server : response to heartbeat check
```
{
    "op": opcodeHeartbeat,
    "d": nil
}
```

### Authentication - OPCODE 1
Authentication is only one event since websocket is connected
If authentication is failed with serveral reason (in-valid form, not exist user etc.), websocket will be disconnected

1. Client -> Server : authentication request with valid json form
```
{
    "op": opcodeAuthentication,
    "d" : {
        "token": "user-jwt-token",
    }
}
```
2. Sever -> Client : Authentication success response
```
{
    "op": opcodeEvent,
    "d" : "" // this need change
}
```