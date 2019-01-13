package test

import (
	"github.com/golang/glog"
	"io/ioutil"
	"k8s.io/perf-tests/clusterloader2/pkg/measurement"
	"path"
	"time"
	"k8s.io/perf-tests/clusterloader2/api"
)

type ReportLocal struct {}

func NewReportLocal() *ReportLocal {
	return &ReportLocal{}
}

func (r *ReportLocal) ReportSummary(ctx Context, conf *api.Config, summary measurement.Summary) error {
	summaryText, err := summary.PrintSummary()
	if err != nil {
		return err
	}

	if ctx.GetClusterLoaderConfig().ReportDir == "" {
		glog.Infof("%v: %v", summary.SummaryName(), summaryText)
	} else {
		// TODO(krzysied): Remember to keep original filename style for backward compatibility.
		filePath := path.Join(ctx.GetClusterLoaderConfig().ReportDir, summary.SummaryName()+"_"+conf.Name+"_"+time.Now().Format(time.RFC3339)+".txt")
		if err := ioutil.WriteFile(filePath, []byte(summaryText), 0644); err != nil {
			return err
		}
	}
	return nil
}


