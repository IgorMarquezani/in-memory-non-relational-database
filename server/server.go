package server

import (
	"app/commands"
	"app/config"
	"app/database"
	"app/utils"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Server struct {
	config   *config.DBConfig
	mutex    sync.Mutex
	database *database.Database
	clients  uint
}

func NewServer(conf *config.DBConfig, db *database.Database) (*Server, error) {
	if conf == nil {
		return nil, errors.New("nil config")
	}

	return &Server{
		config:   conf,
		database: db,
	}, nil
}

func (s *Server) RunServer() {
	fmt.Printf("running server on %s at port %d\n", s.config.Host, s.config.Port)

	go s.database.StartDatabase()

	l, err := net.Listen("tcp", s.config.Host+":"+strconv.Itoa(int(s.config.Port)))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error handling connection:", err)
			continue
		}

		s.ChangeClient("sum", 1)

		fmt.Println("connection accepted:", conn.RemoteAddr())

		go s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	for {
		arr := make([]byte, 1024)

		_, err := conn.Read(arr)
		if err != nil && errors.Is(err, io.EOF) {
			s.ChangeClient("sub", 1)
			return
		}
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println("error reading from connection:", err)
			continue
		}

		if err := s.HandleCommand(conn, arr); err != nil {
			conn.Write([]byte(fmt.Sprint("error:", err.Error())))
		}

		s.ChangeClient("sub", 1)
	}
}

func (s *Server) HandleCommand(conn net.Conn, line []byte) error {
	buff := bytes.NewBuffer(line)

	command, err := buff.ReadString(' ')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("(error) ERR wrong number of arguments for 'echo' command")
		}

		return err
	}

	command = command[:len(command)-1]

	line, err = buff.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return errors.New("Error reading echo data:" + err.Error())
	}

	args := strings.Split(string(line), " ")

	switch command {

	case "echo":
		v, err := commands.Echo(args)
		if err != nil {
			return err
		}

		_, err = conn.Write([]byte(`"` + v + `"`))
		if err != nil {
			return errors.New("Error while sending echo:" + err.Error())
		}

	case "set":
		ch, err := commands.Set(args, database.InstructionQueue)
		if err != nil {
			return err
		}

		msg := utils.WaitForChan(ch)
		if msg.Err != nil {
			conn.Write([]byte("(error) ERR " + msg.Err.Error()))
		}

		_, err = conn.Write([]byte("OK"))
		if err != nil {
			fmt.Println("error:", err)
		}

	case "get":
		ch, err := commands.Get(args, database.InstructionQueue)
		if err != nil {
			return err
		}

		msg := utils.WaitForChan(ch)

		if msg.Err != nil {
			conn.Write([]byte("(error) ERR " + msg.Err.Error()))
		}

		if len(msg.Data) == 0 {
			msg.Data = "(nil)"
		}

		_, err = conn.Write([]byte(fmt.Sprintf(`"%s"`, msg.Data)))
		if err != nil {
			fmt.Println("error:", err)
		}

	default:
		return errors.New(fmt.Sprintf("(error) ERR unknown command '%s', with args beginning with: %v", command, line))
	}

	return nil
}

func (s *Server) ChangeClient(op string, value uint) {
	s.mutex.Lock()
	if op == "add" || op == "sum" {
		s.clients += value
	}
	if op == "remove" || op == "sub" {
		s.clients -= value
	}
	s.mutex.Unlock()
}
