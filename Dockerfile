FROM golang:alpine as builder

RUN mkdir -p $GOPATH/gin-test
WORKDIR $GOPATH/gin-test
#COPY go.mod .
#COPY go.sum .

#RUN go mod download

#COPY . .

ADD . $GOPATH/gin-test
RUN go get -u github.com/gin-gonic/gin && \
  go get github.com/jinzhu/gorm && \
  go get github.com/jinzhu/gorm/dialects/postgres

#ENV PORT=${PORT}
