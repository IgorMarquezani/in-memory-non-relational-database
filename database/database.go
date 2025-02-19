package database

import (
	"app/instruction"
	"app/utils"
	"fmt"
)

type Database struct {
	hash map[any]Value
}

func NewDatabase() *Database {
	return &Database{
		hash: make(map[any]Value),
	}
}

func (d *Database) Set(key any, value Value) {
	d.hash[key] = value
}

func (d *Database) Del(key any) {
	delete(d.hash, key)
}

func (d *Database) Get(key any) Value {
	v, _ := d.hash[key]
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
        fmt.Println(d.hash)

				msg.Channel <- instruction.Message{}

			case "del":

			case "get":
				data := d.Get(msg.Key)

				str, _ := utils.AnyToString(data.Data)

				msg.Channel <- instruction.Message{
					Data: str,
				}

			default:
				fmt.Println("Invalid command")
			}

		default:
		}
	}
}
