========== VPN setup ===============
Скачать и установить WireGuard клиента с сайта https://www.wireguard.com/install/
Запустить клиента (файлы конфигурации предоставил)
Open WireGuard, click on "Import tunnel from file," and select the WireGuard compressed file or .conf file.
After import, click on "Activate" to connect to the WireGuard VPN server.
Проверка: ping 10.200.1.2
Для windows на англ языке инструкция: https://www.tp-link.com/kz/support/faq/3989/

========= Jump servers =============
 С этих ВМ можно подключится к серверам Казпочты
jump vm1-linux - 10.200.1.2  root:Zxcv159*  (SSH)
jump vm2-windows -10.200.1.3 администратор:Zxcv159*  (RDP)

Учетные записи Qazpost:
Логин: g.rakhmanov
Пароль: M8#Qe!2$ZrA9xKp

Логин: a.shukurov
Пароль: T4@Wm$P7!KXz#2e

Логин: h.gaparov
Пароль: 9R!$ZpK@eM2xQ7A

============= Qazpost VM list =======================
Server List (Qazpost):
NGINX =>                        		baraiq-p-api01.kazpost.kz  172.30.75.69
Kafka + Harbor =>              	baraiq-p-w01.kazpost.kz    172.30.75.78
Workflow svc =>                 	baraiq-p-w02.kazpost.kz    172.30.75.85
TAF svc + consul + casdoor=>     	baraiq-p-dbp01.kazpost.kz  172.30.75.91
PGSQL =>                        		baraiq-p-dbpr01.kazpost.kz 172.30.75.92
Mongo DB =>                     	baraiq-p-dbpr02.kazpost.kz 172.30.75.93
REDIS =>                        		baraiq-p-cch01.kazpost.kz  172.30.75.94
NATS JetStream =>              	baraiq-p-que01.kazpost.kz  172.30.75.97

dokcer-compose.yml файлы для сервисов лежат в директории /opt
например: /opt/kafka или /opt/taf




============= Обновление образов ===================
# build images on local PC 
docker build -f build/Processing.dockerfile --build-arg SERVICE=gatewaysvc --build-arg SOURCE_DIR=internal/processing/gatewaysvc -t 172.30.75.78:9080/taf/gatewaysvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=storagesvc --build-arg SOURCE_DIR=internal/processing/storagesvc -t 172.30.75.78:9080/taf/storagesvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=amlsvc --build-arg SOURCE_DIR=internal/processing/SVCs/amlsvc -t 172.30.75.78:9080/taf/amlsvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=lstsvc --build-arg SOURCE_DIR=internal/processing/SVCs/lstsvc -t 172.30.75.78:9080/taf/lstsvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=fcsvc --build-arg SOURCE_DIR=internal/processing/SVCs/fcsvc -t 172.30.75.78:9080/taf/fcsvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=finalizersvc --build-arg SOURCE_DIR=internal/processing/finalizersvc -t 172.30.75.78:9080/taf/finalizersvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=retrysvc --build-arg SOURCE_DIR=internal/processing/retrysvc -t 172.30.75.78:9080/taf/retrysvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=monitoringsvc --build-arg SOURCE_DIR=internal/monitoring -t 172.30.75.78:9080/taf/monitoringsvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=adminsvc  --build-arg SOURCE_DIR=internal/admin -t 172.30.75.78:9080/taf/adminsvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=notificationsvc --build-arg SOURCE_DIR=internal/notification -t 172.30.75.78:9080/taf/notificationsvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=enrichersvc --build-arg SOURCE_DIR=internal/enricher -t 172.30.75.78:9080/taf/enrichersvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=configurersvc --build-arg SOURCE_DIR=internal/configurersvc -t 172.30.75.78:9080/taf/configurersvc:latest .
docker build -f build/Processing.dockerfile --build-arg SERVICE=kycsvc  --build-arg SOURCE_DIR=internal/kyc -t 172.30.75.78:9080/taf/kycsvc:latest .

docker build -t 172.30.75.78:9080/taf/incidentssvc:latest .
docker build -t 172.30.75.78:9080/taf/taf-mlsvc:latest .

docker build -t 172.30.75.78:9080/taf/workflow-worker:latest -f Dockerfile.worker .
docker build -t 172.30.75.78:9080/taf/workflow-registry:latest -f Dockerfile.registry .
docker build -t 172.30.75.78:9080/taf/workflow-executor:latest -f Dockerfile.executor .

docker build -t 172.30.75.78:9080/taf/frontend:latest . 

