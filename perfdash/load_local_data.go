package main

import (
	"fmt"
	"io/ioutil"
	"k8s.io/contrib/test-utils/utils"
	"k8s.io/kubernetes/test/e2e/perftype"
	"os"
	"sync"
)

// LocalDownloader that gets data about test results from the local dir
type LocalDownloader struct {
}

// NewLocalDownloader creates a new LocalDownloader
func NewLocalDownloader() *LocalDownloader {
	return &LocalDownloader{}
}

// TODO(random-liu): Only download and update new data each time.
func (g *LocalDownloader) getData() (JobToCategoryData, error) {
	newJobs, err := getProwConfig()
	if err == nil {
		TestConfig[utils.KubekinsBucket] = newJobs
	} else {
		fmt.Fprintf(os.Stderr, "Failed to refresh config: %v", err)
	}
	//fmt.Printf("---------newJobs: %v\n", newJobs)
	fmt.Print("Getting Data from %s...\n", *localTestDir)
	result := make(JobToCategoryData)
	var resultLock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(TestConfig[utils.KubekinsBucket]))
	for job, tests := range TestConfig[utils.KubekinsBucket] {
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
		go g.getJobData(&wg, result, &resultLock, job, tests)
	}
	wg.Wait()
	return result, nil
}

func (g *LocalDownloader) getJobData(wg *sync.WaitGroup, result JobToCategoryData, resultLock *sync.Mutex, job string, tests Tests) {
	defer wg.Done()

	fmt.Printf("Fetching %s test data\n", job)

	for categoryLabel, categoryMap := range tests.Descriptions {
		for testLabel, testDescriptions := range categoryMap {
			for _, testDescription := range testDescriptions {
				fileStem := fmt.Sprintf("%v/%v_%v.json", *localTestDir, testDescription.OutputFilePrefix, testDescription.Name)
				testData, err := ioutil.ReadFile(fileStem)
				if err != nil {
					//fmt.Printf("read file %s err: %v\n", fileStem, err)
					continue
				} else {
					fmt.Printf("%s-%s-%s, read file %s\n",job, categoryLabel, testLabel, fileStem)
				}

				func() {
					resultLock.Lock()
					buildData := result[tests.Prefix][categoryLabel][testLabel]
					resultLock.Unlock()
					testDescription.Parser(testData, 0, buildData)
				}()
				break
			}
		}
	}

}

