package utils

import (
	"encoding/json"
	"io"
	"os"

	"github.com/leyou240/speedle-plus/api/pms"
	"github.com/leyou240/speedle-plus/pkg/errors"
)

func ReadFilePolicyStore(policyStoreFile string) (*pms.PolicyStore, error) {
	file, err := os.Open(policyStoreFile)
	defer file.Close()
	if err != nil {
		return nil, errors.Wrapf(err, errors.StoreError, "unable to open file %q", policyStoreFile)
	}
	ret, err := readPolicyStore(file)
	return ret, err
}

func readPolicyStore(reader io.Reader) (*pms.PolicyStore, error) {
	decoder := json.NewDecoder(reader)
	var policyStore pms.PolicyStore
	if err := decoder.Decode(&policyStore); err != nil {
		return nil, errors.Wrap(err, errors.SerializationError, "unable to decode poilcy store")
	}
	return &policyStore, nil
}
