package internal

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	EventCounter      *prometheus.GaugeVec
	SubscriberCounter *prometheus.GaugeVec
	PublishCounter    *prometheus.CounterVec
	DistributeCounter *prometheus.CounterVec
}

func MetricHandler(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Println(http.ListenAndServe(":"+port, nil))
}

func NewMetrics(enabled bool, namespace string, port string) Metrics {
	var metric Metrics

	if enabled {
		go MetricHandler(port)
	}

	metric.EventCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{ //nolint:exhaustruct
		Namespace: namespace,
		Subsystem: "qsse",
		Name:      "topic_event_count",
		Help:      "count of events in eventsource",
	}, []string{"topic"})

	metric.SubscriberCounter = promauto.NewGaugeVec(prometheus.GaugeOpts{ //nolint:exhaustruct
		Namespace: namespace,
		Subsystem: "qsse",
		Name:      "topic_subscriber_count",
		Help:      "count of topic's subscribers",
	}, []string{"topic"})

	metric.PublishCounter = promauto.NewCounterVec(prometheus.CounterOpts{ //nolint:exhaustruct
		Namespace: namespace,
		Subsystem: "qsse",
		Name:      "event_publish_total",
		Help:      "count total success published events",
	}, []string{"topic"})

	metric.DistributeCounter = promauto.NewCounterVec(prometheus.CounterOpts{ //nolint:exhaustruct
		Namespace: namespace,
		Subsystem: "qsse",
		Name:      "event_distribute_total",
		Help:      "count total success distributed events",
	}, []string{"topic"})

	return metric
}

func (m Metrics) IncEvent(topic string) {
	m.EventCounter.With(map[string]string{
		"topic": topic,
	}).Inc()
}

func (m Metrics) DecEvent(topic string) {
	m.EventCounter.With(map[string]string{
		"topic": topic,
	}).Dec()
}

func (m Metrics) IncSubscriber(topic string) {
	m.SubscriberCounter.With(map[string]string{
		"topic": topic,
	}).Inc()
}

func (m Metrics) DecSubscriber(topic string) {
	m.SubscriberCounter.With(map[string]string{
		"topic": topic,
	}).Dec()
}

func (m Metrics) IncPublishEvent(topic string) {
	m.PublishCounter.With(map[string]string{
		"topic": topic,
	}).Inc()
}

func (m Metrics) IncDistributeEvent(topic string) {
	m.DistributeCounter.With(map[string]string{
		"topic": topic,
	}).Inc()
}
