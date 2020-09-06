FROM ubuntu:18.04

MAINTAINER JamesClonk

EXPOSE 3000

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY jcio-frontend /jcio-frontend
COPY public /public
COPY templates /templates

ENV JCIO_ENV production
ENV PORT 3000
ENV JCIO_CMS_DATA https://github.com/jamesclonk-io/content/archive/master.zip

CMD ["/jcio-frontend"]
