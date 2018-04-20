FROM scratch

MAINTAINER John Weldon <johnweldon4@gmail.com>

COPY public /public/
ADD logsrv logsrv

ENV PORT 11181
ENV PUBLIC_DIR /public
ENV VERBOSE=true
ENV IGNORED_HOSTS=

EXPOSE 11181

ENTRYPOINT ["/logsrv"]
