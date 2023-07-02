# Continue Watching API

[![build](https://github.com/learn-video/continue-watching-api/actions/workflows/build.yml/badge.svg)](https://github.com/learn-video/continue-watching-api/actions/workflows/build.yml)

Do you enjoy watching videos on demand platforms like YouTube, HBO, Prime Video, or Globoplay? If so, you may have noticed a convenient feature called "Continue Watching." This feature saves you from the hassle of remembering where you left off in a video.

In this project, we aim to develop a video player application that incorporates the "Continue Watching" functionality. The video player will interact with an API to store and retrieve the current position being watched in a video.

Once you start watching a video, the video player will periodically make requests to the API to record the current position. This ensures that even if you leave the video and come back later, the player will resume playback from the exact moment you last watched.

The API will be responsible for handling the requests from the video player and storing the positions in a reliable data store, such as Redis. When requested, the API will retrieve the last known position for a specific video and provide it to the video player, enabling seamless playback from where you left off.

Here is a sequence diagram to help you visualize how the information flows:

```mermaid
sequenceDiagram
    participant VideoPlayer
    participant API
    participant Redis

    VideoPlayer->>+API: Fetch(videoId)
    API->>+Redis: Fetch(videoId)
    Redis-->>-API: PositionResponse(position)
    alt Position exists
        API->>-VideoPlayer: PositionResponse(position)
        VideoPlayer-->>API: HTTP 200 OK
    else Position does not exist
        VideoPlayer-->>API: HTTP 404 Not Found
    end
    opt Redis query failed
        VideoPlayer-->>API: HTTP 500 Internal Server Error
    end

    loop Record position every X seconds
        VideoPlayer->>+API: Record(videoId, position)
        API->>+Redis: Record(videoId, position)
        Redis-->>-API: PositionResponse(success)
        alt Position recorded successfully
            API-->>VideoPlayer: HTTP 200 OK
        else Failed to record position
            API-->>VideoPlayer: HTTP 500 Internal Server Error
        end
    end
```
