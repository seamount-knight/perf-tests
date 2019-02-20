package test

import (
	"k8s.io/perf-tests/clusterloader2/api"
	"k8s.io/perf-tests/clusterloader2/pkg/measurement"
	"path"
	"github.com/softlns/ks3-sdk-go/ks3"
	"strings"
)

type ReportKS3 struct {
	*ks3.Bucket
	Encrypt   bool
	RootDirectory string
}

type KS3Config struct {
	AccessKeyID string
	AccessKeySecret string
	Region string
	Endpoint string
	Secure bool
	Internal bool
	Bucket string
	RootDirectory string
}

func NewReportKS3 (config *KS3Config) *ReportKS3 {
	ks3obj, err := ks3.New(config.AccessKeyID, config.AccessKeySecret, config.Region, config.Secure,
		config.Internal, config.Endpoint)
	if err != nil {
		panic(err)
	}
	ks3.SetDebug(false)

	bucket := ks3obj.Bucket(config.Bucket)
	return &ReportKS3{bucket, true, config.RootDirectory}
}

func (r *ReportKS3) ReportSummary(ctx Context, conf *api.Config, summary measurement.Summary) error {
	summaryText, err := summary.PrintSummary()
	if err != nil {
		return err
	}
	if summaryText == "" {
		return nil
	}

	fileNmae := path.Join(ctx.GetClusterLoaderConfig().TestJob, ctx.GetClusterLoaderConfig().BuildNumber, summary.SummaryName()+"_"+conf.Name+".txt")
	err = r.Put2KS3(fileNmae, summaryText)
	if err != nil {
		return err
	}
	lastBuildFileName := path.Join(ctx.GetClusterLoaderConfig().TestJob, "lastBuildNum")
	return r.Put2KS3(lastBuildFileName, ctx.GetClusterLoaderConfig().BuildNumber)
}

func (r *ReportKS3) Put2KS3(path, summaryText string)error  {
	if summaryText == "" {
		return nil
	}
	return r.Put(r.ks3Path(path), []byte(summaryText), "application/octet-stream", "private", r.getOptions())
}

func (r *ReportKS3) ks3Path(path string) string {
	return strings.TrimLeft(strings.TrimRight(r.RootDirectory, "/")+path, "/")
}

func (r *ReportKS3) getOptions() ks3.Options {
	return ks3.Options{
		SSE: r.Encrypt,
	}
}