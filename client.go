package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	var ip string
	var port int
	fmt.Scanf("%s %d", &ip, &port)
	url := fmt.Sprintf("http://%s:%d", ip, port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error processing url: ", url)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status, resp.StatusCode, resp.Proto,
			resp.ProtoMajor, resp.ProtoMinor)
	}

	robots, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", robots)
}
