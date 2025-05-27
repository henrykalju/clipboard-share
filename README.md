# clipboard-share
A desktop application to share clipboard with history between devices

## Public server
- I have set up the server side on `clipboard-share.northeurope.cloudapp.azure.com`.
- Be careful when choosing server, because while the traffic is protected with TLS, the data is not encrypted in the database.

## Usage (in case of self hosted server, replace url with your own)
- Go to [https://clipboard-share.northeurope.cloudapp.azure.com/](https://clipboard-share.northeurope.cloudapp.azure.com/)
- Create an account
- Download the executable from [GitHub](https://github.com/henrykalju/clipboard-share/releases)
- Run the executable
- In the settings, use the created account and `clipboard-share.northeurope.cloudapp.azure.com` as url (no scheme in the beginning and no `/` in the end)
  
[![Demo Video](https://img.youtube.com/vi/3YQ6cWbY0Sg/0.jpg)](https://www.youtube.com/watch?v=3YQ6cWbY0Sg)

## How it works
### Server
- Go API
    - Serves index.html website
    - Has endpoints to register user, check username+password, save and get clipboard items
    - Saves users and data to Postgres db
- index.html
    - Allows registering new users
    - Allows user to access history from unsupported devices

### Client
Supported devices: Windows, X11
- Wails
    - Binds Svelte frontend and Go backend
- Svelte
    - GUI to see history and select item to be put on the clipbaord
    - Configure application (username, password, server url)
- Go
    - Clipboard (windows package for Windows or C, Xlib and XFixes for Linux)
        - Listens to copying and puts all data into a channel
        - Can write any data to the clipboard
    - When clipboard gives new data, save it and refresh GUI list
    - When GUI selects an item, write it to clipboard

## Self hosting
### Server
- `curl -fsSL https://get.docker.com -o get-docker.sh` (or install docker some other way)
- `git clone https://github.com/henrykalju/clipboard-share.git`
- `cd clipboard-share`
- `cd server`
- `cp sample.env .env`
- Change anything needed in the `.env` file (e.g. postgres login details)
- Set up TLS proxy (a simple one is set up using the `docker-compose.yaml` and `nginx.conf`)
    - `cp sample.nginx.conf nginx.conf`
    - Replace `[domain-name]` in `nginx.conf` file
    - `apt install certbot`
    - `certbot certonly --standalone -d [domain-name]`
- `docker compose up -d`

### Client
- `cd client`
- `wails dev --appargs "-dev"` (for local development)
- `wails build` (for building)
    - built binaries are in the `build/bin` folder (-dev flag still works for disabling TLS requirement)

