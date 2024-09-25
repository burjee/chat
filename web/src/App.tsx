import WebSocketContextProvider from './contexts/WebSocketContext';
import Rooms from './components/Rooms';
import Messages from './components/Messages';
import Users from './components/Users';
import Nickname from './components/Nickname';

function App() {
  return (
    <WebSocketContextProvider>
      <div className="w-screen h-screen bg-gray-700 flex">
        <Rooms />
        <Messages />
        <Users />
        <Nickname />
      </div>
    </WebSocketContextProvider>
  )
}

export default App
