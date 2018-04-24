# FROM ubuntu
FROM ubuntu
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY https-book-server /bin/https-book-server
EXPOSE 8443
CMD [ "/bin/https-book-server", "/var/book-serving-cert/srv.key", "/var/book-serving-cert/srv.crt", "/var/book-serving-cert/ca.crt" ]