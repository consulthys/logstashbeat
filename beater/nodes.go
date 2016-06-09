package beater

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const NODE_EVENTS_STATS = "/_node/stats/events"
const NODE_JVM_STATS = "/_node/stats/jvm"
const NODE_PROCESS_STATS = "/_node/stats/process"

type EventsStats struct {
    Events struct {
        In uint64 `json:"in"`
        Filtered uint64 `json:"filtered"`
        Out uint64 `json:"out"`
    }
}

type JvmStats struct {
    Jvm struct {
        Threads struct {
            Count      uint64 `json:"count"`
            Peak_count uint64 `json:"peak_count"`
        } `json:"threads"`
    }
}
type ProcessStats struct {
    Process struct {
        Open_file_descriptors int64 `json:"open_file_descriptors"`
        Peak_open_file_descriptors int64 `json:"peak_open_file_descriptors"`
        Max_file_descriptors  int64 `json:"max_file_descriptors"`
        Cpu struct {
            Percent         uint64 `json:"percent"`
            Total_in_millis uint64 `json:"total_in_millis"`
        } `json:"cpu"`
        Mem struct {
            Total_virtual_in_bytes uint64 `json:"total_virtual_in_bytes"`
        } `json:"mem"`
    }
}

func (bt *Logstashbeat) GetEventsStats(u url.URL) (*EventsStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + NODE_EVENTS_STATS)
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

    stats := &EventsStats{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }

    return stats, nil
}

func (bt *Logstashbeat) GetJvmStats(u url.URL) (*JvmStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + NODE_JVM_STATS)
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

    stats := &JvmStats{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }

    return stats, nil
}

func (bt *Logstashbeat) GetProcessStats(u url.URL) (*ProcessStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + NODE_PROCESS_STATS)
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

    stats := &ProcessStats{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }

    return stats, nil
}