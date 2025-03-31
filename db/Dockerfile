FROM postgres:latest

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin
ENV POSTGRES_DB=bank

EXPOSE 5432

COPY init.sql /docker-entrypoint-initdb.d/
