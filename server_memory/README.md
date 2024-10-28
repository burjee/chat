# SERVER (MEMORY)

A chat server that stores chat room information and broadcasts messages in the server's memory. Since all data is handled in memory, the server cannot scale horizontally to improve performance.


The web folder contains the frontend's built assets.

## Performance

With 1 CPU and 128MB of memory, the server can handle 1000 connections and 23,000 messages per second. During connection testing, message transmission delays can be observed at peak times.

## Dev

```bash
go run .
```

## Run

Open the webpage on local port 8000.

```bash
docker compose up -d --build
```

## Message Metrics

You can view the logs of the server, which display the average messages per second for the past one minute, five minutes, fifteen minutes, and overall.

```bash
docker logs -f --tail=0 simple-chat-memory
```

## Delete

```bash
docker compose down
```
