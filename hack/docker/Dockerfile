FROM ubuntu
COPY https-book-server /bin/https-book-server
EXPOSE 10010
CMD [ "/bin/https-book-server", "/var/book-serving-cert/srv.key", "/var/book-serving-cert/srv.crt", "/var/book-serving-cert/ca.crt" ]