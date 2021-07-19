FROM golang:1.16

RUN go install golang.org/dl/gotip@latest && gotip download dev.fuzz

ENTRYPOINT [ "/bin/bash" ]