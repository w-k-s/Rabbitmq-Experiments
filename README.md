# RabbitMQ Experiments

Hands-on experiments to get a better understanding of RabbitMQ.

## Getting Started

To run this project, you'll need to have an instance of rabbitmq running. If you have Docker installed, the easiest way to launch an instance of rabbitmq is to run the command:

```sh
docker run -p 5672:5672 -p 15672:15672 -d --name rabbitmq rabbitmq:3-management
```

You should then be able to access the admin portal on your browser at the address: `http://localhost:15672` using the default username: `guest`, and password: `guest`.

---