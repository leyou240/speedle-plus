package etcd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/leyou240/speedle-plus/api/pms"
	"github.com/leyou240/speedle-plus/pkg/errors"
	"github.com/leyou240/speedle-plus/pkg/store"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	StoreType = "etcd"

	//Following are keys of etcd store properties
	EtcdEndpointKey               = "EtcdEndpoint"
	EtcdKeyPrefixKey              = "EtcdKeyPrefix"
	EtcdTLSClientCertFileKey      = "EtcdTLSCertFile"
	EtcdTLSClientKeyFileKey       = "EtcdTLSKeyFile"
	EtcdTLSClientTrustedCAFileKey = "EtcdTLSTrustedCAFile"
	EtcdTLSAllowedCNKey           = "EtcdTLSAllowedCN"
	EtcdTLSServerNameKey          = "EtcdTLSServerName"
	EtcdTLSCRLFileKey             = "EtcdTLSCRLFile"
	EtcdTLSInsecureSkipVerifyKey  = "EtcdTLSInsecureSkipVerify"

	IsEmbeddedEtcdFlagName             = "etcdstore-isembedded"
	EmbeddedEtcdDataDirFlagName        = "etcdstore-embeddedDataDir"
	EtcdEndpointFlagName               = "etcdstore-endpoint"
	EtcdKeyPrefixFlagName              = "etcdstore-keyprefix"
	EtcdTLSClientCertFileFlagName      = "etcdstore-tls-cert"
	EtcdTLSClientKeyFileFlagName       = "etcdstore-tls-key"
	EtcdTLSClientTrustedCAFileFlagName = "etcdstore-tls-ca"
	EtcdTLSAllowedCNFlagName           = "etcdstore-tls-allowedCN"
	EtcdTLSServerNameFlagName          = "etcdstore-tls-serverName"
	EtcdTLSCRLFileFlagName             = "etcdstore-tls-crlFile"
	EtcdTLSInsecureSkipVerifyFlagName  = "etcdstore-tls-insecureSkipVerify"

	//default property values
	DefaultKeyPrefix           = "/speedle_ps/"
	DefaultEtcdStoreEndpoint   = "localhost:2379"
	DefaultEtcdStoreKeyPrefix  = "/speedle_ps/"
	DefaultEtcdStoreIsEmbedded = false
)

type Etcd3StoreBuilder struct{}

func (esb Etcd3StoreBuilder) NewStore(config map[string]interface{}) (pms.PolicyStoreManager, error) {
	keyPrefix, ok := config[EtcdKeyPrefixKey].(string)
	if !ok {
		keyPrefix = DefaultKeyPrefix
	}

	etcdStore := Store{}
	var etcd3ClientConf clientv3.Config
	etcdEndpoint, ok := config[EtcdEndpointKey].(string)
	if !ok {
		return nil, errors.New(errors.ConfigError, "configure item EtcdEndpoint is not found")
	}
	log.Debugf("new Etcd etcdStore: etcdEndpoint = %q, keyPrefix = %q\n", etcdEndpoint, keyPrefix)
	etcd3ClientConf.Endpoints = []string{etcdEndpoint}
	if strings.HasPrefix(etcdEndpoint, "https") {
		tlsInfo := transport.TLSInfo{}

		if certFile, ok := config[EtcdTLSClientCertFileKey].(string); ok {
			tlsInfo.CertFile = certFile
		}

		if keyFile, ok := config[EtcdTLSClientKeyFileKey].(string); ok {
			tlsInfo.KeyFile = keyFile
		}

		if trustedCAFile, ok := config[EtcdTLSClientTrustedCAFileKey].(string); ok {
			tlsInfo.TrustedCAFile = trustedCAFile
		}

		if cRLFile, ok := config[EtcdTLSCRLFileKey].(string); ok {
			tlsInfo.CRLFile = cRLFile
		}

		if val, ok := config[EtcdTLSInsecureSkipVerifyKey]; ok {
			var err error
			tlsInfo.InsecureSkipVerify, err = convertValueToBool(val, EtcdTLSInsecureSkipVerifyKey)
			if err != nil {
				return nil, err
			}
		}

		if serverName, ok := config[EtcdTLSServerNameKey].(string); ok {
			tlsInfo.ServerName = serverName
		}

		if allowedCN, ok := config[EtcdTLSAllowedCNKey].(string); ok {
			tlsInfo.AllowedCN = allowedCN
		}

		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}

		fmt.Printf("tlsInfo: %v\n", tlsInfo)
		etcd3ClientConf.TLS = tlsConfig
	}
	cli, err := clientv3.New(etcd3ClientConf)
	if err != nil {
		return nil, errors.Wrap(err, errors.StoreError, "failed to connect to etcd server")
	}
	etcdStore.client = cli
	etcdStore.Config = &etcd3ClientConf
	etcdStore.KeyPrefix = keyPrefix
	fmt.Println("Etcd etcdStore...")
	return &etcdStore, nil
}

