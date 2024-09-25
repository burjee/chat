# SIMPLE CHAT

A simple chat server written in Go. The server is divided into two modes: one that stores and broadcasts messages in memory, and another that uses Redis for storage and broadcasting. Additionally, there is a frontend page and connection testing.


Under the condition of 1 CPU and 128MB of memory, the server can handle 1000 concurrent connections, sending 25,000 messages per second.

### Folder

``` bash
├── server_memory # Server (using Memory)
├── server_redis  # Server (using Redis)
├── web           # Frontend page
└── ws-test       # Connection testing