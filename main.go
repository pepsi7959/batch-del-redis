package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Queue struct {
	RedisHost  string
	FileName   string
	MaxReaders int
	Buffer     chan string
	Done       chan bool
}

func (q Queue) producer() {
	f, err := os.Open(q.FileName)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		f.Close()
	}()

	scanner := bufio.NewScanner(f)

	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		q.Buffer <- line
		idx = idx + 1
		fmt.Printf("\r\033[0;35m deleted %d record \033[0m", idx)
	}

	q.Done <- true

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (q Queue) reader() {
	rc := ClusterConnect(q.RedisHost)
	for {
		select {
		case line := <-q.Buffer:
			ClusterDelete(rc, line)
			time.Sleep(10 * time.Millisecond)
		case <-q.Done:
			fmt.Printf("read: finish\n")
		}
	}
}

func (q Queue) startReader() {
	if q.MaxReaders == 0 {
		q.MaxReaders = 1
	}

	for i := 0; i < q.MaxReaders; i = i + 1 {
		go q.reader()
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Println("::help::")
		log.Println("  - filename    path to list of keys file. ex: 'del_redis_key.txt'")
		log.Println("  - redis-host  list of hosts. ex: '10.100.200.10:26379,10.100.201.10:26379,10.100.202.10:26379'")
		log.Fatalln("./batch-del-redis <filename> <redis-host>")
	}

	filename := args[0]
	redisHost := args[1]

	q := Queue{
		Buffer:     make(chan string, 1),
		Done:       make(chan bool, 1),
		MaxReaders: 10,
		FileName:   filename,
		RedisHost:  redisHost,
	}

	q.startReader()
	q.producer()
}
