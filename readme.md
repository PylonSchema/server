# Pylon Server

### Requirements
MySql  
Redis  
Go  

### Make config file

ref ./example.toml file and change filename to conf.toml  

### Set Scylla database with docker
Checkout more options - [Scylla Docker Image](https://hub.docker.com/r/scylladb/scylla/)  
```powershell
docker run --name scylla -d -p 22:22 -p 7000:7000 -p 7001:7001 -p 9042:9042 -p 9160:9160 -p 9180:9180 scylladb/scylla --smp 1
```
Check security list - [Scylla Security](https://opensource.docs.scylladb.com/stable/operating-scylla/security/security-checklist.html)  

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
