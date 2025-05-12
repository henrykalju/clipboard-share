# clipboard-share
A desktop application to share clipboard with history between devices

## Setting it up
### Server
- `cd server`
- `cp sample.env .env`
- Change anything needed in the `.env` file
- Add a TLS proxy in front of the `server` container (unless developing locally)
- `docker compose up -d`

### Client
- `cd client`
- `wails dev --appargs "-dev"` (for local development)
- `wails build`
- built binaries are in the `build` folder (-dev flag still works for disabling TLS requirement)

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

