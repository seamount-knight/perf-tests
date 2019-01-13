package test

import "k8s.io/perf-tests/clusterloader2/pkg/config"

func GetReport(ctx Context) SummaryReporter {
	switch  ctx.GetClusterLoaderConfig().ReportType {
	case config.Local:
		return NewReportLocal()
	case config.KS3:
		config := &KS3Config{
			"",
			"",
			"",
			"",
			true,
			false,
			"",
			"",
		}
		return NewReportKS3(config)
	default:
		return NewReportLocal()
	}
	return nil
}
