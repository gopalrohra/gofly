package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gopalrohra/flyapi/log"
)

type E2ETest struct {
	Name         string
	APIRequest   Request
	Expectations Expectations
}

func makeRequest(r Request) map[string]interface{} {
	client := &http.Client{}
	url, reqBody := getURLAndBody(r)
	req, _ := http.NewRequest(r.Method, url, reqBody)
	for _, header := range r.Headers {
		req.Header.Add(getHeader(header))
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error while making the request: %v\n", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error while parsing response: %v\n", err.Error())
		return nil
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Errorf("Error while unmarshiling response: %v\n", err.Error())
		return nil
	}
	return response
}
func getHeader(header string) (key string, value string) {
	return strings.Split(header, ":")[0], strings.Split(header, ":")[0]
}
func getURLAndBody(r Request) (url string, body io.Reader) {
	if r.Method == "GET" {
		return r.GetURLWithQueryParams(), nil
	} else if r.Method == "POST" || r.Method == "PUT" {
		b, _ := json.Marshal(r.Params)
		return r.URL, bytes.NewBuffer(b)
	}
	return r.URL, nil
}
func Run(t *testing.T) {
	tests, err := parseTestCases()
	if err != nil {
		log.Fatalf("Error occured %v\n", err.Error())
	}
	for _, test := range tests {
		t.Run(test.Name, func(it *testing.T) {
			res := makeRequest(test.APIRequest)
			log.Debug(res)
			expect := ExpectationChecker{response: res, expectations: test.Expectations}
			if !expect.shouldHave(test.Expectations.ShouldHave) {
				it.Errorf("%s test failed with response: %v\n", test.Name, res)
			}
		})
	}
}
func parseTestCases() ([]E2ETest, error) {
	var testDir = os.Getenv("testDir")
	var testCaseFile = os.Getenv("testCaseFile")
	log.Debugf("Value of testDir: %v and value of testCaseFile: %v\n", testDir, testCaseFile)
	var tests []E2ETest
	if testCaseFile != "" {
		testCasesFile := fmt.Sprintf("%s/%s.json", testDir, testCaseFile)
		tests, err := readTestCases(testCasesFile)
		log.Infof("Found %v test cases.\n", len(tests))
		return tests, err
	} else {
		fileInfo, err := ioutil.ReadDir(testDir)
		if err != nil {
			return tests, nil
		}
		for _, file := range fileInfo {
			result, err := readTestCases(testDir + file.Name())
			if err != nil {
				return tests, err
			}
			tests = append(tests, result...)
		}
		log.Infof("Found %v test cases.\n", len(tests))
		return tests, nil
	}
}
func readTestCases(testCasesFile string) ([]E2ETest, error) {
	var tests []E2ETest
	data, err := ioutil.ReadFile(testCasesFile)
	if err != nil {
		return tests, err
	}
	err = json.Unmarshal(data, &tests)
	return tests, err
}
