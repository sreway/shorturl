# Make cert and key

The self sign key and certificate can be generated with the following command:

```bash
 openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes -keyout server.key -out server.crt -subj "/C=RU/ST=SPB/L=SPB/O=YaPraktikum/OU=Cohort15/CN=shorturl.sreway.com/emailAddress=andrey.oleynik@sreway.com" -addext "subjectAltName=IP:127.0.0.1"
```
