package beater

import (
    "errors"
    "fmt"
    "net/url"
    "time"

    "github.com/elastic/beats/libbeat/beat"
    "github.com/elastic/beats/libbeat/common"
    "github.com/elastic/beats/libbeat/logp"
    "github.com/elastic/beats/libbeat/publisher"

    "github.com/consulthys/logstashbeat/config"
)

const selector = "logstashbeat"

type Logstashbeat struct {
    period          time.Duration
    urls            []*url.URL

    beatConfig      *config.Config

    done            chan struct{}
    client          publisher.Client

    eventsStats     bool
    jvmStats        bool
    processStats    bool
    memStats        bool
    pipelineStats   bool
}

// Creates beater
func New() *Logstashbeat {
    return &Logstashbeat{
        done: make(chan struct{}),
    }
}

/// *** Beater interface methods ***///

func (bt *Logstashbeat) Config(b *beat.Beat) error {

    // Load beater beatConfig
    err := b.RawConfig.Unpack(&bt.beatConfig)
    if err != nil {
        return fmt.Errorf("Error reading config file: %v", err)
    }

    //define default URL if none provided
    var urlConfig []string
    if bt.beatConfig.Logstashbeat.URLs != nil {
        urlConfig = bt.beatConfig.Logstashbeat.URLs
    } else {
        urlConfig = []string{"http://127.0.0.1:9600"}
    }

    bt.urls = make([]*url.URL, len(urlConfig))
    for i := 0; i < len(urlConfig); i++ {
        u, err := url.Parse(urlConfig[i])
        if err != nil {
            logp.Err("Invalid Logstash URL: %v", err)
            return err
        }
        bt.urls[i] = u
    }

    if bt.beatConfig.Logstashbeat.Stats.Events != nil {
        bt.eventsStats = *bt.beatConfig.Logstashbeat.Stats.Events
    } else {
        bt.eventsStats = true
    }

    if bt.beatConfig.Logstashbeat.Stats.JVM != nil {
        bt.jvmStats = *bt.beatConfig.Logstashbeat.Stats.JVM
    } else {
        bt.jvmStats = true
    }

    if bt.beatConfig.Logstashbeat.Stats.Process != nil {
        bt.processStats = *bt.beatConfig.Logstashbeat.Stats.Process
    } else {
        bt.processStats = true
    }

    if bt.beatConfig.Logstashbeat.Stats.Mem != nil {
        bt.memStats = *bt.beatConfig.Logstashbeat.Stats.Mem
    } else {
        bt.memStats = true
    }

    if bt.beatConfig.Logstashbeat.Stats.Pipeline != nil {
        bt.pipelineStats = *bt.beatConfig.Logstashbeat.Stats.Pipeline
    } else {
        bt.pipelineStats = true
    }

    if !bt.eventsStats && !bt.jvmStats && !bt.processStats && !bt.memStats && !bt.pipelineStats {
        return errors.New("Invalid statistics configuration")
    }

    return nil
}

func (bt *Logstashbeat) Setup(b *beat.Beat) error {

    // Setting default period if not set
    if bt.beatConfig.Logstashbeat.Period == "" {
        bt.beatConfig.Logstashbeat.Period = "10s"
    }

    bt.client = b.Publisher.Connect()

    var err error
    bt.period, err = time.ParseDuration(bt.beatConfig.Logstashbeat.Period)
    if err != nil {
        return err
    }

    logp.Debug(selector, "Init logstashbeat")
    logp.Debug(selector, "Period %v\n", bt.period)
    logp.Debug(selector, "Watch %v", bt.urls)
    logp.Debug(selector, "Events statistics %t\n", bt.eventsStats)
    logp.Debug(selector, "JVM statistics %t\n", bt.jvmStats)
    logp.Debug(selector, "Process statistics %t\n", bt.processStats)
    logp.Debug(selector, "Memory statistics %t\n", bt.memStats)
    logp.Debug(selector, "Pipeline statistics %t\n", bt.pipelineStats)

    return nil
}

