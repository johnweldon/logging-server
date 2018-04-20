FROM scratch

LABEL maintainer="John Weldon <johnweldon4@gmail.com>" \
      company="John Weldon Consulting"

COPY public /public/
ADD logsrv logsrv

ENV PORT 11181
ENV PUBLIC_DIR /public
ENV VERBOSE=true
ENV IGNORE_HOSTS=

EXPOSE 11181

ENTRYPOINT ["/logsrv"]
