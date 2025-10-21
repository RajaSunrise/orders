# Order Service simple with kafka

## create topic with command
```bash
# orders-topic
docker exec -it kafka kafka-topics.sh --create --topic processed-orders --bootstrap-server kafka:9092

# processed-order topic
docker exec -it kafka kafka-topics.sh --create --topic processed-orders --bootstrap-server kafka:9092
```
### Test producer
```bash
# orders-topic producer
docker exec -it kafka opt/kafka/bin/kafka-console-producer.sh --topic orders-topic --bootstrap-server kafka:9092

# processed-order
docker exec -it kafka opt/kafka/bin/kafka-console-producer.sh --topic orders-topic --bootstrap-server kafka:9092
```
### Testing
Succes response
```json
{"id":"order001","user_id":"user123","product":"laptop","quantity":1,"total_amount":15000.00}

{"id":"order004","user_id":"user456","product":"smartphone","quantity":2,"total_amount":15090.00}

{"id":"order006","user_id":"user456","product":"laptop","quantity":3,"total_amount":150980.00}

{"id":"order009","user_id":"user123","product":"smartphone","quantity":4,"total_amount":15200.00}
```
Failed response
```json
{"id":"order0012","user_id":"user456","product":"pillow","quantity":3,"total_amount":1200.00}

{"id":"order0011","user_id":"user921","product":"laptop","quantity":3,"total_amount":1200.00}

{"id":"order0010","user_id":"user900","product":"laptop","quantity":3,"total_amount":1200.00}

{"id":"order0014","user_id":"user456","product":"ice cream","quantity":3,"total_amount":1200.00}
```