func (bt *Logstashbeat) Run(b *beat.Beat) error {
    logp.Info("logstashbeat is running! Hit CTRL-C to stop it.")

    for _, u := range bt.urls {
        go func(u *url.URL) {

            ticker := time.NewTicker(bt.period)
            counter := 1
            for {
                select {
                case <-bt.done:
                    goto GotoFinish
                case <-ticker.C:
                }

                timerStart := time.Now()

                if bt.eventsStats {
                    logp.Debug(selector, "Events stats for url: %v", u)
                    events_stats, err := bt.GetEventsStats(*u)

                    if err != nil {
                        logp.Err("Error reading events stats: %v", err)
                    } else {
                        logp.Debug(selector, "Events stats detail: %+v", events_stats)

                        event := common.MapStr{
                            "@timestamp": common.Time(time.Now()),
                            "type": "events",
                            "counter": counter,
                            "events": events_stats.Events,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Logstash events stats sent")
                        counter++
                    }
                }

                if bt.jvmStats {
                    logp.Debug(selector, "JVM stats for url: %v", u)
                    jvm_stats, err := bt.GetJvmStats(*u)

                    if err != nil {
                        logp.Err("Error reading JVM stats: %v", err)
                    } else {
                        logp.Debug(selector, "JVM stats detail: %+v", jvm_stats)

                        event := common.MapStr{
                            "@timestamp":   common.Time(time.Now()),
                            "type":         "jvm",
                            "counter":      counter,
                            "jvm":          jvm_stats.Jvm,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Logstash JVM stats sent")
                        counter++
                    }
                }

                if bt.processStats {
                    logp.Debug(selector, "Process stats for url: %v", u)
                    process_stats, err := bt.GetProcessStats(*u)

                    if err != nil {
                        logp.Err("Error reading process stats: %v", err)
                    } else {
                        logp.Debug(selector, "Process stats detail: %+v", process_stats)

                        event := common.MapStr{
                            "@timestamp": common.Time(time.Now()),
                            "type":       "process",
                            "counter":    counter,
                            "process": process_stats.Process,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Logstash process stats sent")
                        counter++
                    }
                }

                if bt.memStats {
                    logp.Debug(selector, "Memory stats for url: %v", u)
                    mem_stats, err := bt.GetMemStats(*u)

                    if err != nil {
                        logp.Err("Error reading memory stats: %v", err)
                    } else {
                        logp.Debug(selector, "Memory stats detail: %+v", mem_stats)

                        event := common.MapStr{
                            "@timestamp": common.Time(time.Now()),
                            "type":       "mem",
                            "counter":    counter,
                            "mem": mem_stats.Mem,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Logstash memory stats sent")
                        counter++
                    }
                }

                if bt.pipelineStats {
                    logp.Debug(selector, "Pipeline stats for url: %v", u)
                    pipeline_stats, err := bt.GetPipelineStats(*u)

                    if err != nil {
                        logp.Err("Error reading pipeline stats: %v", err)
                    } else {
                        logp.Debug(selector, "Pipeline stats detail: %+v", pipeline_stats)

                        event := common.MapStr{
                            "@timestamp": common.Time(time.Now()),
                            "type":       "pipeline",
                            "counter":    counter,
                            "pipeline": pipeline_stats.Pipeline,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Logstash pipeline stats sent")
                        counter++
                    }
                }

                timerEnd := time.Now()
                duration := timerEnd.Sub(timerStart)
                if duration.Nanoseconds() > bt.period.Nanoseconds() {
                    logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
                }
            }

        GotoFinish:
        }(u)
    }

    <-bt.done
    return nil
}

func (bt *Logstashbeat) Cleanup(b *beat.Beat) error {
    return nil
}

func (bt *Logstashbeat) Stop() {
    logp.Debug(selector, "Stop logstashbeat")
    close(bt.done)
}
