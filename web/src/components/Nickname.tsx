import { useCallback, useEffect, useState } from 'react';
import { LuLoader2 } from "react-icons/lu";
import { ActionType, useWebsocket } from '../contexts/WebSocketContext';
import { v4 as uuidv4 } from "uuid";

const RULE = /^[0-9a-zA-Z]{1,12}$/;

function Nickname() {
    const { connection, isConnected, isJoined, setIsJoined, responses, dispatch } = useWebsocket();
    const [nickname, setNickname] = useState("");
    const [nonce, setNonce] = useState("");
    const [errorText, setErrorText] = useState("");
    const [disabled, setDisabled] = useState(true);
    const [isPending, setIsPending] = useState(false);
    const [sid, setSid] = useState<NodeJS.Timeout>();

    useEffect(() => {
        setDisabled(!isConnected || !nickname.match(RULE) || isPending);
    }, [isConnected, nickname, isPending]);

    useEffect(() => {
        if (nickname.length > 0 && !nickname.match(RULE)) {
            setErrorText("nickname format error");
        } else {
            setErrorText("");
        }
    }, [nickname]);

    useEffect(() => {
        if (responses[nonce] !== undefined) {
            if (responses[nonce] === "") {
                setIsJoined(true);
            } else {
                setErrorText("please try again later");
            }
            dispatch({ actionType: ActionType.DeleteResponse, payload: nonce });
            reset();
            clearTimeout(sid);
        }
    }, [responses, sid]);

    const requestJoin = useCallback((nickname: string) => {
        if (disabled) return;
        let nonce = uuidv4();
        let request = {
            method: "JOIN",
            name: nickname,
            nonce,
        }
        setNonce(nonce);
        setIsPending(true);
        connection.send(JSON.stringify(request));
        setSid(setTimeout(checkTimeout, 5000));
    }, [connection, disabled]);

    const checkTimeout = useCallback(() => {
        reset();
        setErrorText("please try again later");
    }, [isJoined]);

    const reset = useCallback(() => {
        setIsPending(false);
        setNonce("");
    }, []);

    return (
        <>
            {
                !isJoined &&
                <div className="fixed inset-0 w-screen h-screen flex items-center justify-center backdrop-blur">
                    <div className="flex flex-col items-center justify-center h-60 w-96 bg-gray-600 rounded-lg">
                        <div className="flex-1" />
                        <div className="flex-1" />
                        <input className="p-2 rounded text-gray-800 text-lg" type="text" placeholder="nickname" maxLength={12} onChange={e => setNickname(e.target.value)} disabled={!isConnected || isPending} />
                        <div className="flex-1"><span className="text-red-300">{errorText}</span></div >
                        <button className={`bg-gray-700 px-4 py-2 rounded-lg w-20 disabled:cursor-not-allowed ${disabled ? "" : "hover:bg-gray-500 active:bg-gray-400"}`} onClick={() => requestJoin(nickname)} disabled={disabled}>
                            {!isConnected || isPending ? <LuLoader2 className="inline animate-spin" /> : "Join"}
                        </button>
                        <div className="flex-1" />
                    </div>
                </div >
            }
        </>
    )
}

export default Nickname;
