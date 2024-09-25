package libs

import (
	"context"
	"encoding/json"
	"server_redis/structure"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Information struct {
	redis_client *redis.Client
}

func NewInformation(redis_client *redis.Client) *Information {
	rooms := viper.GetStringSlice("chat.rooms")

	if err := redis_client.FlushAll(context.Background()).Err(); err != nil {
		panic(err)
	}

	room_list := structure.RoomList{Type: "ROOM_LIST", Rooms: rooms}
	user_list := structure.UserList{Type: "USER_LIST", Users: make([]string, 0)}

	if err := redis_client.SAdd(context.Background(), "rooms", rooms).Err(); err != nil {
		panic(err)
	} else if room_list_value, err := json.Marshal(room_list); err != nil {
		panic(err)
	} else if user_list_value, err := json.Marshal(user_list); err != nil {
		panic(err)
	} else if err := redis_client.Set(context.Background(), "room_list", room_list_value, 0).Err(); err != nil {
		panic(err)
	} else if err := redis_client.Set(context.Background(), "user_list", user_list_value, 0).Err(); err != nil {
		panic(err)
	}

	return &Information{redis_client}
}

func (i *Information) GetList() ([]byte, []byte, error) {
	room_list, err0 := i.getRoomList()
	user_list, err1 := i.getUserList()

	if err0 != nil {
		return nil, nil, err0
	}
	if err1 != nil {
		return nil, nil, err1
	}

	return room_list, user_list, nil
}

func (i *Information) getRoomList() ([]byte, error) {
	return i.redis_client.Get(context.Background(), "room_list").Bytes()
}

func (i *Information) getUserList() ([]byte, error) {
	return i.redis_client.Get(context.Background(), "user_list").Bytes()
}

func (i *Information) HasRoom(room string) (bool, error) {
	return i.redis_client.SIsMember(context.Background(), "rooms", room).Result()
}

func (i *Information) AddUser(name string) bool {
	result, err := i.redis_client.SAdd(context.Background(), "users", name).Result()
	if err != nil {
		return false
	}

	i.updateUserList()
	return result == 1
}

func (i *Information) RemUser(name string) bool {
	result, err := i.redis_client.SRem(context.Background(), "users", name).Result()
	if err != nil {
		return false
	}

	i.updateUserList()
	return result == 1
}

func (i *Information) updateUserList() error {
	users, err := i.redis_client.SMembers(context.Background(), "users").Result()
	if err != nil {
		return err
	}

	user_list := structure.UserList{Type: "USER_LIST", Users: users}

	value, err := json.Marshal(user_list)
	if err != nil {
		return err
	}

	return i.redis_client.Set(context.Background(), "user_list", value, 0).Err()
}
