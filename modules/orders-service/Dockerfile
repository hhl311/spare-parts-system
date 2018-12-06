FROM scratch

ADD orders_service_unix /

ENV PORT 8080
ENV GIN_MODE release

EXPOSE 8080

CMD ["/orders_service_unix"]
