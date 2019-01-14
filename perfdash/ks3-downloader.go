package main

import (
	"fmt"
	"github.com/softlns/ks3-sdk-go/ks3"
	"k8s.io/kubernetes/test/e2e/perftype"
	"os"
	"strconv"
	"strings"
	"sync"
)

const kS3Bucket  =  "test333"
const listMax = 1000


type Ks3Downloader struct {
	DefaultBuildsCount int
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

//accessKeyID: "AKLTVsI6SBANS3ypMan93Aon3w"
//accessKeySecret: "OJbvuwBHlzrui2/4fb+C+cv74d9WAMorittupvVYPGcFXRWW+3D0akhnFCCArOgaIQ=="
//bucket: "test333"
//region: "ks3-cn-beijing"

func getDefaultKS3Config() *KS3Config {
	return &KS3Config {
		"AKLTVsI6SBANS3ypMan93Aon3w",
		"OJbvuwBHlzrui2/4fb+C+cv74d9WAMorittupvVYPGcFXRWW+3D0akhnFCCArOgaIQ==",
		"ks3-cn-beijing",
		"",
		true,
		false,
		kS3Bucket,
		"",
	}
}

func NewKs3Downloader(defaultBuildsCount int) *Ks3Downloader {
	config := getDefaultKS3Config()

	ks3obj, err := ks3.New(config.AccessKeyID, config.AccessKeySecret, config.Region, config.Secure,
		config.Internal, config.Endpoint)
	if err != nil {
		panic(err)
	}
	ks3.SetDebug(false)

	bucket := ks3obj.Bucket(config.Bucket)

	return &Ks3Downloader{
		DefaultBuildsCount: defaultBuildsCount,
		Bucket: bucket,
		Encrypt: true,
		RootDirectory: "",
	}
}

func (k *Ks3Downloader) getData() (JobToCategoryData, error) {
	data, err := k.Get("periodics.yaml")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	newJobs, err := parseProwConfig(data)
	if err == nil {
		TestConfig[kS3Bucket] = newJobs
	} else {
		fmt.Fprintf(os.Stderr, "Failed to refresh config: %v", err)
	}
	fmt.Printf("---------newJobs: %v\n", newJobs)

	fmt.Print("Getting Data from KS3...\n")
	result := make(JobToCategoryData)
	var resultLock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(TestConfig[kS3Bucket]))
	for job, tests := range TestConfig[kS3Bucket] {
		if tests.Prefix == "" {
			return result, fmt.Errorf("Invalid empty Prefix for job %s", job)
		}
		for categoryLabel, categoryMap := range tests.Descriptions {
			for testLabel := range categoryMap {
				resultLock.Lock()
				if _, found := result[tests.Prefix]; !found {
					result[tests.Prefix] = make(CategoryToMetricData)
				}
				if _, found := result[tests.Prefix][categoryLabel]; !found {
					result[tests.Prefix][categoryLabel] = make(MetricToBuildData)
				}
				if _, found := result[tests.Prefix][categoryLabel][testLabel]; found {
					return result, fmt.Errorf("Duplicate name %s for %s", testLabel, tests.Prefix)
				}
				result[tests.Prefix][categoryLabel][testLabel] = &BuildData{Job: job, Version: "", Builds: map[string][]perftype.DataItem{}}
				resultLock.Unlock()
			}
		}
		go k.getJobData(&wg, result, &resultLock, job, tests)
	}
	wg.Wait()
	return result, nil
}


func (k *Ks3Downloader) getJobData(wg *sync.WaitGroup, result JobToCategoryData, resultLock *sync.Mutex, job string, tests Tests) {
	defer wg.Done()
	lastBuildNo, err := k.GetLastBuildNumber(job)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Last build no for %v: %v\n", job, lastBuildNo)
	buildsToFetch := tests.BuildsCount
	if buildsToFetch < 1 {
		buildsToFetch = k.DefaultBuildsCount
	}
	fmt.Printf("Builds to fetch for %v: %v\n", job, buildsToFetch)

	for buildNumber := lastBuildNo; buildNumber > lastBuildNo-buildsToFetch && buildNumber > 0; buildNumber-- {
		fmt.Printf("Fetching %s build %v...\n", job, buildNumber)
		files, _ := k.ListBuckets(fmt.Sprintf("%s/%d", job, buildNumber))
		if len(files) == 0 {
			continue
		}

		for categoryLabel, categoryMap := range tests.Descriptions {
			for testLabel, testDescriptions := range categoryMap {
				for _, testDescription := range testDescriptions {
					testData, err := k.GetFileFromBucket(job, buildNumber, fmt.Sprintf("%s_%s.txt", testDescription.OutputFilePrefix, testDescription.Name))
					if err != nil {
						continue
					}
					fmt.Println(job, buildNumber, fmt.Sprintf("%s_%s.txt", testDescription.OutputFilePrefix, testDescription.Name))
					func() {
						resultLock.Lock()
						buildData := result[tests.Prefix][categoryLabel][testLabel]
						resultLock.Unlock()
						testDescription.Parser(testData, buildNumber, buildData)
					}()
					break
				}
			}
		}
	}
}

func (k *Ks3Downloader) GetLastBuildNumber(job string) (int, error) {
	filePath := fmt.Sprintf("%s/lastBuildNum", job)
	data, err := k.Bucket.Get(filePath)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	lastBuildNum, err := strconv.Atoi(strings.Trim(string(data), "\n"))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return lastBuildNum, nil
}


func (k *Ks3Downloader) ListBuckets(path string) ([]string, error) {
	files := []string{}

	if path != "/" && path[len(path)-1] != '/' {
		path = path + "/"
	}
	prefix := ""
	if k.ks3Path("") == "" {
		prefix = "/"
	}
	listResponse, err := k.Bucket.List(k.ks3Path(path), "/", "", listMax)
	if err != nil {
		return files, err
	}

	for {
		for _, key := range listResponse.Contents {
			files = append(files, strings.Replace(key.Key, k.ks3Path(""), prefix, 1))
		}

		if listResponse.IsTruncated {
			listResponse, err = k.Bucket.List(k.ks3Path(path), "/", listResponse.NextMarker, listMax)
			if err != nil {
				return files, err
			}
		} else {
			break
		}
	}
	fmt.Println(files)
	return files, nil
}


func (k *Ks3Downloader) GetFileFromBucket(job string, buildNumber int, fileName string) ([]byte, error) {
	filePath := fmt.Sprintf("%s/%d/%s", job, buildNumber, fileName)
	data, err := k.Bucket.Get(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}


func (k *Ks3Downloader) ks3Path(path string) string {
	return strings.TrimLeft(strings.TrimRight(k.RootDirectory, "/")+path, "/")
}

func (k *Ks3Downloader) getOptions() ks3.Options {
	return ks3.Options{
		SSE: k.Encrypt,
	}
}
