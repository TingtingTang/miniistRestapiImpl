# Start from a Debian image with the latest version of Go installed 
# and a workspace (GOPATH) configured at /go
FROM golang

# Gopy the local package files to the container's workspace
ADD  .  /user

# Install third-party go packages inside container
RUN mkdir -p /go/src
ADD install/gopkg.in  /go/src/gopkg.in 

RUN go get github.com/gorilla/mux 
RUN go get github.com/stretchr/testify
RUN go get github.com/asaskevich/govalidator
RUN go get github.com/dgrijalva/jwt-go


EXPOSE 80

