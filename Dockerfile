FROM golang:latest

# Install commonlib when internet connected
RUN git clone https://34c1a18a731ddfe9b3f3b3fb1dccc3ee43afc465@github.com/kurtd5105/SENG-468-Common-Lib.git ${GOPATH}/src/github.com/kurtd5105/SENG-468-Common-Lib

# Build directories and binary
RUN mkdir /http_server
WORKDIR /http_server
ADD . /http_server/
RUN go build -o main

# Ensure that all the arguments are all provided by erroring out if not.
ARG PORT
ENV PORT ${PORT}
RUN test -n "$PORT"
ARG DB_ADDRESS
ENV DB_ADDRESS ${DB_ADDRESS}
RUN test -n "$DB_ADDRESS"
ARG TX_ADDRESS
ENV TX_ADDRESS ${TX_ADDRESS}
RUN test -n "$TX_ADDRESS"
ARG LOG_ADDRESS
ENV LOG_ADDRESS ${LOG_ADDRESS}
RUN test -n "$LOG_ADDRESS"

# Make port 8084 available
EXPOSE 8084

# Run it
CMD /http_server/main -port $PORT -tx $TX_ADDRESS -db $DB_ADDRESS -log $LOG_ADDRESS
