IPADDRESS=192.168.1.10
CONSUMER_API_KEY="U3tuMr9txi4jtqtaZFHaKC2RO"
CONSUMER_API_SECRET="sKZwNIxDKj9kNVwO9K7oPncDqTExUWioY0yHxeVfY1j11xKfMc"
ACCESS_TOKEN_KEY="89549425-CYSCF9BEsM9BspXMt2sDrE3lNefBRSqMnp23t8psy"
ACCESS_TOKEN_SECRET="htBgTB7UfAyguykyN2smGCvIwLVwunEIq25v57376QI9j"
DB_PORT="2379"

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

run:
	@echo "Starting crawler-tags"
	docker run -d --name crawler-tags --rm -p 8080:8080 -e CONSUMER_API_KEY=$(CONSUMER_API_KEY) -e CONSUMER_API_SECRET=$(CONSUMER_API_SECRET) -e ACCESS_TOKEN_KEY=$(ACCESS_TOKEN_KEY) -e ACCESS_TOKEN_SECRET=$(ACCESS_TOKEN_SECRET) -e DB_HOST=$(IPADDRESS) -e DB_PORT="2379" crawler-tags:1.0.0-alpha

stop:
	@echo "Stopping crawler-tags"
	docker stop crawler-tags

etcd-up:
	@echo "Starting etcd3"
	docker run -d -p 2379:2379 -p 2380:2380 --name etcd3 --rm -v $(PWD)/data/etcd:/var/lib/etcd/data quay.io/coreos/etcd:v3.3.10 etcd --listen-client-urls http://0.0.0.0:2379 --initial-advertise-peer-urls http://$(IPADDRESS):2380 --listen-peer-urls http://0.0.0.0:2380 --initial-cluster-token etcd-cluster --initial-cluster etcd0=http://$(IPADDRESS):2380 --initial-cluster-state new --auto-compaction-retention=1 --advertise-client-urls=http://$(IPADDRESS):2379 --name=etcd0 --data-dir=/var/lib/etcd/data

etcd-down:
	@echo "Stopping etcd3"
	docker stop etcd3

es-up:
	@echo "Starting ES"
	docker run --name elasticsearch --rm -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.5.1

es-down:
	@echo "Stopping ES"
	docker stop elasticsearch

logstash-up:
	@echo "Starting Logstash"
	sed -e 's/@IPADDRESS@/$(IPADDRESS)/g' $(PWD)/logstash/logstash.config > $(PWD)/logstash/logstash.conf
	docker run -d -p 9600:9600 --name logstash --rm -v $(PWD)/logstash:/usr/share/logstash/config/ -v $(PWD)/logstash:/usr/share/logstash/pipeline/ -v /tmp/crawler-tags:/tmp/crawler-tags docker.elastic.co/logstash/logstash:6.4.2

logstash-down:
	@echo "Stopping logstash"
	docker stop logstash

kibana-up:
	@echo "Starting Kibana"
	docker run -d --rm -p 5601:5601 --name kibana -e ELASTICSEARCH_HOSTS=http://$(IPADDRESS):9200 docker.elastic.co/kibana/kibana:7.5.1

kibana-down:
	@echo "Stopping Kibana"
	docker stop kibana