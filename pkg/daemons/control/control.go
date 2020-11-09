package control

import (
	"net"
	"os"
	"text/template"
	"crypto/x509"
	"github.com/xiaods/k8e/lib/tcplistener/cert"
)

var (
	localhostIP        = net.ParseIP("127.0.0.1")
	requestHeaderCN    = "system:auth-proxy"
	kubeconfigTemplate = template.Must(template.New("kubeconfig").Parse(`apiVersion: v1
clusters:
- cluster:
    server: {{.URL}}
    certificate-authority: {{.CACert}}
  name: local
contexts:
- context:
    cluster: local
    namespace: default
    user: user
  name: Default
current-context: Default
kind: Config
preferences: {}
users:
- name: user
  user:
    client-certificate: {{.ClientCert}}
    client-key: {{.ClientKey}}
`))
)

func KubeConfig(dest, url, caCert, clientCert, clientKey string) error {
	data := struct {
		URL        string  
		CACert     string
		ClientCert string
		ClientKey  string
	}{
		URL:        url,
		CACert:     caCert,
		ClientCert: clientCert,
		ClientKey:  clientKey,
	}

	output, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer output.Close()

	return kubeconfigTemplate.Execute(output, &data)    
}


type signedCertFactory = func(commonName string, organization []string, certFile, keyFile string) (bool, error)

func GetSigningCertFactory(regen bool, altNames *cert.AltNames, extKeyUsage []x509.ExtKeyUsage, caCertFile, caKeyFile string) signedCertFactory {
	return func(commonName string, organization []string, certFile, keyFile string) (bool, error) {
		return cert.CreateClientCertKey(regen, commonName, organization, altNames, extKeyUsage, caCertFile, caKeyFile, certFile, keyFile)  
	}  
}




