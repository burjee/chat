import { useState, useEffect, createContext, ReactNode, useContext, useReducer } from "react";
import Connection from "../libs/socket";
import { EventType, WSEvent } from "../libs/socket";
import { v4 as uuidv4 } from "uuid";

type IRooms = string[];

interface IUsers {
    [key: string]: "joined" | "left";
};

interface IMessage {
    [key: string]: WSEvent[];
}

interface IResponse {
    [key: string]: string;
}

interface IWebsocketContext {
    connection: Connection;
    isConnected: boolean;
    isJoined: boolean;
    setIsJoined: React.Dispatch<React.SetStateAction<boolean>>;
    room: string;
    setRoom: React.Dispatch<React.SetStateAction<string>>;
    rooms: IRooms;
    users: IUsers;
    messages: IMessage;
    responses: IResponse;
    dispatch: (value: Actions) => void
}

export enum ActionType {
    RoomList,
    UserList,
    Join,
    Leave,
    Message,
    Response,
    DeleteResponse,
}

interface RoomsAction {
    actionType: ActionType.RoomList;
    payload: IRooms;
}

interface UsersAction {
    actionType: ActionType.UserList;
    payload: string[];
}

interface JoinAction {
    actionType: ActionType.Join;
    payload: WSEvent;
}

interface LeaveAction {
    actionType: ActionType.Leave;
    payload: WSEvent;
}

interface MessageAction {
    actionType: ActionType.Message;
    payload: WSEvent;
}

interface ResponseAction {
    actionType: ActionType.Response;
    payload: WSEvent;
}

interface DeleteResponseAction {
    actionType: ActionType.DeleteResponse;
    payload: string;
}

type Actions = RoomsAction | UsersAction | JoinAction | LeaveAction | MessageAction | ResponseAction | DeleteResponseAction;

interface State {
    rooms: IRooms;
    users: IUsers;
    messages: IMessage;
    responses: IResponse;
}

function reducer(state: State, action: Actions) {
    switch (action.actionType) {

        case ActionType.RoomList:
            action.payload.forEach(room => {
                if (state.messages[room] === undefined) {
                    state.messages[room] = [];
                }
            })

            return {
                ...state,
                rooms: action.payload
            };

        case ActionType.UserList:
            let users: IUsers = {};
            action.payload.forEach(user => {
                users[user] = "joined";
            });
            return { ...state, users };

        case ActionType.Join: {
            limitMessageNumber(state.messages);
            let user = action.payload.content.from;
            if (state.users[user] === "left") {
                delete state.users[user];
            } else {
                let joinedMessage = { type: EventType.Join, content: action.payload.content, nonce: action.payload.nonce, error: "" };
                state.users[user] = "joined";
                state.rooms.forEach(room => {
                    state.messages[room] = [...state.messages[room], joinedMessage];
                });
            }
            return { ...state, messages: { ...state.messages } };
        }

        case ActionType.Leave: {
            limitMessageNumber(state.messages);
            let user = action.payload.content.from;
            if (state.users[user] === "joined") {
                delete state.users[user];
                let leftMessage = { type: EventType.Leave, content: action.payload.content, nonce: uuidv4(), error: "" };
                state.rooms.forEach(room => {
                    state.messages[room] = [...state.messages[room], leftMessage];
                });
            } else {
                state.users[user] = "left";
            }
            return { ...state, messages: { ...state.messages } };
        }

        case ActionType.Message:
            limitMessageNumber(state.messages);
            let room = action.payload.content.message!.room;
            return {
                ...state,
                messages: {
                    ...state.messages,
                    [room]: [...state.messages[room], action.payload]
                },
            };

        case ActionType.Response:
            return {
                ...state,
                responses: {
                    ...state.responses,
                    [action.payload.nonce]: action.payload.error
                }
            };

        case ActionType.DeleteResponse:
            delete state.responses[action.payload];
            return { ...state, responses: { ...state.responses } };

        default:
            return state
    }
}

function limitMessageNumber(messages: IMessage) {
    Object.keys(messages).forEach(room => {
        if (messages[room].length > 99) {
            messages[room].shift();
        }
    });
}

const WebSocketContextProvider = ({ children }: { children: ReactNode | ReactNode[] }) => {
    const [connection, setConnection] = useState({} as Connection);
    const [trigger, setTrigger] = useState(0);
    const [isConnected, setIsConnected] = useState(false);
    const [isJoined, setIsJoined] = useState(false);
    const [room, setRoom] = useState("");
    const [state, dispatch] = useReducer(reducer, { rooms: [], users: {}, messages: {}, responses: {} });

    useEffect(() => {
        let connection = new Connection();

        connection.onOpen = () => {
            setIsConnected(true);
        };

        connection.onClose = () => {
            setIsConnected(false);
            setIsJoined(false);
            setTimeout(() => {
                setTrigger(trigger + 1);
            }, 3000);
        };

        connection.onError = () => {
            connection.close();
            setIsConnected(false);
        };

        setConnection(connection);

        return () => { connection.close() }
    }, [trigger]);

    useEffect(() => {
        connection.onRoomList = e => {
            setRoom(e.rooms![0]);
            dispatch({ actionType: ActionType.RoomList, payload: e.rooms! });
        };
        connection.onUserList = e => {
            dispatch({ actionType: ActionType.UserList, payload: e.users! });
        };
        connection.onJoin = e => {
            dispatch({ actionType: ActionType.Join, payload: e });
        };
        connection.onLeave = e => {
            dispatch({ actionType: ActionType.Leave, payload: e });
        };
        connection.onMessage = e => {
            dispatch({ actionType: ActionType.Message, payload: e });
        };
        connection.onResponse = e => {
            dispatch({ actionType: ActionType.Response, payload: e });
        };
    }, [connection]);

    return (
        <WebsocketContext.Provider value={{ connection, isConnected, isJoined, setIsJoined, room, setRoom, rooms: state.rooms, users: state.users, messages: state.messages, responses: state.responses, dispatch }}>
            {children}
        </WebsocketContext.Provider>
    );
};

export const WebsocketContext = createContext<IWebsocketContext>({} as IWebsocketContext);

export function useWebsocket() {
    return useContext(WebsocketContext);
}

export default WebSocketContextProvider;