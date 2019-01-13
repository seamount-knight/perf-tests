package test

import (
	"fmt"
	"strings"
	"testing"
)
const listMax = 1000

func getKS3Config () *KS3Config {
	return &KS3Config{
		"AKLTVsI6SBANS3ypMan93Aon3w",
		"OJbvuwBHlzrui2/4fb+C+cv74d9WAMorittupvVYPGcFXRWW+3D0akhnFCCArOgaIQ==",
		"ks3-cn-beijing",
		"",
		true,
		false,
		"test333",
		"",
	}
}

func getReportKS3() (*ReportKS3, error) {
	return NewReportKS3(getKS3Config()), nil
}

func TestReportKS3_Put2KS3(t *testing.T) {
	reportKS3, err := getReportKS3()
	if err != nil {
		t.Fatal(err)
	}
	fileName := "kubemark-100/3/apiresponse_load.txt"
	summaryText := `{"test":"ok"}`
	err = reportKS3.Put2KS3(fileName, summaryText)

	if err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	reportKS3, err := getReportKS3()
	if err != nil {
		t.Fatal(err)
	}

	path := "kubemark-100/3"
	if path != "/" && path[len(path)-1] != '/' {
		path = path + "/"
	}
	prefix := ""
	if reportKS3.ks3Path("") == "" {
		prefix = "/"
	}

	listResponse, err := reportKS3.List(reportKS3.ks3Path(path), "/", "", listMax)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{}
	directories := []string{}

	for {
		for _, key := range listResponse.Contents {
			files = append(files, strings.Replace(key.Key, reportKS3.ks3Path(""), prefix, 1))
		}

		for _, commonPrefix := range listResponse.CommonPrefixes {
			directories = append(directories, strings.Replace(commonPrefix[0:len(commonPrefix)-1], reportKS3.ks3Path(""), prefix, 1))
		}

		if listResponse.IsTruncated {
			listResponse, err = reportKS3.Bucket.List(reportKS3.ks3Path(path), "/", listResponse.NextMarker, listMax)
			if err != nil {
				t.Fatal(err)
			}
		} else {
			break
		}
	}
	fmt.Println(files)
	fmt.Println(directories)
}

func TestGetObject(t *testing.T) {
	reportKS3, err := getReportKS3()
	if err != nil {
		t.Fatal(err)
	}
	data, err := reportKS3.Get("/kubemark-100/1/apiresponse_load.txt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}