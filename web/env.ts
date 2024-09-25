const dev = process.env.NODE_ENV === "development";

export const WSURL = dev ? "ws://localhost:8000/ws" : "/ws";