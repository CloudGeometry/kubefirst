FROM kubefirst/chubbo:0.1

ADD scripts/nebulous /scripts/nebulous
ADD gitops /gitops
ADD metaphor /metaphor

RUN apt-get update
RUN apt-get install dnsutils -y

CMD [ "/bin/bash" ]
