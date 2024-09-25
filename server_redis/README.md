# SERVER (REDIS)

The chat server uses Redis to store and broadcast messages, allowing for horizontal scaling to improve server performance. In the Docker Compose setup, two backend servers are configured, and connections are load-balanced through Nginx. Overall, the performance is better compared to handling everything in memory.


The `counter1.txt` and `counter2.txt` files record the number of messages sent per second by each of the two servers, and the `web` folder contains the frontend's built assets.

### Performance

Each server, with 1 CPU and 128MB of memory, can handle 1000 connections and 25,000 messages per second. During connection testing, the latency is significantly lower compared to the other server type.

### Run

```
docker compose up -d
```

Open the webpage on local port 8000.
