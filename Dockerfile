FROM golang:1.16rc1-alpine3.13 AS build
RUN apk update
RUN apk --no-cache add ca-certificates spamassassin=3.4.4-r2 vim wget unzip db db-dev make gcc perl-dev libc-dev
WORKDIR /mailpie
RUN wget https://github.com/da-coda/email-spam-ham-learning-dataset/archive/main.zip && unzip -q main.zip
RUN curl -L https://cpanmin.us | perl - App::cpanminus && cpanm -v DB_File::HASHINFO
RUN sa-update && sa-learn --ham /mailpie/email-spam-ham-learning-dataset-main/ham && sa-learn --spam /mailpie/email-spam-ham-learning-dataset-main/spam && sa-learn --backup > /mailpie/bayes.db
COPY ./ .
RUN go build github.com/da-coda/mailpie .

FROM alpine:3.13
EXPOSE 1025
EXPOSE 1143
EXPOSE 8000
EXPOSE 783
RUN apk update
RUN apk --no-cache add ca-certificates spamassassin=3.4.4-r2 vim db db-dev make gcc perl-dev libc-dev
WORKDIR /root/
COPY --from=build /mailpie/mailpie .
COPY --from=build /mailpie/bayes.db .
RUN curl -L https://cpanmin.us | perl - App::cpanminus && cpanm -v DB_File::HASHINFO
VOLUME /root
ENTRYPOINT sa-update && sa-learn --restore /root/bayes.db --sync && spamd -d -i 0.0.0.0 -L -A 172. && /root/mailpie -config /root/mailpie.yml -logLevel 5
