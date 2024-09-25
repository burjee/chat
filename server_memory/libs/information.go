package libs

import (
	"encoding/json"
	"server_memory/structure"
	"sync"

	"github.com/spf13/viper"
)

type Information struct {
	lock      *sync.RWMutex
	rooms     map[string]bool
	users     map[string]bool
	room_list []byte
	user_list []byte
}

func NewInformation() *Information {
	i := &Information{
		&sync.RWMutex{},
		make(map[string]bool),
		make(map[string]bool),
		make([]byte, 0),
		make([]byte, 0),
	}

	rooms := viper.GetStringSlice("chat.rooms")
	for _, room := range rooms {
		i.rooms[room] = true
	}

	room_list := structure.RoomList{Type: "ROOM_LIST", Rooms: rooms}
	user_list := structure.UserList{Type: "USER_LIST", Users: make([]string, 0)}
	if room_list_value, err := json.Marshal(room_list); err != nil {
		panic(err)
	} else if user_list_value, err := json.Marshal(user_list); err != nil {
		panic(err)
	} else {
		i.room_list = room_list_value
		i.user_list = user_list_value
	}

	return i
}

func (i *Information) GetList() ([]byte, []byte) {
	return i.getRoomList(), i.getUserList()
}

func (i *Information) getRoomList() []byte {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.room_list
}

func (i *Information) getUserList() []byte {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.user_list
}

func (i *Information) HasRoom(room string) bool {
	i.lock.RLock()
	defer i.lock.RUnlock()
	_, ok := i.rooms[room]
	return ok
}

func (i *Information) AddUser(name string) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.users[name]; ok {
		return false
	} else {
		i.users[name] = true
		i.updateUserList()
		return true
	}
}

func (i *Information) RemUser(name string) bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	if _, ok := i.users[name]; ok {
		delete(i.users, name)
		i.updateUserList()
		return true
	} else {
		return false
	}
}

func (i *Information) updateUserList() error {
	users := make([]string, 0, len(i.users))
	for user := range i.users {
		users = append(users, user)
	}

	user_list := structure.UserList{Type: "USER_LIST", Users: users}

	value, err := json.Marshal(user_list)
	if err != nil {
		return err
	}

	i.user_list = value
	return nil
}
