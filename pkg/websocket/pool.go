package websocket

import (
	"fmt"
	"time"

	"strconv"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

type Proci struct {
	Pid     int32
	Nombre  string
	Estado  string
	Memoria float32
	User    string
}
type Myresponse struct {
	Totalprocs int32
	Procsr     int32
	Procss     int32
	Procst     int32
	Procsz     int32
	Infop      []Proci
	Cpuu       float64
	Memtotal   uint64
	Memfree    uint64
	Mempercent float64
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				//client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			//for client, _ := range pool.Clients {
			//client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			//}
			break
		case message := <-pool.Broadcast:
			//En teoria aqui va a llegar solo cuando se le de click al boton...
			fmt.Printf("Se va a matar al proceso %v", message.Body)
			var p []*process.Process
			p, _ = process.Processes()
			for _, proc := range p {
				i, err := strconv.ParseInt(message.Body, 10, 32)
				if err != nil {
					panic(err)
				}
				result := int32(i)
				if proc.Pid == result {
					fmt.Println("SI se va a a matar..")
					proc.Kill()
				}
			}
			/*
				for client, _ := range pool.Clients {


					if err := client.Conn.WriteJSON(message); err != nil {
						fmt.Println(err)
						return
					}

				}
			*/
		default:
			time.Sleep(2000 * time.Millisecond)
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(getInfo()); err != nil {
					fmt.Println("Entro en error")
					fmt.Println(err)
					//return
				}
			}

		}

	}
}

func getInfo() Myresponse {
	/*
		type Myresponse struct {
			Totalprocs int32
			Procsr     int32
			Procss     int32
			Procst     int32
			Procsz     int32
			Infop      []Proci
			Cpuu       float64
			Memtotal   uint64
			Memfree    uint64
			Mempercent float64
		}
	*/
	var resp = Myresponse{}
	var cont int32 = 0
	var contr int32 = 0
	var conts int32 = 0
	var contt int32 = 0
	var contz int32 = 0
	var p []*process.Process
	var _ error
	p, _ = process.Processes()
	var infop []Proci
	for _, proc := range p {
		cont++
		nombre, _ := proc.Name()
		status, _ := proc.Status()
		switch status {
		case "S":
			conts++
		case "R":
			contr++
		case "T":
			contt++
		case "Z":
			contz++

		}

		namesp, _ := proc.Username()
		ramp, _ := proc.MemoryPercent()
		infop = append(infop, Proci{proc.Pid, nombre, status, ramp, namesp})

	}

	resp.Totalprocs = cont
	resp.Procsr = contr
	resp.Procss = conts
	resp.Procst = contt
	resp.Procsz = contz
	resp.Infop = infop
	mycpu, _ := cpu.Percent(0, false)
	resp.Cpuu = mycpu[0]
	v, _ := mem.VirtualMemory()
	resp.Memtotal = v.Total
	resp.Memfree = v.Free
	resp.Mempercent = v.UsedPercent
	// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	//return infop
	fmt.Println("Se genero...............")
	return resp

}
