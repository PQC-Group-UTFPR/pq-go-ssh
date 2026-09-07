FROM ubuntu:latest

WORKDIR /pq-ssh
RUN apt-get update --fix-missing
RUN apt-get install -y bash sudo

#instalar deps 


#Copia o código para a imagem
COPY . .

#executa o install and run
RUN chmod +x ./install_and_run.sh

CMD ["sh", "-c", "./install_and_run.sh $NAME $IPSERVER"   ]