func convertValueToBool(val interface{}, keyName string) (bool, error) {
	switch x := val.(type) {
	case bool:
		boolValule := val.(bool)
		return boolValule, nil
	case string:
		boolValule, err := strconv.ParseBool(val.(string))
		if err != nil {
			return false, errors.Wrapf(err, errors.ConfigError, "failed to convert configure %q", keyName)
		}
		return boolValule, nil
	default:
		return false, errors.Errorf(errors.ConfigError, "unsupported data type %T for configuration item %q", x, keyName)
	}
}

func (esb Etcd3StoreBuilder) GetStoreParams() map[string]string {
	return map[string]string{
		EtcdEndpointFlagName:               EtcdEndpointKey,
		EtcdKeyPrefixFlagName:              EtcdKeyPrefixKey,
		EtcdTLSClientCertFileFlagName:      EtcdTLSClientCertFileKey,
		EtcdTLSClientKeyFileFlagName:       EtcdTLSClientKeyFileKey,
		EtcdTLSClientTrustedCAFileFlagName: EtcdTLSClientTrustedCAFileKey,
		EtcdTLSAllowedCNFlagName:           EtcdTLSAllowedCNKey,
		EtcdTLSServerNameFlagName:          EtcdTLSServerNameKey,
		EtcdTLSCRLFileFlagName:             EtcdTLSCRLFileKey,
		EtcdTLSInsecureSkipVerifyFlagName:  EtcdTLSInsecureSkipVerifyKey,
	}

}

func init() {
	pflag.String(EtcdEndpointFlagName, DefaultEtcdStoreEndpoint, "Store config: endpoint of etcd store.")
	pflag.String(EtcdKeyPrefixFlagName, DefaultEtcdStoreKeyPrefix, "Store config: key prefix to store speedle policy data in etcd store.")
	pflag.Bool(IsEmbeddedEtcdFlagName, DefaultEtcdStoreIsEmbedded, "Store config: is embedded etcd store or not.")
	pflag.String(EmbeddedEtcdDataDirFlagName, "", "Store config: data dir for embedded etcd store.")
	pflag.String(EtcdTLSClientCertFileFlagName, "", "Store config: etcd x509 client cert.")
	pflag.String(EtcdTLSClientKeyFileFlagName, "", "Store config: etcd x509 client key.")
	pflag.String(EtcdTLSClientTrustedCAFileFlagName, "", "Store config: etcd x509 client CA cert.")
	pflag.String(EtcdTLSAllowedCNFlagName, "", "Store config: etcd x509 allowed CN.")
	pflag.String(EtcdTLSServerNameFlagName, "", "Store config: etcd x509 server name.")
	pflag.String(EtcdTLSCRLFileFlagName, "", "Store config: etcd x509 CRL file.")
	pflag.Bool(EtcdTLSInsecureSkipVerifyFlagName, false, "Store config: etcd x509 insecure skip verify.")

	store.Register(StoreType, Etcd3StoreBuilder{})
}
