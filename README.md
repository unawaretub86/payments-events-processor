# Lambda en Golang para Procesar y enviar Eventos de SQS, hacer interaccion con DynamoDB

![GitHub](https://github.com/unawaretub86/payments-events-processor)
![GitHub contributors](https://github.com/unawaretub86)

Esta es una función Lambda escrita en Golang que se encarga de procesar eventos proporcionados desde un SQS que desencadena una lambda. Puede utilizarse para manejar mensajes en SQS y ejecutar lógica personalizada en función de las solicitudes entrantes, como guardar en dynamoDB y actualizar items y enviar mensajes via SQS.

## Requisitos

- Go 1.13 o superior
- AWS CLI y AWS SAM CLI  configurada con las credenciales adecuadas
- API Gateway configurada para enrutar eventos a esta función Lambda

## Estructura del Proyecto

- `cmd/api/main.go`: El archivo principal de la función Lambda que contiene la lógica de procesamiento.
- `template.yaml`: Un archivo de plantilla SAM que define los recursos necesarios para desplegar la función Lambda y la API Gateway.

## Despliegue

Siga estos pasos para desplegar la función Lambda, SQS y tabla en base de datos utilizando el archivo `template.yaml`.


```bash

1. Clona este repositorio:

git clone https://github.com/unawaretub86/payments-events-processor

2. Asegúrese de tener la AWS CLI configurada correctamente con las credenciales adecuadas: 

- aws configure

3. Despliegue la función Lambda y la API Gateway utilizando CloudFormation:

- sam deploy --guided

4. Una vez se haya creado una orden en https://github.com/unawaretub86/order-processor-events esta enviara un mensaje via SQS , recibiremos el mensaje y lo procesaremos guardando un objeto payment en DynamoDB con el orderId y el status "PENDING".

4. Una vez completado el pago, se debera haber enviado una peticion http PATCH con el orderId a su función Lambda creada en el servicio https://github.com/unawaretub86/payments-processor.

5. Una vez enviada la peticion esta lambda procesara el mensaje enviado, actualizara un pago en dynamoDB con el status "PAID" , esta enviara un mensaje via SQS a https://github.com/unawaretub86/order-processor-events. 
