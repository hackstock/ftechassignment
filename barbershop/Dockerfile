FROM golang:1.19
ENV NAME=barbershop 
ENV APP_DIR=/${NAME}
ENV GOOS=linux
ENV GO_LINKER_FLAGS=-ldflags="-s -w"
COPY . ${APP_DIR}
WORKDIR ${APP_DIR}
RUN go mod download && go mod verify
RUN go build -o ${NAME}  -ldflags="-s -w" ${APP_DIR}/cmd/barbershop/main.go 

ENV SEAT_CAPACITY=10
ENV ARRIVAL_RATE_IN_SEC=100
ENV CUT_DURATION_IN_MSEC=1000
ENV OPENED_UNTIL_IN_SEC=10
ENV ENVIRONMENT="development"
EXPOSE 9000
ENTRYPOINT [ "./barbershop" ]