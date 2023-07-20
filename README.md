# Computer Manager API

A basic application for system administrators to keep track of the computers issued by their company.
The application provides a REST API to manage computer-related datasets that are stored in a postgres database.

### Build

#### Build binary

```shell
go build
```

Creates the binary `computer-manager-api`.

#### Build Docker image

```shell
docker build . -t computer-manager-api
```

### Run

```shell
docker-compose up
```

The application listens on port `8080`.

```text
http://localhost:8080/
```

The database can be reached via the DSN `postgres://postgres:postgres@localhost:5433/computer_manager_api`.

```shell
psql postgres://postgres:postgres@localhost:5433/computer_manager_api
```

### Endpoints

| Action                       | HTTP Method | Path                                         | Content-Type       | Request Body                                                          | Description                                            |
|------------------------------|-------------|----------------------------------------------|--------------------|-----------------------------------------------------------------------|--------------------------------------------------------|
| Create(Computers)            | **POST**    | /computers                                   | `application/json` | see [JSON](#computer-model-json)                                      | Store the data of a new computer                       |
| Read(Computers)              | **GET**     | /computers/{mac}                             | _none_             | _empty_                                                               | Get the data of a computer                             |
| Update(Computers)            | **PUT**     | /computers/{mac}                             | `application/json` | like [JSON](#computer-model-json), but field `macAddr` can be omitted | Update the data of a computer                          |
| Delete(Computers)            | **DELETE**  | /computers/{mac}                             | _none_             | _empty_                                                               | Delete the data of a computer                          |
| Read All(Computers)          | **GET**     | /computers                                   | _none_             | _empty_                                                               | Get the data of all computers                          |
| Read All Computers(Employee) | **GET**     | /employees/{employee-abbreviation}/computers | _none_             | _empty_                                                               | Get the data of all assigned computers for an employee |

> Update(Computers):
>   - operation does not support update of the MAC address
>   - operation will fail if any required field (other than MAC address) is not provided in request body
>   - operation will overwrite optional fields with empty value if field is not provided in request body

> Delete(Computers):
>   - repeatedly calling delete on same resource will return `200`

### Computer Model JSON

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

- MAC address is considered to be unique and used as primary key in the database table
- employee abbreviation is considered to be unique for each employee

### Next Steps

- add more tests
- validate and normalize incoming data, i.e. parse types, check constraints, store in standardized representation etc.
  - sanitizing input to prevent SQL injection is not needed since GORM takes care of this
- return meaningful HTTP status codes and responses, especially for client-sided errors
- add client authentication, maybe user management
- enable TLS
- make service more robust, for example with request timeouts and throttling
- add exhaustive logging
- add metrics (Prometheus)
- add readiness and health check endpoints
- support graceful shutdown
- consider using advanced database migration framework like `github.com/golang-migrate/migrate`
- automate workflows
  - add makefile
  - add CI (GitHub Actions)
