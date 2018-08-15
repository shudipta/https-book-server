FROM ubuntu
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY https-book-server /bin/https-book-server
EXPOSE 8443

CMD [ "/bin/https-book-server", "/var/book-serving-cert/server.key", "/var/book-serving-cert/server.crt", "/var/book-serving-cert/ca.crt" ]

# CMD [ "/bin/bash", "-ec", "while :; do echo '.'; sleep 5 ; done" ]