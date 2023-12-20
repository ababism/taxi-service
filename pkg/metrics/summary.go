package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"strings"
	"sync"
)

const quantileLabel = "quantile"

var DefaultObjectives = map[float64]float64{
	0.1:    0.09,    // quantile 10
	0.5:    0.05,    // quantile 50
	0.75:   0.025,   // quantile 75
	0.9:    0.01,    // quantile 90
	0.95:   0.005,   // quantile 95
	0.99:   0.001,   // quantile 99
	0.999:  0.0001,  // quantile 99.9
	0.9999: 0.00001, // quantile 99.99
	1.0:    0.0,     // quantile 100
}

// Summary (Сводка) - это метрика собирает отдельные наблюдения из потока событий или выборок и суммирует их способом,
// аналогичным традиционной сводной статистике: 1. сумма наблюдений, 2. количество наблюдений, 3. ранговые оценки.
type Summary interface {
	Observe(val float64)
	getMetric() ([]*MetricDTO, error)
}

// SummaryMetric - дополнительная оболочка необходимая для получения Write
type SummaryMetric interface {
	Observe(val float64)
	Write(*dto.Metric) error
}

type summary struct { // implement Summary
	prom SummaryMetric
}

func (s *summary) Observe(val float64) {
	s.prom.Observe(val)
}

// getMetric дает альтернативу по получению метрик как пакет promhttp
func (s *summary) getMetric() ([]*MetricDTO, error) {
	mS := &dto.Metric{}
	if err := s.prom.Write(mS); err != nil {
		return nil, err
	}

	var metrics []*MetricDTO
	for _, quantile := range mS.Summary.GetQuantile() {
		metric := &MetricDTO{
			Type:         "Summary",
			Labels:       []string{quantileLabel},
			LabelsValues: []string{fmt.Sprintf("%0.2f", quantile.GetQuantile()*100)},
			Value:        quantile.GetValue(),
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// SummaryOpts опции для создания метрики
type SummaryOpts struct {
	Namespace   string
	Name        string
	Description string
	Objectives  map[float64]float64
}

// SummaryVec метрика сводки с поддержкой labels
type SummaryVec struct { // implement RegistryMetric
	opts    SummaryOpts
	mutex   sync.RWMutex
	labels  []string
	metrics map[string]Summary

	constructor *prometheus.SummaryVec
}

type SummaryVecMetric struct { // implement SummaryMetric
	observer prometheus.Observer
	writer   prometheus.Metric
}

// GetOrRegisterSummaryVec создает или возвращает уже созданную фабрику метрик Summary.
// HistogramVec это множественная метрика сводок. Сводки можно выделять в группы по разным лейблам.
// С помощью WithLabelValues вы получите необходимую группу сводки и будете изменять только ее.
func GetOrRegisterSummaryVec(opts SummaryOpts, labels []string) *SummaryVec {
	if rm := globalRegistry.getMetric(opts.Name); rm != nil {
		if vm, ok := rm.(*SummaryVec); ok {
			return vm
		}
		globalRegistry.markCollision(opts.Name, "SummaryVec")
	}

	promConstructor := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  opts.Namespace,
			Name:       opts.Name,
			Help:       opts.Description,
			Objectives: opts.Objectives,
		},
		labels,
	)

	summaryVec := &SummaryVec{
		opts:        opts,
		labels:      labels,
		metrics:     map[string]Summary{},
		constructor: promConstructor,
	}

	globalRegistry.addMetric(opts.Name, summaryVec)

	return summaryVec
}

// WithLabelValues возвращает необходимую сводку Summary в зависимости от переданных значений
func (s *SummaryVec) WithLabelValues(labelValues ...string) Summary {
	return s.getOrRegisterLabelMetric(labelValues)
}

// getMetric дает альтернативу по получению метрик как пакет promhttp
func (s *SummaryVec) getMetrics() ([]*MetricDTO, error) {
	var metrics []*MetricDTO
	for group, summary := range s.metrics {
		groupMetrics, _ := summary.getMetric()
		for _, metric := range groupMetrics {
			metric.Namespace = s.opts.Namespace
			metric.Name = s.opts.Name
			metric.Description = s.opts.Description

			metric.Labels = append(s.labels, metric.Labels...)

			// Получаем нормальный список labels
			labelsValues := strings.Split(group, GroupKeySeparator)
			metric.LabelsValues = append(labelsValues, metric.LabelsValues...)

			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

// getOrRegisterLabelMetric это специальная конструкция - фабрика,
// которая в зависимости от переданных значений создаст вам новую и полноценную сводку
func (s *SummaryVec) getOrRegisterLabelMetric(labelValues []string) Summary {
	globalRegistry.checkLabels(labelValues, s.labels, s.opts.Name)

	key := strings.Join(labelValues, GroupKeySeparator)

	s.mutex.RLock()
	cm, exist := s.metrics[key]
	s.mutex.RUnlock()

	if !exist {
		// Игнорируем ошибку, тк проверка в начале этого метода делает то же самое
		writer, _ := s.constructor.MetricVec.GetMetricWithLabelValues(labelValues...)
		prom := SummaryVecMetric{
			observer: s.constructor.WithLabelValues(labelValues...),
			writer:   writer,
		}
		// оболочка метрики с которой и взаимодействует приложение
		cm = &summary{
			prom: prom,
		}

		// write
		s.mutex.Lock()
		s.metrics[key] = cm
		s.mutex.Unlock()
	}

	return cm
}

func (s SummaryVecMetric) Observe(val float64) {
	s.observer.Observe(val)
}

func (s SummaryVecMetric) Write(dto *dto.Metric) error {
	return s.writer.Write(dto)
}
