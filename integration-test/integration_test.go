package integration_test

import (
	. "github.com/Eun/go-hit"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	// Attempts connection
	host       = "app:8081"
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

func TestHTTPDoCommon(t *testing.T) {
	// call clear
	Test(t,
		Delete(basePath+"/statistics/clear"),
		Send().Headers("Content-Type").Add("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().Bytes().Equal([]byte{}),
	)

	// add some row
	Test(t,
		Post(basePath+"/statistics/save"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(`{ "date": "1234-12-12" }`),
		Expect().Status().Equal(http.StatusAccepted),
		Expect().Body().Bytes().Equal([]byte{}),
	)

	// make sure that row is there
	Test(t,
		Post(basePath+"/statistics/get"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(`{ "from": "1234-12-12", "to": "1234-12-12" }`),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Equal([]map[string]interface{}{{"date": "1234-12-12", "views": 0, "clicks": 0, "cost": "0", "cpc": 0, "cpm": 0}}),
	)

	// call clear
	Test(t,
		Delete(basePath+"/statistics/clear"),
		Send().Headers("Content-Type").Add("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().Bytes().Equal([]byte{}),
	)

	// expected empty array
	Test(t,
		Post(basePath+"/statistics/get"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(`{ "from": "1234-12-12", "to": "1234-12-12" }`),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().Equal([]map[string]interface{}{}),
	)
}

func TestHTTPDoSave(t *testing.T) {
	doTestCase := func(reqString string, expectStatus int64) {
		Test(t,
			Post(basePath+"/statistics/save"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().String(reqString),
			Expect().Status().Equal(expectStatus),
		)
	}

	input := []struct {
		reqString    string
		expectStatus int64
	}{
		{`{ "date": "2022-12-31", "views": 9, "cost": "2.00", "clicks": 1  }`, 202},
		{`{ "date": "2022-12-31" }`, 202},
		{`{ "date": "2022-12-31", "cost": "2.00" }`, 202},
		{`{ "date": "2021-12-31", "views": 9, "clicks": 1, "cost": "3.00" }`, 202},
		{`{ "date": "2022-11-30", "views": 9, "clicks": 1, "cost": "2.0" }`, 202},
		{`{ "date": "2022-12-30", "views": 9, "clicks": 1, "cost": "1" }`, 202},
		{`{ "date": "2023-12-30", "views": 9, "clicks": 1, "cost": "1.23" }`, 202},
		{`{ "date": "2013-12-30", "views": 9, "clicks": 1, "cost": "0.0" }`, 202},
		{`{ "date": "2013-12-30", "views": 9, "clicks": 1, "cost": "0.00" }`, 202},
		{``, 400},
		{`{ "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "202a-12-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-13-31", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-13-31", "views": 9, "clicks": 1, "cost": "2,00" }`, 400},
		{`{ "date": "2022-12-32", "views": 9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": "9", "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": -9, "clicks": 1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": "1", "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": -1, "cost": "2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": 2.00 }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "-2.00" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2.001" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2.a0" }`, 400},
		{`{ "date": "2022-12-31", "views": 9, "clicks": 1, "cost": "2." }`, 400},
	}

	for _, row := range input {
		doTestCase(row.reqString, row.expectStatus)
	}
}

func TestHTTPDoGet(t *testing.T) {
	// call clear
	Test(t,
		Delete(basePath+"/statistics/clear"),
		Send().Headers("Content-Type").Add("application/json"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().Bytes().Equal([]byte{}),
	)

	doSaveCase := func(reqString string, expectStatus int64) {
		Test(t,
			Post(basePath+"/statistics/save"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().String(reqString),
			Expect().Status().Equal(expectStatus),
		)
	}

	input := []struct {
		reqString    string
		expectStatus int64
	}{
		{`{ "date": "1111-11-11", "views": 15, "clicks": 0, "cost": "15" }`, 202},
		{`{ "date": "2222-12-12", "views": 2, "clicks": 4, "cost": "0" }`, 202},
		{`{ "date": "2222-12-12", "views": 1, "clicks": 1, "cost": "0.01" }`, 202},
		{`{ "date": "2222-12-13", "views": 9, "clicks": 9, "cost": "0" }`, 202},
		{`{ "date": "2222-12-12", "views": 1, "clicks": 1, "cost": "2" }`, 202},
	}

	for _, row := range input {
		doSaveCase(row.reqString, row.expectStatus)
	}

	var getCases = []struct {
		reqString    string
		expectStatus int64
		expectJSON   interface{}
	}{
		{`{ "from": "0000-00-00", "to": "0000-00-00" }`, 400, map[string]interface{}{"error": "invalid request body"}},
		{`{ "from": "0001-01-01", "to": "0001-13-01" }`, 400, map[string]interface{}{"error": "invalid request body"}},
		{`{ "from": "0001-01-01", "to": "0001-12-1" }`, 400, map[string]interface{}{"error": "invalid request body"}},
		{`{ "from": "0001-01-01", "to": "0001-01-01" }`, 200, []interface{}{}},
		{`{ "from": "0001-01-01", "to": "0001-02-01" }`, 200, []interface{}{}},

		{`{ "from": "2222-12-12", "to": "2222-12-12" }`, 200,
			[]map[string]interface{}{{"date": "2222-12-12", "views": 4, "clicks": 6, "cost": "2.01", "cpc": 0.34, "cpm": 502.5}},
		},

		{`{ "from": "1111-11-11", "to": "2222-12-13" }`, 200,
			[]map[string]interface{}{
				{"date": "1111-11-11", "views": 15, "clicks": 0, "cost": "15", "cpc": 0, "cpm": 1000},
				{"date": "2222-12-12", "views": 4, "clicks": 6, "cost": "2.01", "cpc": 0.34, "cpm": 502.5},
				{"date": "2222-12-13", "views": 9, "clicks": 9, "cost": "0", "cpc": 0, "cpm": 0}},
		},

		{`{ "from": "1111-11-11", "to": "2222-12-13", "order": "Views" }`, 200,
			[]map[string]interface{}{
				{"date": "2222-12-12", "views": 4, "clicks": 6, "cost": "2.01", "cpc": 0.34, "cpm": 502.5},
				{"date": "2222-12-13", "views": 9, "clicks": 9, "cost": "0", "cpc": 0, "cpm": 0},
				{"date": "1111-11-11", "views": 15, "clicks": 0, "cost": "15", "cpc": 0, "cpm": 1000}},
		},
	}
	doTestCase := func(reqString string, expectStatus int64, expectJSON interface{}) {
		Test(t,
			Post(basePath+"/statistics/get"),
			Send().Headers("Content-Type").Add("application/json"),
			Send().Body().String(reqString),
			Expect().Status().Equal(expectStatus),
			Expect().Body().JSON().Equal(expectJSON),
		)
	}

	for _, testCase := range getCases {
		doTestCase(testCase.reqString, testCase.expectStatus, testCase.expectJSON)
	}
}
