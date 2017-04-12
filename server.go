package main

import (
  "net"
  "os"
  "strconv"
  "sync"
  "time"
)

type Server struct {
  listener net.Listener
}

var total int
var mu sync.Mutex

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
  mu = sync.Mutex{}
  cid := 0
  for {
    fd, err := s.listener.Accept()
    if err != nil {
      break;
    }
    cid = cid + 1
    go s.Process(fd,cid,log)
  }
}

func (s *Server) Process(fd net.Conn,cid int,log chan string) error{
  var length = 0
  defer fd.Close()
  go Writelog("start",cid,length,log)
  for {
    buf := make([]byte,4096)
    nr, err := fd.Read(buf)
    if err != nil {
      break
    }
    data := buf[0:nr]
    length += len(data)
    go Amount(len(data))
  }
  go Writelog("exit",cid,length,log)
  return nil
}

func Amount(len int){
    mu.Lock()
    total +=len
    mu.Unlock()
}

func Writelog(status string,cid int,length int,log chan string){
  mu.Lock()
  defer mu.Unlock()
  log <- strconv.FormatInt(time.Now().Unix(),10)+","+status+","+strconv.Itoa(cid)+","+strconv.Itoa(length)+","+strconv.Itoa(total) 
}
