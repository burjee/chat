import { WSURL } from "@/env";

export enum EventType {
    RoomList = "ROOM_LIST",
    UserList = "USER_LIST",
    Join = "JOIN",
    Leave = "LEAVE",
    Message = "MESSAGE",
    Response = "RESPONSE",
}

export interface WSEvent {
    type: string;
    content: Content;
    nonce: string;
    rooms?: string[];
    users?: string[];
    error: string;
}

export interface Content {
    from: string;
    message?: Message;
    datetime?: string;
}

export interface Message {
    room: string;
    text: string;
}

class Connection {
    socket!: WebSocket;

    onOpen: (event: globalThis.Event) => void = () => { };
    onClose: (event: CloseEvent) => void = () => { };
    onRoomList: (e: WSEvent) => void = () => { };
    onUserList: (e: WSEvent) => void = () => { };
    onJoin: (e: WSEvent) => void = () => { };
    onLeave: (e: WSEvent) => void = () => { };
    onMessage: (e: WSEvent) => void = () => { };
    onResponse: (e: WSEvent) => void = () => { };
    onError: (event: globalThis.Event) => void = () => { };

    constructor() {
        this.open();
    }

    open() {
        this.socket = new WebSocket(WSURL);
        this.socket.onerror = (event) => this.onError(event);
        this.socket.onopen = (event) => this.onOpen(event);
        this.socket.onclose = (event) => this.onClose(event);
        this.socket.onmessage = (event: MessageEvent<any>) => {
            let e: WSEvent = JSON.parse(event.data);
            switch (e.type) {
                case EventType.RoomList: { this.onRoomList(e); break; }
                case EventType.UserList: { this.onUserList(e); break; }
                case EventType.Join: { this.onJoin(e); break; }
                case EventType.Leave: { this.onLeave(e); break; }
                case EventType.Message: { this.onMessage(e); break; }
                case EventType.Response: { this.onResponse(e); break; }
            }
        };
    }

    send(message: string) {
        this.socket.send(message);
    }

    close() {
        this.socket.close();
    }
}

export default Connection;