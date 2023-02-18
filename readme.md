# Pylon Server

### Requirements
MySql  
Redis  
Go  

### Make config file

ref ./example.toml file and change filename to conf.toml  

### Run Server
```
go run main.go
```


### Hot Reloading
```
go install github.com/zzwx/fresh@latest
fresh
```
or hot reloading with service state check (for Windows)
```
go install github.com/zzwx/fresh@latest
./runner.ps1
```
run as administrator  
no hot realoding with service stat check for linux ( TODO )

### Supported OAuth2
- Github  
