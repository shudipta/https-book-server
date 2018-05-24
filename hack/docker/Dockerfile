FROM alpine

RUN set -x \
  && apk add --update --no-cache ca-certificates

COPY https-book-server /bin/https-book-server

EXPOSE 8443

USER nobody:nobody

CMD [ "/bin/https-book-server", "/var/book-serving-cert/srv.key", "/var/book-serving-cert/srv.crt", "/var/book-serving-cert/ca.crt" ]

#ENTRYPOINT ["scanner"]

#FROM ubuntu
#RUN apt-get update
#RUN apt-get install -y ca-certificates
#COPY https-book-server /bin/https-book-server
#EXPOSE 8443
#CMD [ "/bin/https-book-server", "/var/book-serving-cert/srv.key", "/var/book-serving-cert/srv.crt", "/var/book-serving-cert/ca.crt" ]