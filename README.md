# SIMPLE CHAT

A simple chat server written in Go. Under the condition of 1 CPU and 128MB of memory, the server can handle 1000 concurrent connections, sending 25,000 messages per second.


In addition to the frontend page and connection testing, the project includes two types of servers: one that stores and broadcasts messages in memory, and another that uses Redis for storage and broadcasting.

### Folder

``` bash
├── server_memory # Server (using Memory)
├── server_redis  # Server (using Redis)
├── web           # Frontend page
└── ws-test       # Connection testing