FROM ubuntu:22.04

ENV USER root

ENV HOME /root

WORKDIR $HOME

ARG VERSION

ARG VERSION_URL

ARG SAVE_DIR=$HOME/saves

ARG LOG_FILE=$HOME/server.log

ARG SERVER_NAME

ARG WORLD_NAME

RUN apt-get update && \
    apt-get install -y \
    wget \
    software-properties-common \
    xvfb

RUN dpkg --add-architecture i386

## Add steamcmd
RUN echo steam steam/question select "I AGREE" | debconf-set-selections  && echo steam steam/license note '' | debconf-set-selections

RUN apt-get update && apt-get install -y --no-install-recommends \
    steamcmd \
    locales

# Add unicode support
RUN locale-gen en_US.UTF-8
ENV LANG 'en_US.UTF-8'
ENV LANGUAGE 'en_US:en'

RUN ln -s /usr/games/steamcmd /usr/bin/steamcmd

## Add wine
RUN wget -nc https://dl.winehq.org/wine-builds/winehq.key && \
    mv winehq.key /usr/share/keyrings/winehq-archive.key

RUN wget -nc https://dl.winehq.org/wine-builds/ubuntu/dists/jammy/winehq-jammy.sources && \
    mv winehq-jammy.sources /etc/apt/sources.list.d/

RUN apt-get update &&  \
    apt-get install -y --install-recommends winehq-staging

RUN wine --version

ENV WINEARCH=win64

RUN mkdir "vrising"

RUN steamcmd +force_install_dir $HOME/vrising +login anonymous +app_update 1829350 +quit

COPY start.sh start.sh

RUN chmod +x start.sh

LABEL org.opencontainers.image.description="V Rising Linux version ${VERSION}. See changelog here: ${VERSION_URL}."
LABEL org.opencontainers.image.url='ghcr.io/hostfactor/vrising-server'
LABEL org.opencontainers.image.version=${VERSION}
LABEL org.opencontainers.image.authors='eddie@hostfactor.io'

EXPOSE 9876/udp

EXPOSE 9877/udp

CMD ["./start.sh"]
