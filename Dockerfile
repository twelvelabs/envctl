FROM alpine:3.23

COPY envctl /usr/local/bin/envctl

CMD ["/usr/local/bin/envctl"]
