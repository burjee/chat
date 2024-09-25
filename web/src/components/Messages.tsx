import { useCallback, useEffect, useRef, useState } from "react";
import { LuLoader2 } from "react-icons/lu";
import { ActionType, useWebsocket } from "../contexts/WebSocketContext";
import { v4 as uuidv4 } from "uuid";
import { EventType } from "../libs/socket";

const RULE = /^.{1,300}$/;

function timeToMMSS(datetime: string) {
    const date = new Date(datetime);
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${hours}:${minutes}`;
}

function Messages() {
    const { connection, room, messages, responses, dispatch } = useWebsocket();
    const [text, setText] = useState("");
    const [nonce, setNonce] = useState("");
    const [disabled, setDisabled] = useState(false);
    const [sid, setSid] = useState<NodeJS.Timeout>();
    const [showHint, setShowHint] = useState(false);
    const inputRef = useRef<HTMLInputElement>(null);
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (responses[nonce] !== undefined) {
            if (responses[nonce] === "") {
                setText("");
            } else {
                setShowHint(true);
            }
            dispatch({ actionType: ActionType.DeleteResponse, payload: nonce });
            reset();
            clearTimeout(sid);
        }
    }, [responses, sid]);

    useEffect(() => {
        if (!disabled) {
            inputRef.current?.focus();
        }
    }, [disabled])

    useEffect(() => {
        scrollRef.current?.scrollTo(0, scrollRef.current?.scrollHeight);
    }, [room, messages]);

    const onKeyDown = useCallback((e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            if (!!text.match(RULE)) {
                sendMessage();
                setDisabled(true);
            }
        }
    }, [connection, room, text]);

    const sendMessage = useCallback(() => {
        let nonce = uuidv4();
        let request = {
            method: "MESSAGE",
            message: { room, text },
            nonce,
        }
        setShowHint(false);
        setNonce(nonce);
        connection.send(JSON.stringify(request));
        setSid(setTimeout(checkTimeout, 5000));
    }, [connection, room, text]);

    const checkTimeout = useCallback(() => {
        reset();
        setShowHint(true);
    }, []);

    const reset = useCallback(() => {
        setNonce("");
        setDisabled(false);
    }, []);

    return (
        <div className="flex-[3_3_0%] px-2 relative flex flex-col">
            <div ref={scrollRef} className="mt-6 px-16 flex-1 overflow-y-auto">
                {
                    messages[room]?.map((e) => {
                        if (e.type === EventType.Message) {
                            return (
                                <div key={e.nonce} className="p-2 rounded-lg hover:bg-gray-600">
                                    <span className="text-xs text-gray-400">{timeToMMSS(e.content.datetime!)}</span>
                                    <span className="pl-2 font-bold">{e.content.from}:</span>
                                    <span className="pl-2" >{e.content.message?.text}</span>
                                </div>
                            );
                        } else {
                            return <div key={e.nonce} className="p-2 rounded-lg text-xs text-gray-300">{e.content.from} {e.type === EventType.Join ? "joined" : "left"} </div>;
                        }
                    })
                }
            </div>
            <div className="bg-gray-600 px-5 mx-16 my-4 h-12 rounded-lg flex flex-row relative">
                {showHint && <span className="absolute top-[-0.625rem] text-red-300" >please try again later</span>}
                <input ref={inputRef} className="bg-gray-600 w-full h-12 rounded-lg" maxLength={300} placeholder={`Message #${room}`} value={text} onChange={e => setText(e.target.value)} onKeyDown={onKeyDown} disabled={disabled} />
                {disabled && <LuLoader2 className="h-12 animate-spin" />}
            </div>
        </div>
    );
}

export default Messages;
