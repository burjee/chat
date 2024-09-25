import { GoHash } from "react-icons/go";
import { useWebsocket } from "../contexts/WebSocketContext";

function Rooms() {
    const { room, setRoom, rooms } = useWebsocket();

    return (
        <div className="flex-1 bg-gray-800 flex flex-col p-5 overflow-y-auto">
            {
                rooms.sort().map(_room =>
                    <button key={_room} onClick={() => setRoom(_room)} className={`my-1 py-2 pl-2 rounded-lg text-left active:bg-gray-600 ${_room === room ? "bg-gray-600" : "hover:bg-gray-700"}`}>
                        <GoHash className="inline mr-4" />{_room}
                    </button>)
            }
        </div>
    );
}

export default Rooms;
