package main

import (
	"fmt"
	"os"
	consul "github.com/hashicorp/consul/api"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Consul struct{ consul.Client }

var session *r.Session

func prestart() {
	config := consul.DefaultConfig()
	config.Address = os.Getenv("CONSUL_ADDRESS")
	fmt.Println(config.Address)
	client, err := consul.NewClient(config)
	if err != nil {
	    panic(err)
	}
	services, _, sErr := client.Health().Service(os.Getenv("SERVICE_NAME"), ``, true, nil)
	if sErr != nil {
	    panic(sErr)
	}
	if (len(services) == 0) {
		//first node
		fmt.Println("This is the first node")
		return
	}
	f, fErr := os.OpenFile("/etc/rethink.conf", os.O_RDWR, 0)
	defer f.Close()
	if fErr != nil {
	    panic(fErr)
	}
	fmt.Println("Other nodes exist. Writing configuration file")
	for i := 0; i < len(services); i++ {
		service := services[i].Service
		_, wErr := f.WriteString(fmt.Sprintf("join=%s:29015\n", service.Address))
		if wErr != nil {
		    panic(wErr)
		}
	}
	f.Sync()
}

type ServerStatusNetwork struct {
	Hostname string `gorethink:"hostname"`
}

type ServerStatus struct {
	Name string `gorethink:"name"`
	Network ServerStatusNetwork `gorethink:"network"`
}

func healthCheck() {
	session, err := r.Connect(r.ConnectOpts{
	    Address: "172.17.0.4:28015",
	})
	if err != nil {
	    panic(err)
	}
	res, dbErr := r.DB("rethinkdb").Table("server_status").Run(session)
	defer res.Close()
	if dbErr != nil {
		os.Exit(1)
	}
	found := false
	hostname, hErr := os.Hostname()
	if hErr != nil {
	    panic(hErr)
	}
	servers := []ServerStatus{}
	err = res.All(&servers)
	for i := 0; i < len(servers); i++ {
		server := servers[i]
		if (server.Network.Hostname == hostname) {
			found = true
		}
	}
	if (found) {
		fmt.Println("Found self")
		os.Exit(0)
	}
	fmt.Println("Unable to find self")
	os.Exit(1);
}

func main() {
	argsWithoutProg := os.Args[1:]
	action := `healthCheck`
	if (len(argsWithoutProg) > 0) {
		action = argsWithoutProg[0];
	}
	switch action {
		case `healthCheck`:
			healthCheck()
		case `prestart`:
			prestart()
		default:
			panic("Unknown task")
	}
}