FROM alpine:3.22

COPY envctl /usr/local/bin/envctl

CMD ["/usr/local/bin/envctl"]
