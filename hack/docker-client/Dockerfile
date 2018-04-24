# FROM ubuntu
FROM ubuntu
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY https-client /bin/https-client
EXPOSE 8443
CMD [ "/bin/https-client", "/var/client-serving-cert/cl.key", "/var/client-serving-cert/cl.crt", "/var/client-serving-cert/ca.crt" ]