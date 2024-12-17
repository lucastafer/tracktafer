### Events

### 1st Event Received

RouteCreated

- id
- distance
- directions
- - lat
- - lng

### Side effect for 1st event (calculate freigth and return another event)

FreigthCalculated

- route_id
- amount

---

### 2nd Event Received

DeliveryStarted

- route_id

### Side effect for 2nd event (return another event with driver coordinates)

DriverMoved

- route_id
- lat
- lng
