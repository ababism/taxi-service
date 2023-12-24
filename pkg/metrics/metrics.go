package metrics

import "go.uber.org/zap"

const GroupKeySeparator = "||"

type AppInfo struct {
	Name        string
	Environment string
	Version     string
}

// MetricDTO является неким DTO, которую возвращают все метрики при запросе у них текущего значения
type MetricDTO struct {
	Namespace    string
	Name         string
	Description  string
	Type         string
	Labels       []string
	LabelsValues []string
	Value        float64
}

// InitOnce конфигурирует работу пакета. Аналог init()
func InitOnce(cfg *Config, logger *zap.Logger, app AppInfo) {
	// Создание реестра метрик
	initRegistry(logger.Sugar().Errorf)
}
