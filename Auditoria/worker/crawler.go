package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	WORKERS     = 6
	CRAWLERNAME = "TEST"
	TIMEWORK    = 5 * time.Minute
	DELAY       = 1 * time.Second
)

//Crawler struct
type Crawler struct {
	Count       int
	Atualizado  time.Time
	RunChan     chan bool
	Endpoint    string
	NumRequests int
}

//Crawler Start
func NewCrawler(endpoint string, num int) (*Crawler, error) {
	if len(endpoint) == 0 {
		return nil, fmt.Errorf("Endpoint vazio")
	}

	f := &Crawler{}

	//Channels
	f.RunChan = make(chan bool, 0)
	f.Endpoint = endpoint
	f.NumRequests = num

	return f, nil
}

//Ativa - Run once
func (c *Crawler) Ativa() {
	c.RunChan <- true
}

func (c *Crawler) timer() {
	for {
		time.Sleep(TIMEWORK)
		c.RunChan <- true
	}
}

type Request struct {
	ID int
}

// Run - Crawler Loop
func (c *Crawler) Run() {

	//Request para workers
	workChan := make(chan Request, 0)

	for i := 0; i < WORKERS; i++ {
		log.Printf("Ativando worker full %d\n", i)
		go c.crawlerWorker(workChan, i)
	}

	go c.timer()

	for {
		log.Printf("Ciclo inicio\n")
		go c.Atualizar(workChan)
		<-c.RunChan
	}
}

// Atualizar - Gerar requisicao de atializacao
func (c *Crawler) Atualizar(wr chan Request) {
	c.Atualizado = time.Now()

	for i := 1; i < 100; i++ {
		request := Request{
			ID: i,
		}
		wr <- request
	}

}

func (c *Crawler) crawlerWorker(workReqChan chan Request, id int) {
	for wrk := range workReqChan {

		fmt.Printf("Processando %d\n", wrk.ID)

		req, err := http.NewRequest("GET", c.Endpoint, nil)

		client := &http.Client{
			Timeout: 300 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error %s", err)
			continue
		}

		if resp.StatusCode != 200 {
			log.Printf("response Status: %s ", resp.Status)
			resp.Body.Close()
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error %s", err)
			continue
		}
		fmt.Printf("Body %s", body)

		resp.Body.Close()

		time.Sleep(DELAY * time.Duration(int64(wrk.ID)))

	}
}