# save images on local PC
docker save -o gatewaysvc_latest.tar 172.30.75.78:9080/taf/gatewaysvc:latest
docker save -o storagesvc_latest.tar 172.30.75.78:9080/taf/storagesvc:latest
docker save -o amlsvc_latest.tar 172.30.75.78:9080/taf/amlsvc:latest
docker save -o lstsvc_latest.tar 172.30.75.78:9080/taf/lstsvc:latest
docker save -o fcsvc_latest.tar 172.30.75.78:9080/taf/fcsvc:latest
docker save -o finalizersvc_latest.tar 172.30.75.78:9080/taf/finalizersvc:latest
docker save -o retrysvc_latest.tar 172.30.75.78:9080/taf/retrysvc:latest
docker save -o monitoringsvc_latest.tar 172.30.75.78:9080/taf/monitoringsvc:latest
docker save -o adminsvc_latest.tar 172.30.75.78:9080/taf/adminsvc:latest
docker save -o notificationsvc_latest.tar 172.30.75.78:9080/taf/notificationsvc:latest
docker save -o enrichersvc_latest.tar 172.30.75.78:9080/taf/enrichersvc:latest
docker save -o configurersvc_latest.tar 172.30.75.78:9080/taf/configurersvc:latest
docker save -o kycsvc_latest.tar 172.30.75.78:9080/taf/kycsvc:latest
docker save -o incidentssvc_latest.tar 172.30.75.78:9080/taf/incidentssvc:latest

docker save -o taf-mlsvc_latest.tar 172.30.75.78:9080/taf/taf-mlsvc:latest

docker save -o workflow-worker_latest.tar 172.30.75.78:9080/taf/workflow-worker:latest
docker save -o workflow-registry_latest.tar 172.30.75.78:9080/taf/workflow-registry:latest
docker save -o workflow-executor_latest.tar 172.30.75.78:9080/taf/workflow-executor:latest

docker save -o frontend_latest.tar 172.30.75.78:9080/taf/frontend:latest

#
#  copy saved images to /tmp directory on jump server 10.200.1.2
#


# single command copies all at once to Kazpost harbor VM 172.30.75.78 from jump server 10.200.1.2
scp /tmp/*_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
# =========OR ===============
# copy one service at a time to Kazpost harbor VM 172.30.75.78 from jump server 10.200.1.2
scp /tmp/gatewaysvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/storagesvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/amlsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/lstsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/fcsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/finalizersvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/retrysvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/monitoringsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/adminsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/notificationsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/enrichersvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/configurersvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/kycsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/incidentssvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/

scp /tmp/taf-mlsvc_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/

scp /tmp/workflow-worker_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/workflow-registry_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/
scp /tmp/workflow-executor_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/

scp /tmp/frontend_latest.tar zh.akhmetkarimov@172.30.75.78:/tmp/

# load commands Kazpost harbor VM 172.30.75.78
sudo docker load -i /tmp/gatewaysvc_latest.tar
sudo docker load -i /tmp/storagesvc_latest.tar
sudo docker load -i /tmp/amlsvc_latest.tar
sudo docker load -i /tmp/lstsvc_latest.tar
sudo docker load -i /tmp/fcsvc_latest.tar
sudo docker load -i /tmp/finalizersvc_latest.tar
sudo docker load -i /tmp/retrysvc_latest.tar
sudo docker load -i /tmp/monitoringsvc_latest.tar
sudo docker load -i /tmp/adminsvc_latest.tar
sudo docker load -i /tmp/notificationsvc_latest.tar
sudo docker load -i /tmp/enrichersvc_latest.tar
sudo docker load -i /tmp/configurersvc_latest.tar
sudo docker load -i /tmp/kycsvc_latest.tar
sudo docker load -i /tmp/incidentssvc_latest.tar

sudo docker load -i /tmp/taf-mlsvc_latest.tar

sudo docker load -i /tmp/workflow-worker_latest.tar
sudo docker load -i /tmp/workflow-registry_latest.tar
sudo docker load -i /tmp/workflow-executor_latest.tar

sudo docker load -i /tmp/frontend_latest.tar

# push commands Kazpost harbor VM 172.30.75.78
sudo docker push 172.30.75.78:9080/taf/gatewaysvc:latest
sudo docker push 172.30.75.78:9080/taf/storagesvc:latest
sudo docker push 172.30.75.78:9080/taf/amlsvc:latest
sudo docker push 172.30.75.78:9080/taf/lstsvc:latest
sudo docker push 172.30.75.78:9080/taf/fcsvc:latest
sudo docker push 172.30.75.78:9080/taf/finalizersvc:latest
sudo docker push 172.30.75.78:9080/taf/retrysvc:latest
sudo docker push 172.30.75.78:9080/taf/monitoringsvc:latest
sudo docker push 172.30.75.78:9080/taf/adminsvc:latest
sudo docker push 172.30.75.78:9080/taf/notificationsvc:latest
sudo docker push 172.30.75.78:9080/taf/enrichersvc:latest
sudo docker push 172.30.75.78:9080/taf/configurersvc:latest
sudo docker push 172.30.75.78:9080/taf/kycsvc:latest
sudo docker push 172.30.75.78:9080/taf/incidentssvc:latest

sudo docker push 172.30.75.78:9080/taf/taf-mlsvc:latest

sudo docker push 172.30.75.78:9080/taf/workflow-worker:latest
sudo docker push 172.30.75.78:9080/taf/workflow-registry:latest
sudo docker push 172.30.75.78:9080/taf/workflow-executor:latest

sudo docker push 172.30.75.78:9080/taf/frontend:latest
====== check on 172.30.75.78 HARBOR VM ==================
sudo docker compose images -a