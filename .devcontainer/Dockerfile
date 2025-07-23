FROM golang:1.24.5-alpine

ARG USERNAME=vscode
ARG UID=1000

# addgroup ile grup, adduser ile kullanıcı oluşturun
RUN addgroup -g ${UID} ${USERNAME} \
 && adduser  -D \
             -u ${UID} \
             -G ${USERNAME} \
             -h /home/${USERNAME} \
             -s /bin/sh \
             ${USERNAME}

# Temel paketler (git, ca-certificates, vs.)
RUN apk add --no-cache \
      git \
      ca-certificates

# Çalışma dizinini ayarlayın
WORKDIR /workspace

# vscode kullanıcısına geçiş
USER ${USERNAME}

# Modül önbelleğini kopyalayarak zaman kazanın (isteğe bağlı)
# ADD go.mod go.sum ./
# RUN go mod download

# Container başlatıldığında default komut
CMD ["sleep", "infinity"]