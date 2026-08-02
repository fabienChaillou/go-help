# docker 
docker run -d \
  --name book-api-postgres \
  -e POSTGRES_USER=bookuser \
  -e POSTGRES_PASSWORD=bookpass \
  -e POSTGRES_DB=bookdb \
  -p 5432:5432 \
  postgres:16-alpine

# vérifier que le container tourne
docker ps

# se connecter en psql
docker exec -it book-api-postgres psql -U bookuser -d bookdb

# voir les logs
docker logs book-api-postgres

# arrêter / supprimer
docker stop book-api-postgres
docker rm book-api-postgres
