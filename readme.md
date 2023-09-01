# Pylon Server
Server for chat app [Pylon](https://github.com/PylonSchema)  

### Requirements
<div>
<img src="https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=Go&logoColor=ffffff"/>
<img src="https://img.shields.io/badge/MySQL-4479A1?style=flat-square&logo=MySQL&logoColor=ffffff"/>
<img src="https://img.shields.io/badge/PowerShell-5391FE?style=flat-square&logo=PowerShell&logoColor=ffffff"/>
<img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=Docker&logoColor=ffffff"/>
<img src="https://img.shields.io/badge/ScyllaDB-6CD5E7?style=flat-squarse&logo=ScyllaDB&logoColor=ffffff"/>
</div>  


## How to start
Download [MySQL](https://www.mysql.com/downloads/)  
Download [Redis](https://redis.io/download/) - If OS is window, you can download in [here](https://github.com/microsoftarchive/redis/releases) or consider Docker  

### Make config file
ref [example.toml](./example.toml) file and change file name to conf.toml

### Set Scylla database with docker
```powershell
docker run --name scylla -d -p 22:22 -p 7000:7000 -p 7001:7001 -p 9042:9042 -p 9160:9160 -p 9180:9180 scylladb/scylla --smp 1
# by default port, development setting (single core), no specific volumn
```
Checkout more options - [Scylla Docker Image](https://hub.docker.com/r/scylladb/scylla/)  
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

## Develop guide
Use [initdb.ps1](./initdb.ps1) for remove all data from db  
Use [ruuner.ps1](./runner.ps1) for auto start service and hot reloading (fresh installation required)  