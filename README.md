# ModStats
This is a website which displays information about tModloader mods.
![image](https://user-images.githubusercontent.com/26361108/205339062-78e3a42f-1e96-4abf-b224-2c113f3e71c9.png)

The site can be found under https://modstats.le0n.dev/

## Features
- Mod Search (internal- and display name)
- Author Search (through steamid64 and steam name)
- Mod List
- Random Mod
- 1.3 and 1.4 support

### How to run locally

Requirements:

 - Go

 1. Clone the repo `git clone https://github.com/NotLe0n/ModStats.git`
 2. Navigate to it `cd ModStats`
 3. create `config.json`, for the default config write: `{}`
 4. Run `go run ./server`
 5. Browse `http://localhost:8080`

### Configuration
This is how the default config looks like:
```json
{
	"API-URL": "https://tmlapis.le0n.dev",
	"port": 8080,
	"useHTTPS": false,
	"certPath": "",
	"keyPath": ""
}
```
