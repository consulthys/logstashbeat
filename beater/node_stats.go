package beater

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const NODE_JVM_STATS = "/_node/stats/jvm"
const NODE_PROCESS_STATS = "/_node/stats/process"
const NODE_PIPELINE_STATS = "/_node/stats/pipeline"

type MemPool struct {
    Peak_used_in_bytes uint64 `json:"peak_used_in_bytes"`
    Used_in_bytes uint64 `json:"used_in_bytes"`
    Peak_max_in_bytes uint64 `json:"peak_max_in_bytes"`
    Max_in_bytes uint64 `json:"max_in_bytes"`
    Committed_in_bytes uint64 `json:"committed_in_bytes"`
}

type JvmStats struct {
    Jvm struct {
        Threads struct {
            Count      uint64 `json:"count"`
            Peak_count uint64 `json:"peak_count"`
        } `json:"threads"`
        Mem struct {
            Heap_used_in_bytes int64 `json:"heap_used_in_bytes"`
            Heap_used_percent int64 `json:"heap_used_percent"`
            Heap_committed_in_bytes  int64 `json:"heap_committed_in_bytes"`
            Heap_max_in_bytes  int64 `json:"heap_max_in_bytes"`
            Non_heap_used_in_bytes  int64 `json:"non_heap_used_in_bytes"`
            Non_heap_committed_in_bytes  int64 `json:"non_heap_committed_in_bytes"`
            Pools struct {
                Survivor MemPool `json:"survivor"`
                Old MemPool `json:"old"`
                Young MemPool `json:"young"`
            } `json:"pools"`
        } `json:"mem"`
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

type PipelineElementEvents struct {
    Duration_in_millis uint64 `json:"duration_in_millis"`
    In uint64 `json:"in"`
    Out uint64 `json:"out"`
}

type PipelineEvents struct {
    Name string `json:"name"`
    Id string `json:"id"`
    Matches uint64 `json:"matches,omitempty"`
    Failures uint64 `json:"failures,omitempty"`
    Formats uint64 `json:"formats,omitempty"`
    PatternsPerField interface{} `json:"patterns_per_field,omitempty"`
    Events PipelineElementEvents `json:"events"`
}

type Events struct {
    In uint64 `json:"in"`
    Filtered uint64 `json:"filtered"`
    Out uint64 `json:"out"`
}

type PipelineStats struct {
    Pipeline struct {
        Events   Events `json:"events"`
        Plugins struct {
            Inputs  []*PipelineEvents `json:"inputs"`
            Filters []*PipelineEvents `json:"filters"`
            Outputs []*PipelineEvents `json:"outputs"`
        } `json:"plugins"`
    }
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

func (bt *Logstashbeat) GetPipelineStats(u url.URL) (*PipelineStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + NODE_PIPELINE_STATS)
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
    print([]byte(body))
    stats := &PipelineStats{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }

    return stats, nil
}