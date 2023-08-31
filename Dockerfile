FROM ubuntu:22.04

LABEL author="devhoodit"

RUN apt-get update && \
    apt-get -y install sudo && \
    sudo apt-get -y install systemctl && \
    sudo apt-get -y install wget

# install services
RUN sudo apt-get -y install mysql-server && \
    sudo apt-get -y install redis-server && \
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz && \
    sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

ENV PATH "/usr/local/go/bin:${PATH}"
ENV DB_PWD=""

RUN mkdir server

# copy files
COPY ./ ./server

# install packages
RUN cd server && go mod download

# CMD
CMD ["/bin/bash", "-c", "sudo systemctl start mysql", "&&", "mysql -u root -e \"alter user 'root'@'localhost' identified with mysql_native_password by '${DB_PWD}';\"" ]

EXPOSE 5500