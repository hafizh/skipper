package run

import (
	"github.com/zalando/skipper/dispatch"
	"github.com/zalando/skipper/etcd"
	"github.com/zalando/skipper/filters"
	"github.com/zalando/skipper/proxy"
	"github.com/zalando/skipper/settings"
	"github.com/zalando/skipper/skipper"
	"log"
	"net/http"
)

const storageRoot = "/skipper"

// Run skipper. Expects address to listen on and one or more urls to find
// the etcd service at. If the flag 'insecure' is true, skipper will accept
// invalid TLS certificates from the backends.
func Run(address string, etcdUrls []string, insecure bool, customFilters ...skipper.FilterSpec) error {
	// create a client to etcd
	dataClient, err := etcd.Make(etcdUrls, storageRoot)
	if err != nil {
		return err
	}

	// create a filter registry with the available filter specs registered,
	// and register the custom filters
	registry := filters.RegisterDefault()
	registry.Add(customFilters...)

	// create a settings dispatcher instance
	// create a settings source
	// create the proxy instance
	dispatcher := dispatch.Make()
	settingsSource := settings.MakeSource(dataClient, registry, dispatcher)
	proxy := proxy.Make(settingsSource, insecure)

	// subscribe to new settings
	settingsChan := make(chan skipper.Settings)
	dispatcher.Subscribe(settingsChan)

	// start the http server
	log.Printf("listening on %v\n", address)
	return http.ListenAndServe(address, proxy)
}
