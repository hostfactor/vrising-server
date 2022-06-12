# V Rising Linux Docker Server

<p>
  <a aria-label="Deploy on Host Factor" href="https://hostfactor.io/games/vrising">
    <img src="https://img.shields.io/badge/Deploy-Host%20Factor-%234f6ac6?labelColor=1b1c1d&style=for-the-badge" alt="">
  </a>
  <a aria-label="Build status" href="https://github.com/hostfactor/vrising-server/actions/workflows/build_latest.yml">
    <img src="https://img.shields.io/github/workflow/status/hostfactor/vrising-server/Build%20latest?style=for-the-badge&labelColor=1b1c1d" alt="">
  </a>
  <a aria-label="Latest Docker version" href="https://hub.docker.com/repository/docker/hostfactor/vrising-server">
    <img src="https://img.shields.io/docker/v/hostfactor/vrising-server?style=for-the-badge&labelColor=1b1c1d" alt="">
  </a>
</p>

## Basic usage

```
docker run -v /my/folder:/root/saves -p 9876:9876/udp hostfactor/vrising-server
```

Where `/my/folder` is an empty folder on your computer. If you don't add the `-v` bit, your server save will be lost
after the Docker container exits.

> You can also run a specific version by using `docker run [...] hostfactor/vrising-server:0.5.41821`

Once the server is ready you'll see a message appear that says.

> Your V Rising server is ready!

Once you see that message, go to V Rising

1. Click `Online Play`
2. Select a game mode
3. Click `Display all Servers & Settings`
4. `Direct Connect`
5. enter `127.0.0.1`

## Using existing saves

If you want to bring an existing save, simply find the folder housing your server saves

Copy that path in your folder's address bar and run the server with

```
docker run -e SAVE_NAME=world1 -v /my/save/path:/root/saves -p 9876:9876/udp hostfactor/vrising-server
```

Replace `world1` with the name of the folder in your saves folder.

The above command will allow your Docker container to access the save folder on your computer

## Configuring your server

The server is configured via the `ServerGameSettings.json` and `ServerHostSettings.json` files. You can find more about
their
structure online, but they can be used within your Docker server with the following command.

```
docker run -v /my/settings/path:/root/vrising/VRisingServer_Data/StreamingAssets/Settings -p 9876:9876/udp hostfactor/vrising-server
```

Where `/my/settings/path` is the absolute path to your settings folder which houses the `*.json` files.

### Custom args

The server also supports custom arguments that can be passed into the server command line.
See [here](https://github.com/StunlockStudios/vrising-dedicated-server-instructions#configuring-the-server) for
additional arguments.

```
docker run -v /my/folder:/root/saves -p 9876:9876/udp hostfactor/vrising-server -lan
```

The above puts the server in LAN mode allowing for offline play.

## Variables

You can configure quite a bit through passing `-e` flags into your `docker run` commands e.g.

| Name                   | Default                       | Description                                                                            |
|------------------------|-------------------------------|----------------------------------------------------------------------------------------|
| `SAVE_NAME`            | `world1`                      | The name of the folder under the `SAVE_DIR`                                            |
| `SAVE_DIR`             | `/root/saves`                 | The absolute path to the save folder.                                                  |
| `SERVER_NAME`          | `Host Factor V Rising Server` | The name of the server as it appears in-game                                           |
| `LOG_FILE`             | `/root/server.log`            | The absolute path where the log file is stored.                                        |
| `RENFIELD_SERVER_PORT` |                               | The port number to start the HTTP server. If not specified, the server is not started. |

## HTTP API

HTTP APIs are made available through `renfield` if `RENFIELD_SERVER_PORT` is set e.g.

```
docker run -v /my/folder:/root/saves -p 9876:9876/udp -p 8080:8080 -e RENFIELD_SERVER_PORT=8080 hostfactor/vrising-server
```

| Path                | Method | Description                                                  | Request schema | Response schema |
|---------------------|--------|--------------------------------------------------------------|----------------|-----------------|
| `/api/server/ready` | `GET`  | Returns 200 when the server is ready and a non-200 when not. |                |                 |
