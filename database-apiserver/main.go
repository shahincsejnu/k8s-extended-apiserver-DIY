package main

import(
	"fmt"
	"net"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"github.com/shahincsejnu/k8s-extended-apiserver-DIY/lib/certstore"
	"github.com/shahincsejnu/k8s-extended-apiserver-DIY/lib/server"
	"k8s.io/client-go/util/cert"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "oka")
}

func main() {
	fs := afero.NewOsFs()
	store, err := certstore.NewCertStore(fs, "/tmp/k8s-extended-apiserver-DIY")
	if err != nil {
		log.Fatalln(err)
	}
	err = store.NewCA("database")
	if err != nil {
		log.Fatalln(err)
	}

	serverCert, serverKey, err := store.NewServerCertPair(cert.AltNames{
		IPs: []net.IP{net.ParseIP("127.0.0.2")},
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = store.Write("tls", serverCert, serverKey)
	if err != nil {
		log.Fatalln(err)
	}

	clientCert, clientKey, err := store.NewClientCertPair(cert.AltNames{
		DNSNames: []string{"jane"},
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = store.Write("jane", clientCert, clientKey)
	if err != nil {
		log.Fatalln(err)
	}

	cfg := server.Config{
		Address: "127.0.0.2:8443",
		CACertFiles: []string{
			store.CertFile("ca"),
		},
		CertFile: store.CertFile("tls"),
		KeyFile:  store.KeyFile("tls"),
	}
	srv := server.NewGenericServer(cfg)

	r := mux.NewRouter()
	r.HandleFunc("/database/{resource}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Resource: %v\n", vars["resource"])
	})
	r.HandleFunc("/", handler)
	srv.ListenAndServe(r)
}