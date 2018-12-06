FROM iron/go

ADD catalog_service_unix /

ENV PORT 8080
ENV GIN_MODE release

EXPOSE 8080

ENTRYPOINT ["/catalog_service_unix"]
