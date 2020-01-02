IPADDRESS=192.168.1.10

grafana-up:
	@echo Starting grafana...
	docker run --rm -d -p 3000:3000 --name=grafana -e "GF_SERVER_ROOT_URL=http://$(IPADDRESS):3000"  -e "GF_SECURITY_ADMIN_PASSWORD=secret" grafana/grafana
	sleep 2
	curl -v -s -k "http://admin:secret@$(IPADDRESS):3000/api/datasources" -X POST -H "Content-Type: application/json"  -d@grafana-dashboard/datasource.json > /dev/null 2>&1
	curl -v -s -k "http://admin:secret@$(IPADDRESS):3000/api/dashboards/db" -X POST -H "Content-Type: application/json"  -d@grafana-dashboard/dashboard.json > /dev/null 2>&1
	@echo "Grafana is up and running"

grafana-down:
	@echo "Stopping grafana"
	docker stop grafana

prometheus-up:
	@echo "Starting prometheus"
	sed -e 's/@IPADDRESS@/$(IPADDRESS)/g' $(PWD)/prometheus-config/prometheus.yaml > $(PWD)/prometheus-config/prometheus.yml
	docker run --rm -d --name prometheus -p 9090:9090 -v $(PWD)/prometheus-config:/etc/prometheus/ prom/prometheus

prometheus-down:
	@echo "Stopping prometheus"
	docker stop prometheus

build:
	@$(eval VERSION=`cat main.go |grep "version ="| cut -d"=" -f2 |sed -e 's/"//g' -e 's/ //g'`)
	docker build -t crawler-tags:$(VERSION) .