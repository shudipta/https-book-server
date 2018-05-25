#FROM ubuntu
#FROM alpine

#RUN set -x \
#  && apk add --update --no-cache ca-certificates

#COPY https-client /bin/https-client

#EXPOSE 8444

#USER nobody:nobody

#CMD [ "/bin/https-client", "/var/client-serving-cert/client.key", "/var/client-serving-cert/client.crt", "/var/client-serving-cert/ca.crt" ]
# CMD [ "/bin/https-book-server", "/var/book-serving-cert/server.key", "/var/book-serving-cert/server.crt", "/var/book-serving-cert/ca.crt" ]

FROM ubuntu
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY https-client /bin/https-client
EXPOSE 8444
CMD [ "/bin/https-client", "/var/client-serving-cert/client.key", "/var/client-serving-cert/client.crt", "/var/client-serving-cert/ca.crt" ]
