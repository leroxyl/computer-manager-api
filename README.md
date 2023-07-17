# Greenbone Code Challenge

A basic application for system administrators to keep track of the computers issued by their company.
The application provides a REST API to manage computer-related datasets that are stored in a postgres database.

### Computer Model

- MAC address (required)
- computer name (required)
- IP address (required)
- employee abbreviation[^1] (optional)
- description (optional)

[^1]: The employee abbreviation consists of 3 letters. For example Max Mustermann should be mmu.

### Prerequisites

- you need a running postgres database instance and set the DSN via env variable `GREENBONE_POSTGRES_DSN`
  - the database table will be automatically created at application start-up
- start the notification service by running the following command
    ```shell
    docker pull greenbone/exercise-admin-notification && \
    docker run -p 8081:8080 greenbone/exercise-admin-notification
    ```

### Endpoints

| Action                       | HTTP Method | Path                                         | Content-Type       | Request Body                                           | Description                                            |
|------------------------------|-------------|----------------------------------------------|--------------------|--------------------------------------------------------|--------------------------------------------------------|
| Create(Computers)            | **POST**    | /computers                                   | `application/json` | see [JSON](#JSON)                                      | Store the data of a new computer                       |
| Read(Computers)              | **GET**     | /computers/{mac}                             | _none_             | _empty_                                                | Get the data of a computer                             |
| Update(Computers)            | **PUT**     | /computers/{mac}                             | `application/json` | like [JSON](#JSON), but field `macAddr` can be omitted | Update the data of a computer                          |
| Delete(Computers)            | **DELETE**  | /computers/{mac}                             | _none_             | _empty_                                                | Delete the data of a computer                          |
| Read All(Computers)          | **GET**     | /computers                                   | _none_             | _empty_                                                | Get the data of all computers                          |
| Read All Computers(Employee) | **GET**     | /employees/{employee-abbreviation}/computers | _none_             | _empty_                                                | Get the data of all assigned computers for an employee |

> Update(Computers):
>   - operation does not support update of the MAC address
>   - operation will fail if any required field (other than MAC address) is not provided in request body
>   - operation will overwrite optional fields with empty value if field is not provided in request body

> Delete(Computers):
>   - repeatedly calling delete on same resource will return `200`

### JSON

```json
{
  "macAddr": "<MAC address (required)>",
  "computerName": "<computer name (required)>",
  "ipAddr": "<IP address (required)>",
  "employeeAbbr": "<employee abbreviation (optional)>",
  "description": "<description (optional)>"
}
```

### Example Requests

#### Create

    curl localhost:8080/computers \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"macAddr": "00:1B:44:11:3A:B7", "computerName": "localhorst", "ipAddr": "127.0.0.1", "employeeAbbr": "rpm", "description": "hello :)"}' \
    -i

#### Read

    curl localhost:8080/computers/00:1B:44:11:3A:B7 \
    -X GET \
    -i

#### Update

    curl localhost:8080/computers/00:1B:44:11:3A:B7 \
    -X PUT \
    -H "Content-Type: application/json" \
    -d '{"computerName": "localhorst", "ipAddr": "127.0.0.2", "employeeAbbr": "mmu"}' \
    -i

#### Delete

    curl localhost:8080/computers/00:1B:44:11:3A:B7 \
    -X DELETE \
    -i

#### Read All

    curl localhost:8080/computers \
    -X GET \
    -i

#### Read All For Employee

    curl localhost:8080/employees/rpm/computers \
    -X GET \
    -i

### Limitations

- MAC address is considered to be unique and used as primary key in database table
- employee abbreviation is considered to be unique for each employee
- client authentication is considered to be out of scope

