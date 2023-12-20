FROM postgres:15-bookworm

RUN mkdir -p /docker-entrypoint-initdb.d
COPY ./install_earthdistance.sql /docker-entrypoint-initdb.d/install_earthdistance.sql