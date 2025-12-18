#!/bin/bash
set -e

COMPOSE_FILE=docker-compose.test.yml
PROJECT_NAME=chat_service_test

docker compose -f $COMPOSE_FILE -p $PROJECT_NAME build chat_service_app_test

docker compose -f $COMPOSE_FILE -p $PROJECT_NAME up -d \
  chat_service_mongo_test \
  chat_service_app_test

# Mongo ready wait
echo "Waiting for MongoDB..."
until docker compose -f $COMPOSE_FILE -p $PROJECT_NAME exec -T \
  chat_service_mongo_test \
  mongosh --eval "db.runCommand({ ping: 1 })" &>/dev/null
do
  sleep 0.5
done

docker compose -f $COMPOSE_FILE -p $PROJECT_NAME exec \
  chat_service_app_test \
  sh -c "\
PKGS=\$(go list ./... | grep -v test_helper) && \
go test \$PKGS -coverprofile=coverage.out && \
go tool cover -html=coverage.out -o /mount/coverage.html \
"

docker compose -f $COMPOSE_FILE -p $PROJECT_NAME down -v

echo "カバレッジレポートをブラウザで表示しますか？ (y/n)"
read -r answer
if [[ "$answer" == "y" || "$answer" == "Y" ]]; then
  if command -v xdg-open &> /dev/null; then
    xdg-open ./mount/coverage.html
  elif command -v open &> /dev/null; then
    open ./mount/coverage.html
  else
    echo "ブラウザを開くコマンドが見つかりません。ブラウザで ./mount/coverage.html を手動で開いてください。"
  fi
else
  echo "カバレッジレポートは ./mount/coverage.html に保存されました。"
fi

