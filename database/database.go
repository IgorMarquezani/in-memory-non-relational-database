package database

import (
	"app/instruction"
	"app/utils"
	"errors"
	"fmt"
)

type Database struct {
	hash map[string]Value
}

func NewDatabase() *Database {
	return &Database{
		hash: make(map[string]Value),
	}
}

func (d *Database) Set(key string, value Value) {
	d.hash[key] = value
}

func (d *Database) Del(key string) {
	delete(d.hash, key)
}

func (d *Database) Get(key string) Value {
	v := d.hash[key]
	return v
}

func (d *Database) StartDatabase() {
	for {
		select {
		case msg := <-InstructionQueue:
			switch msg.Command {
			case "set":
				value, err := NewValue(msg.Data, msg.Type)
				if err != nil {
					fmt.Println("(error) ERR " + err.Error())
				}

				d.Set(msg.Key, value)

				msg.Channel <- instruction.Message{}

			case "del":

			case "get":
				data := d.Get(msg.Key)

				str, _ := utils.AnyToString(data.Data)

				msg.Channel <- instruction.Message{
					Data: str,
				}

			default:
				msg.Channel <- instruction.Message{
					Err: errors.New("Invalid command"),
				}
			}
		}
	}
}
