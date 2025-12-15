package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const serverURL = "http://localhost:8080"
const verticalTimeout = time.Second * 2

type SensorResponse struct {
	Value bool `json:"value"`
}

func main() {
	// test connection to server
	resp, err := http.Get(fmt.Sprintf("%s/tp/ping", serverURL))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// get gripper start pos value
	resp, err = http.Get(fmt.Sprintf("%s/tp/sensor/%s", serverURL, "ns:1, i:2"))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// parse sensor response
	var sensStart SensorResponse
	err = json.NewDecoder(resp.Body).Decode(&sensStart)
	if err != nil {
		panic(err)
	}
	if !sensStart.Value {
		panic("move gripper to start position first")
	}

	// place puck to start tp
	resp, err = http.Post(fmt.Sprintf("%s/tp/puck", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// move gripper down
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/down", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// wait gripper to get down, cause no sensor provided
	time.Sleep(verticalTimeout)

	// stop gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/stop", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// taking puck by gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/open", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/close", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// enable movement of gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/enable", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// move gripper up
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/up", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// wait gripper to get up, cause no sensor provided
	time.Sleep(verticalTimeout)

	// stop gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/stop", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// enable movement of gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/enable", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// start moving gripper to the left
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/left", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	for {
		// get gripper carousel pos value
		resp, err = http.Get(fmt.Sprintf("%s/tp/sensor/%s", serverURL, "ns:1, i:1"))
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			respString, _ := io.ReadAll(resp.Body)
			panic(string(respString))
		}

		// parse sensor response
		var sensCarousel SensorResponse
		err = json.NewDecoder(resp.Body).Decode(&sensCarousel)
		if err != nil {
			panic(err)
		}

		// if gripper not in carousel pos wait and then repeat
		if sensCarousel.Value {
			break
		} else {
			time.Sleep(time.Millisecond)
		}
	}

	// stop gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/stop", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// enable movement of gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/enable", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// move gripper down
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/down", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// wait gripper to get down, cause no sensor provided
	time.Sleep(verticalTimeout)

	// stop gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/stop", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// place puck by gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/open", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// enable movement of gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/enable", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// move gripper up
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/up", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// wait gripper to get up, cause no sensor provided
	time.Sleep(verticalTimeout)

	// stop gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/stop", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}

	// close gripper
	resp, err = http.Post(fmt.Sprintf("%s/tp/gripper/close", serverURL), "", nil)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		respString, _ := io.ReadAll(resp.Body)
		panic(string(respString))
	}
}
