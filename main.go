package main

import (
  "os"
  "os/signal"
  "io/ioutil"
  "fmt"
  "syscall"
  "strconv"
  "log"
  "log/syslog"
)

func main() {
  log.SetFlags(log.Lshortfile)
  tempDir, err := ioutil.TempDir("", "counter-socket.")
  if err != nil {
    log.Printf("error: %v\n", err)
    return
  }
  msg := make(chan string)
  pid := strconv.Itoa(os.Getpid())
  socket := tempDir + "/server." + pid
  if err := os.Chmod(tempDir, 0700); err != nil {
    log.Printf("error: %v\n", err)
    return
  }
  defer func() {
    if err := os.Remove(tempDir); err != nil {
      log.Printf("error: %v\n", err)
    }
  }()
  server := NewServer()
  if err := server.Open(socket); err != nil {
    log.Printf("error: %v\n", err)
    return;
  }
  go Startlogger(msg)
  registerShutdown(server)
  fmt.Printf("GOLANG_SAMPLE_SOCK=%s\n",os.Getenv("GOLANG_SAMPLE_SOCK"))
  fmt.Printf("GOLANG_SAMPLE_SOCK=%v;export GOLANG_SAMPLE_SOCK;\n", socket)
  fmt.Printf("GOLANG_SAMPLE_PID=%v;export GOLANG_SAMPLE_PID;\n", pid)
  server.Start(msg)
}

func registerShutdown(server *Server) {
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    interrupt := 0
    for {
      s := <-c
      switch s {
      case os.Interrupt:
        if (interrupt == 0) {
          fmt.Println("Interrupt...")
          interrupt++
          continue
        }
      }
      break
    }
    if err := server.Close(); err != nil {
      log.Printf("error: %v\n", err)
    }
  }()
}

func Startlogger(msg chan string){
  logger,_ := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "my-daemon")
  log.SetOutput(logger)
  for{
     log.Println(<-msg)
  }
}

