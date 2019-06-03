package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/D4-project/d4-golang-utils/config"
	"github.com/gomodule/redigo/redis"
)

type (
	conf struct {
		redisHost    string
		redisPort    string
		redisDB      int
		redisQueue   string
		balboaSocket string
	}
)

var (
	confdir      = flag.String("c", "conf.sample", "configuration directory")
	connectRedis = true
	cr           redis.Conn
)

func main() {
	// Control Chan
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)

	// Usage and flags
	flag.Usage = func() {
		fmt.Printf("analyzer-d4-balboa - export D4 Type 8 to Balboa UNIX socket:\n\n")
		fmt.Printf("\n")
		fmt.Printf("Usage:\n\n analyzer-d4-balboa -c  config_directory\n")
		fmt.Printf("\n")
		fmt.Printf("Configuration:\n\n")
		fmt.Printf(" The configuration settings are stored in files in the configuration directory\n")
		fmt.Printf(" specified with the -c command line switch.\n\n")
		fmt.Printf("Files in the configuration directory:\n")
		fmt.Printf("\n")
		fmt.Printf(" redis - d4 server\n")
		fmt.Printf("       | host:port/db\n")
		fmt.Printf(" redis_queue - type and uuid of the redis queue\n")
		fmt.Printf("          | type:uuid \n")
		fmt.Printf(" balboa_socket - socket file to balboa\n")
		fmt.Printf("       | /tmp/balboa.sock\n")
		fmt.Printf("\n")
		flag.PrintDefaults()
	}

	// Config
	c := conf{}
	flag.Parse()
	if flag.NFlag() == 0 || *confdir == "" {
		flag.Usage()
		os.Exit(1)
	} else {
		*confdir = strings.TrimSuffix(*confdir, "/")
		*confdir = strings.TrimSuffix(*confdir, "\\")
	}

	// Parse Redis Config
	tmp := config.ReadConfigFile(*confdir, "redis")
	ss := strings.Split(string(tmp), "/")
	if len(ss) <= 1 {
		log.Fatal("Missing Database in Redis config: should be host:port/database_name")
	}
	c.redisDB, _ = strconv.Atoi(ss[1])
	var ret bool
	ret, ss[0] = config.IsNet(ss[0])
	if !ret {
		sss := strings.Split(string(ss[0]), ":")
		c.redisHost = sss[0]
		c.redisPort = sss[1]
	}
	c.redisQueue = string(config.ReadConfigFile(*confdir, "redis_queue"))
	c.balboaSocket = string(config.ReadConfigFile(*confdir, "balboa_socket"))
	//TODO: handle empty ...

	initRedis(c.redisHost, c.redisPort, c.redisDB)
	defer cr.Close()
	cs, err := net.Dial("unix", c.balboaSocket)
	//defer cs.Close()
	if err != nil {
		panic(err)
	}
	// pop redis queue
	for {
		dnsLine, err := redis.String(cr.Do("LPOP", "analyzer:"+c.redisQueue))
		if err != nil {
			log.Fatal("Queue processed")
		}
		// Write in Balboa socket
		cs.Write([]byte(dnsLine))
		//TODO: Check that it works...

		// Exit Signal Handle
		select {
		case <-s:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			continue
		}
	}

}

func initRedis(host string, port string, d int) {
	err := errors.New("")
	cr, err = redis.Dial("tcp", host+":"+port, redis.DialDatabase(d))
	if err != nil {
		panic(err)
	}
}
