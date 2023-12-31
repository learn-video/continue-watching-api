# Continue Watching API

Please use the issues section to discuss and share your solutions. Remember to tag the issue with the "challenge" label to make it easier to filter and locate solutions in the future.

<p align="center" width="100%">
    <img src="bbb.png" alt="big buck bunny" />
</p>

---

[![build](https://github.com/learn-video/continue-watching-api/actions/workflows/build.yml/badge.svg)](https://github.com/learn-video/continue-watching-api/actions/workflows/build.yml)

Do you enjoy watching videos online? If so, you may have noticed a convenient feature called "Continue Watching." This feature saves you from the hassle of remembering where you left off in a video.

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

## Demo

To experience the "Continue Watching" feature in action, we provide a demo page that allows you to experiment with the functionality firsthand.

Before running the project, make sure you have the following tools installed:

* Go
* Redis

This project utilizes a development environment tool called _devbox_, which simplifies the setup process. It provides a pre-configured environment, including Redis, to ensure a smooth experience.

To start the project, follow these steps:

1. Launch the `devbox shell` by running the appropriate command.
2. Within the _devbox shell_, execute `go run main.go` to start the server. This will also expose the index page containing the video player.

Once the server is running, you can begin testing the "Continue Watching" feature:

1. Start playing the video and watch it for at least 30 seconds. Close this tab.
2. Open a private browser tab or window.
3. Navigate to the index page again in the private tab or window.
4. The video will resume playback from a position as close as possible to where you left off, allowing you to seamlessly continue watching.

By following these steps, you can observe how the project effectively stores and retrieves the last watched position, providing a convenient viewing experience for users. Feel free to explore different scenarios, pause the video, seek to the future to fully grasp the capabilities of the "Continue Watching" feature.

## Challenges to the reader

The most basic scenarios are described. How would you evolve the "Continue Watching API" with the following challenges in mind?

* **Concurrent requests for the same video's position?** Imagine multiple instances of the video player making simultaneous requests to update the position for a specific video. How would you ensure data consistency and prevent conflicts when multiple requests are received concurrently?

* **Scaling the system to handle a large number of users and videos?** As the number of users and videos increases, the system's scalability becomes crucial. How would you design and implement a scalable architecture that can handle a high volume of requests, store and retrieve positions efficiently, and ensure optimal performance across a growing user base?

* **Handling API Downtime** Imagine a scenario where the API goes offline for an extended period. How would you ensure that no position events are lost during this downtime?

* **Dealing with Bad API Responses** What steps would you take if the API starts responding with error status codes like HTTP 5xx? How would you handle these situations to maintain data integrity and ensure the reliability of the "Continue Watching" feature?

* **Optimal Expiring Time in Redis** Determine an appropriate expiration time for the position data stored in Redis. Consider factors such as the frequency of user activity, the maximum duration between user visits, and the balance between storage efficiency and ensuring the freshness of the position data. Find a suitable expiration time that minimizes storage overhead while still providing an accurate representation of the user's most recent position.

* **Preserving Values in Expired Redis Keys** As Redis is a volatile storage caching mechanism, keys can expire based on the configured expiration time. How would you handle situations where Redis keys containing position data have expired?
