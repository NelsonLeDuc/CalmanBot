FROM golang:1.9-alpine 
EXPOSE 4000
EXPOSE 4004

# setup GOPATH and friends
#
# TECHNICALLY, you don't have to do these three cmds as the 
# golang:alpine image actually uses this same directory structure and
# already has $GOPATH set to this same structure.  You could just
# remove these two lines and everything below should continue to work.
#
# But, I like to do it anyways to ensure my proper build
# path in case I experiment with different Docker build images or in
# case the #latest image changes structure (you should really use
# a tag to lock down what version of Go you are using - note that I
# locked you to the docker image golang:1.7-alpine above, since that is
# the current latest you were using, with bug fixes).
#
# RUN  mkdir -p /go/src \
#   && mkdir -p /go/bin \
#   && mkdir -p /go/pkg
ENV GOPATH=/go
# ENV PATH=$GOPATH/bin:$PATH   

# now copy your app to the proper build path
# RUN mkdir -p $GOPATH/src/app 
ADD . $GOPATH/src/github.com/nelsonleduc/calmanbot
# COPY . $GOPATH/src

# should be able to build now
WORKDIR $GOPATH/src/github.com/nelsonleduc/calmanbot/
RUN go build .
CMD ["/go/src/github.com/nelsonleduc/calmanbot/calmanbot", "-v"]