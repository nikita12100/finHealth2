# finHealth2

## Roadamp

1. стата по облигам
2. сохранять файл в бд для постанализа?
3. p/e, p/s как спидометр
4. treemap по ценам от покупки
5. pie for total balance
6. sankey for buys

### FIXME

- attempt to write a readonly database
- читать инфу об экспирации облиги
- CNYRUB_TOM\ поправить кол-во (покупка облиг за юани не уменьшает кол-во бумаг)
- заменить jsonb на операциях

### build

docker buildx build --platform linux/amd64 -f Dockerfile.build -t botbuilder .
docker buildx build --platform linux/amd64 -f Dockerfile_server.build -t botbuilder_server .

docker ps -a

docker cp c765d5a1ec20:app/build/bot_app ./bot_app
