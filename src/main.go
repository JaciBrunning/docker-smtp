package main

import (
	"log"
	"os"
	"os/exec"
	"time"
	"context"
	"bytes"
	"io/ioutil"
	"compress/gzip"
	"encoding/json"
	"encoding/base64"
	"go.etcd.io/etcd/clientv3"
)

type JsonResponse struct {
	DomainC DomainCert `json:"DomainsCertificate"`
}

type DomainCert struct {
	Certs []CertEntry `json:"Certs"`
}

type CertEntry struct {
	Dms Domains `json:"Domains"`
	Cert Certificate `json:"Certificate"`
}

type Domains struct {
	Main string `json:"Main"`
	SANs []string `json:"SANs"`
}

type Certificate struct {
	Dom string `json:"Domain"`
	Priv string `json:"PrivateKey"`
	Cert string `json:"Certificate"`
}

// Args: domain key {endpoints}
func main() {
	key := os.Args[2]
	cfg := clientv3.Config{
		Endpoints: os.Args[3:],
		DialTimeout: 5 * time.Second,
	}

	c, err := clientv3.New(cfg)
	for err != nil {
		log.Panic("Client error: ")
		log.Panic(err)

		time.Sleep(time.Second)
		c, err = clientv3.New(cfg)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	last_value := []byte{}

	v, err := c.Get(ctx, key)
	for true {
		for _, ev := range v.Kvs {
			if (string(ev.Key) == key) {
				current_value := ev.Value
				if (!bytes.Equal(current_value, last_value)) {
					decode_certs(ungz(current_value))
				}
				last_value = current_value
			}
		}
		v, err = c.Get(ctx, key)
		time.Sleep(time.Second)
	}
}

func decode_certs(value []byte) {
	var response JsonResponse
	log.Print("New Certs!")
	err := json.Unmarshal(value, &response)
	if err != nil {
		log.Print(string(value))
		log.Fatal(err)
	}

	for _, cert_entry := range response.DomainC.Certs {
		if cert_entry.Dms.Main == os.Args[1] || string_in_slice(cert_entry.Dms.SANs, os.Args[1]) {
			cert_value, _  := base64.StdEncoding.DecodeString(cert_entry.Cert.Cert)
			priv_value, _ := base64.StdEncoding.DecodeString(cert_entry.Cert.Priv)

			// Write to file		
			cert_name := cert_entry.Dms.Main
			os.MkdirAll("/certs", os.ModePerm)
			ioutil.WriteFile("/certs/" + cert_name + ".crt", cert_value, 0644)
			ioutil.WriteFile("/certs/" + cert_name + ".key", priv_value, 0644)
			log.Print("Wrote cert: " + cert_name)

			update_postfix(cert_name)
		}
	}
}

func update_postfix(cert_name string) {
	postconf("-e", "smtpd_tls_cert_file=/certs/" + cert_name + ".crt")
	postconf("-e", "smtpd_tls_key_file=/certs/" + cert_name + ".key")

	postconf("-M", "submission/inet=submission   inet   n   -   n   -   -   smtpd")
	postconf("-P", "submission/inet/syslog_name=postfix/submission")
	postconf("-P", "submission/inet/smtpd_tls_security_level=encrypt")
	postconf("-P", "submission/inet/smtpd_sasl_auth_enable=yes")
	postconf("-P", "submission/inet/milter_macro_daemon_name=ORIGINATING")
	postconf("-P", "submission/inet/smtpd_recipient_restrictions=permit_sasl_authenticated,permit_mynetworks,reject_unauth_destination")
	postfix_reload()
	log.Print("Updated postconf for certificates!")
}

func postconf(args ...string) {
	cmd := exec.Command("postconf", args...)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func postfix_reload() {
	cmd := exec.Command("postfix", "reload")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func string_in_slice(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ungz(data []byte) []byte {
	gr, _ := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()

	data, _ = ioutil.ReadAll(gr)
	return data
}