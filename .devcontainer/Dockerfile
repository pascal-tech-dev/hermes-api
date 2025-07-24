# Resmi Go imajını kullanıyoruz (versiyonu ihtiyaçlarınıza göre değiştirin)
FROM golang:1.24.5-bookworm

ARG USERNAME=vscode
ARG UID=1000

# Grup ve kullanıcı oluşturma (Debian tabanında addgroup/adduser zaten var)
RUN groupadd --gid ${UID} ${USERNAME} \
 && useradd  --uid ${UID} \
             --gid ${USERNAME} \
             --create-home \
             --shell /bin/bash \
             ${USERNAME}

# Temel paketler (Debian/Ubuntu tabanında apt-get ile)
RUN apt-get update && apt-get install -y --no-install-recommends \
      bash \
      curl \
      wget \
      unzip \
      zip \
      tar \
      gzip \
      bzip2 \
      git \
      ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Çalışma dizinini ayarla
WORKDIR /workspace
RUN chown ${USERNAME}:${USERNAME} /workspace

# vscode kullanıcısına geçiş
USER ${USERNAME}

# Modül önbelleğini kopyalayarak zaman kazanın (isteğe bağlı)
# COPY go.mod go.sum ./
# RUN go mod download

# Container başlatıldığında default komut
CMD ["sleep", "infinity"]