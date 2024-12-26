# Proyecto de Mensajes

Este proyecto es una aplicación escrita en Go que utiliza Kafka y ActiveMQ para la mensajería.

## Requisitos Técnicos

- Go 1.23.2 o superior
- Docker
- Docker Compose

## Instalación

Para poder usar la herramienta en local, sigua los siguientes pasos:

1. Realiza una copia del env.example y renómbralo a .env que se encuentra en la carpeta devEnv

```bash
cp devEnv/env.example devEnv/.env
```

2. Remplaza los valores de los envs del archivo .env con los de entorno local.

Para el Kafka
- KAFKA_HOST

Para el ActiveMQ
- ACTIVEMQ_HOST
- ACTIVEMQ_USER
- ACTIVEMQ_PASSWORD

Puede mirar que otros valores se pueden cambiar en el archivo .env en caso de ser necesario. 

En caso de no tener un ActiveMQ o Kafka local, puede utilizar el docker-compose que se encuentra en la carpeta devEnv

```bash
docker-compose -f devEnv/docker-compose.yml up activemq kafka -d
```

3. Instalar las dependencias del proyecto

```bash
go mod download
```

4. Ejecutar el proyecto
```bash
go run main.go
```

## Correr el proyecto desde docker

Para correr el proyecto desde docker, siga los siguientes pasos

1. Realiza una copia del env.example y renómbralo a .env que se encuentra en la carpeta devEnv, no modifique los valores en este caso.
```bash
cp devEnv/env.example devEnv/.env
```
2. Inicia el proyecto con docker-compose

```bash
docker-compose -f devEnv/docker-compose.yml up --build -d
```