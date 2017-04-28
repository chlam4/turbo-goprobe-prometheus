# turbo-goprobe-prometheus

This is a GO SDK probe developed for the Turbonomic Operations Manager.  This  probe aims to discover applications and
nodes from [Prometheus](https://prometheus.io/).  It is being developed actively and please expect a lot of changes.

As of currently, this probe supports:
* Creating Application entities based on the Prometheus [webdriver](https://github.com/mattbostock/webdriver_exporter)
and the [mysql](https://github.com/prometheus/mysqld_exporter) exporters.  More will be gradually added in the future.
* Collecting web app response time and mysql transaction data.  More will be gradually added in the future.
* Stitching the discovered Application entities with their underlying Virtual Machine entities, provided that they are
discovered by the Turbo OpsMgr.

To try it out:
0. Prerequisites:
  * Install your Turbonomic OpsMgr.  The probe as of currently has been tested against version 5.9.
  * Install your Prometheus server and supported exporters (as listed above).
1. Configuration
  * Customize turbo-server-conf.json to point to your Turbo OpsMgr instance.
  * Customize target-conf.json to point to your Prometheus server.
2. Run `go install ./...` to build and install.
3. Start the probe: `./turbo-goprobe-prometheus`
4. Confirm in your OpsMgr that a target has been created as specified in `target-conf.json`.
5. Browse the Turbo UI for:
  * The discovered Application entities
  * The transaction metric graph of your mysql instances; response time requires configuration changes to be exposed
  * The relationship between the Application and its underlying Virtual Machine - you may want to add a VC target for
  example in your Turbo OpsMgr to discover the corresponding VM.
