FROM  registry.manadrdev.com/ubuntu

LABEL maintainer="Kyo Nguyen <kyo.nguyen@manadr.com>"
ENV   GO_FILE=go1.9.2.linux-amd64.tar.gz
ENV   GOROOT=/usr/local/go GOPATH=/opt/golang
ENV   PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV   APP_PATH=$GOPATH/src/github.com/mobile-health/scheduler-service 

RUN   apt-get update && apt-get install -y curl git wget make ruby mysql-client
RUN   gem install tiller
RUN   wget https://redirector.gvt1.com/edgedl/go/$GO_FILE && tar -C /usr/local -xzf $GO_FILE && rm $GO_FILE
RUN   mkdir $GOPATH/bin -p && mkdir $GOPATH/src -p
RUN   curl https://glide.sh/get | sh

WORKDIR  $APP_PATH
COPY  glide.yaml ./
RUN   glide up
COPY  . ./
RUN   rm $APP_PATH/conf/config.yaml
RUN   make build-linux
COPY  tiller/ /etc/tiller
COPY  conf/config.yaml.example /etc/tiller/templates/config.erb
ADD   entrypoint-ci.sh /tmp/entrypoint.sh
RUN   chmod a+x /tmp/entrypoint.sh

RUN   mkdir /var/log/scheduler
CMD   ["/usr/local/bin/tiller" , "-v"]
HEALTHCHECK --interval=10s --timeout=5s \
      CMD curl -f http://localhost:8080/api/v1/ || exit 1

EXPOSE  8080