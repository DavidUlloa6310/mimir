FROM rapidsai/rapidsai:25.04-cuda12.0-runtime-ubuntu22.04-py3.10

RUN apt-get update && apt-get install -y \
    wget \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp
RUN wget https://go.dev/dl/go1.20.10.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.20.10.linux-amd64.tar.gz \
    && rm go1.20.10.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend/database/clustering_cuml.py ./backend/database/

COPY . .

RUN go build -o /goapp

EXPOSE 8080

CMD ["/goapp"]