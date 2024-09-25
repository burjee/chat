# SERVER (MEMORY)

A chat server that stores chat room information and broadcasts messages in the server's memory. Since all data is handled in memory, the server cannot scale horizontally to improve performance.


The `counter.txt` file records the number of messages sent by the server per second, and the `web` folder contains the frontend's built assets.

### Performance

With 1 CPU and 128MB of memory, the server can handle 1000 connections and 25,000 messages per second. During connection testing, message transmission delays can be observed at peak times.

### Dev

```
go run .
```

### Run

```
docker compose up -d
```

Open the webpage on local port 8000.
