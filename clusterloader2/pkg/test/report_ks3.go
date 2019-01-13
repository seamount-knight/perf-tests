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

	bucket := path.Join(ctx.GetClusterLoaderConfig().TestJob, ctx.GetClusterLoaderConfig().BuildNumber)
	fileNmae := path.Join(summary.SummaryName()+"_"+conf.Name+".txt")
	return r.Put2KS3(bucket, fileNmae, summaryText)
}

func (r *ReportKS3) Put2KS3(path, fileName, summaryText string)error  {
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
