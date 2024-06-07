mkdir -p /srv/prometheus
cd /srv/prometheus


docker run -d \
  --name prometheus \
  -p 9090:9090 \
  -v /srv/prometheus:/srv/prometheus \
  -v /srv/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
