package main

import (
  "net"
  "os"
  "strconv"
  "time"
)

type Server struct {
  listener net.Listener
}

var total = 0

func NewServer() *Server{
  s := new(Server)
  return s
}

func (s *Server) Open(socket string) error {
  listener, err := net.Listen("unix", socket)
  if err != nil {
    return err
  }
  s.listener = listener;
  if err := os.Chmod(socket, 0600); err != nil {
    s.Close()
    return err
  }
  return nil
}

func (s *Server) Close() error{
  if err := s.listener.Close(); err != nil {
    return err;
  }
  return nil
}

func (s *Server) Start(log chan string) {
  ch := make(chan int)
  cid := 0
  go Amount(ch)
  for {
    fd, err := s.listener.Accept()
    if err != nil {
      break;
    }
    cid = cid + 1
    go s.Process(fd,cid,ch,log)
  }
}

func Amount(ch chan int){
  for{
    total += <-ch
  }
}
func (s *Server) Process(fd net.Conn,cid int,ch chan int,log chan string) error{
  var length = 0
  defer fd.Close()
  for {
    buf := make([]byte,512)
    nr, err := fd.Read(buf)
    if err != nil {
      break
    }
    data := buf[0:nr]
    length += len(data)
    ch <- len(data)
    //fmt.Printf("Recieved: %v", string(data));
    /*_, err = fd.Write(data)
    if err != nil {
      return err
    }*/
  }
  log <- strconv.FormatInt(time.Now().Unix(),10)+","+strconv.Itoa(cid)+","+strconv.Itoa(length)+","+strconv.Itoa(total)
  return nil
}
