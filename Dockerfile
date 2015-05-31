FROM ubuntu:14.04

EXPOSE 3000

ADD jcio-frontend /jcio-frontend
ADD public /public
ADD templates /templates

CMD ["/jcio-frontend"]
