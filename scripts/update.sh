docker rm -f freybot

docker rmi ghcr.io/tamper000/freybot
docker pull ghcr.io/tamper000/freybot:latest

docker run --restart=always -d --name freybot -v $(pwd)/config/config.yaml:/app/config/config.yaml -v $(pwd)/database/bot.db:/app/database/bot.db -p 8888:8888 ghcr.io/tamper000/freybot
