package beater

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
    "strconv"
)

const NODE_HOT_THREADS = "/_node/hot_threads"

type Thread struct {
    Name string `json:"name"`
    Percent_of_cpu_time float32 `json:"percent_of_cpu_time"`
    State string `json:"state"`
    Traces []*string `json:"traces"`
}

type HotThread struct {
    Time string `json:"time"`
    Busiest_threads int `json:"busiest_threads"`
    Threads []*Thread `json:"threads"`
}

type HotThreadsStats struct {
    Host string `json:"host"`
    Version string `json:"version"`
    Http_address string `json:"http_address"`
    Hot_threads HotThread `json:"hot_threads"`
}

func (bt *Logstashbeat) GetHotThreads(u url.URL, numThreads int) (*HotThreadsStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + NODE_HOT_THREADS + "?threads=" + strconv.Itoa(numThreads))
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", res.Status)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    threads := &HotThreadsStats{}
    err = json.Unmarshal([]byte(body), &threads)
    if err != nil {
        return nil, err
    }

    return threads, nil
}
