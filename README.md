# go-sro-agent-server

![go-version](https://img.shields.io/github/go-mod/go-version/ferdoran/go-sro-agent-server)
![version](https://img.shields.io/github/v/tag/ferdoran/go-sro-agent-server?label=version)
![build-status](https://img.shields.io/github/workflow/status/ferdoran/go-sro-agent-server/Build%20and%20Publish%20Docker%20Image)
![last-commit](https://img.shields.io/github/last-commit/ferdoran/go-sro-agent-server)

An agent server implementation of the game Silkroad Online ,
based on [go-sro-framework](https://github.com/ferdoran/go-sro-framework) 
and [go-sro-fileutils](https://github.com/ferdoran/go-sro-fileutils),
written in Golang.

It was developed using vSRO 1.88 files.
Using different versions might result in errors and bugs.

## Config
> Environment variable values are prioritized over config file values

Configuration can be done by config file or env variables.
There is some default configuration which you can check here:

- [`go-sro-framework/config/config.go`](https://github.com/ferdoran/go-sro-framework/tree/master/config/config.go)
- [`config/agent.go`](https://github.com/ferdoran/go-sro-agent-server/tree/master/config/agent.go)

## Features

- Character Lobby
    - Character Creation
    - Character Deletion
    - Joining Game
- Movement
  - Currently, just _click-walking_ is possible
  - Terrain collision detection working
  - Object collision detection working (almost perfect)
- Spawn System
  - Automatic spawning of objects/players/monsters/NPCs in range
  - Automatic despawning of objects/players/monsters/NPCs out of range
  - Position updates of objects/players/monsters/NPCs
- Chat
  - Custom Chat Commands (Currently just for accounts with GM status)
  - All Chat
  - Party Chat
  - Notices
  - Whisper
- Party
  - Party invitation
  - Party matching system registration
  - Party matching system application (join request)
  - Party matching system application handling
- Inventory
  - Equipping items
  - Unequipping items
  - Moving items in inventory
- Stall
  - Stall creation
  - Stall item registration
  - Stall item deletion
  - Stall name changing
  - Stall chat

## Acknowledgement

As the development was not a single person's effort,
I want to thank [DaxterSoul](https://www.elitepvpers.com/forum/members/1084164-daxtersoul.html)
for sharing his wide knowledge on the game and its peculiarities.

Without his packet and file structure documentation this would not have been possible.

## Additional Projects

As this is just a framework, there are also projects taking this framework into use:

- [go-sro-fileutils](https://github.com/ferdoran/go-sro-fileutils)
- [go-sro-framework](https://github.com/ferdoran/go-sro-framework)
- [go-sro-gateway-server](https://github.com/ferdoran/go-sro-gateway-server)

## Contribution

If you want to engage in the development, you are free to so.
Simply fork this project and submit your changes via Pull Requests.

There is some documentation around the game located under `docs`,
however it does not give any guideline.
It is more like a collection of information I gathered over time.
Some of it may still be in German, so take my apologies here for not translating them yet.

Providing a usable documentation is still a big item on my TODO list for this project.
Any help on this is highly appreciated.

If you have any more questions, feel free to contact me.

## Licensing

go-sro-agent-server is licensed under the DON'T BE A DICK PUBLIC LICENSE.
See [LICENSE](LICENSE) for the full license text