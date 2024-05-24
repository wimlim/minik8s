docker run -d --privileged --name my-nginx-container -p 80:80 -v /root/minik8s/pkg/nginx:/etc/nginx/conf.d nginx
docker exec my-nginx-container apt-get update && apt-get install -y ipvsadm