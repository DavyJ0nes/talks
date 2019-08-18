FROM alpine:3.10
LABEL Name=feelgood-api
LABEL Author=awesome-corp

RUN mkdir -p /var/log/feelgood
ADD feelgood-api /

EXPOSE 8080
CMD ["./feelgood-api"]
