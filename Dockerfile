FROM alpine:3.16.2

WORKDIR /api

COPY ./api ./api

#prepare migration 
COPY ./database /database

WORKDIR /api
RUN chmod +x ./api

CMD ["/api/api"]