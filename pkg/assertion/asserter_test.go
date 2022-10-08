package assertion

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	adsapi "github.com/leyou240/speedle-plus/api/ads"
)

func TestAssertion(t *testing.T) {

	server := NewTestServer(t, nil)
	defer server.Close()

	asserter, err := getAsserter(server.URL+"/assert", t)

	t.Logf("load asserter: %v, err: %v", asserter, err)
	if asserter == nil || err != nil {
		t.Error("asserter is nil, err ", err)

	}

	ar, errAssert := asserter.AssertToken("testtoken", "WERCKER", "", nil)
	t.Logf("auth result: %v, err: %v", ar, errAssert)
	if errAssert != nil {
		t.Error("assertion, err ", errAssert)
	}

	if len(ar.Principals) != 1 {
		t.Fatalf("One principals should be returned but returned %d principals.\n", len(ar.Principals))
	}

	if ar.Principals[0].Name != "testUser" || ar.Principals[0].Type != adsapi.PRINCIPAL_TYPE_USER {
		t.Fatalf("returned user should be testUser, actually returned %v .\n", ar.Principals[0])
	}

	// no token
	ar, errAssert = asserter.AssertToken("", "", "", nil)
	t.Logf("auth result: %v, err: %v", ar, errAssert)
	if errAssert == nil {
		t.Fatalf("should failed with error")
	}

	// invalid token
	ar, errAssert = asserter.AssertToken("test-token", "", "", nil)
	t.Logf("auth result: %v, err: %v", ar, errAssert)
	if errAssert == nil {
		t.Fatalf("should failed with error")
	}
	if !strings.Contains(errAssert.Error(), strconv.Itoa(http.StatusBadRequest)) {
		t.Fatalf("should failed with error code: %d, actually failed with error code: %s.\n", http.StatusBadRequest, errAssert.Error())
	}

}

func getAsserter(endpoint string, t *testing.T) (TokenAsserter, error) {
	conf := &AsserterConfig{
		Endpoint: endpoint,
	}
	t.Logf("endpoint: %s", endpoint)
	return NewAsserter(conf, nil)

}
