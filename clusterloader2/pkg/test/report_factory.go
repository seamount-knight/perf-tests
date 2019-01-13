package test

import "k8s.io/perf-tests/clusterloader2/pkg/config"

func GetReport(ctx Context) SummaryReporter {
	clusterLoaderConfig := ctx.GetClusterLoaderConfig()
	switch  clusterLoaderConfig.ReportType {
	case config.Local:
		return NewReportLocal()
	case config.KS3:
		config := &KS3Config{
			clusterLoaderConfig.AccessKeyID,
			clusterLoaderConfig.AccessKeySecret,
			clusterLoaderConfig.Region,
			"",
			true,
			false,
			clusterLoaderConfig.Bucket,
			"",
		}
		return NewReportKS3(config)
	default:
		return NewReportLocal()
	}
	return nil
}
