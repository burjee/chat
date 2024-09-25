import { useWebsocket } from "../contexts/WebSocketContext";

function Users() {
    const { users } = useWebsocket();

    return (
        <div className="flex-1 border-l border-l-slate-600 flex flex-col p-5 overflow-y-auto">
            {
                Object.keys(users).
                    filter(k => users[k] !== "left").
                    sort().
                    map(user => <span key={user} className="my-1 py-2 pl-2 rounded-lg text-left hover:bg-gray-500">{user}</span>)
            }
        </div>
    );
}

export default Users;
