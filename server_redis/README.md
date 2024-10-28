# SERVER (REDIS)

The chat server uses Redis to store and broadcast messages, allowing for horizontal scaling to improve server performance. In the Docker Compose setup, two backend servers are configured, and connections are load-balanced through Nginx. Overall, the performance is better compared to handling everything in memory.


The web folder contains the frontend's built assets.

## Performance

Each server, with 1 CPU and 128MB of memory, can handle 1000 connections and 23,000 messages per second. During connection testing, the latency is significantly lower compared to the memory server type.

## Run

Open the webpage on local port 8000.

```bash
docker compose up -d --build
```

## Message Metrics

You can view the logs of the server, which display the average messages per second for the past one minute, five minutes, fifteen minutes, and overall.

```bash
docker logs -f --tail=0 simple-chat-redis-1
```

## Delete

```bash
docker compose down
```
