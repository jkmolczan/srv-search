# srv-search

This service provides an HTTP API to search for the index of a given value from a preloaded dataset.

## Features

- Efficient binary search algorithm.
- Approximate matching within 10% if the exact value is not found.
- Configurable service port and log level via `config.yaml`.
- Logging with support for `Info`, `Debug`, and `Error` levels.
- Includes unit tests.
- Automation with a Makefile.

## Setup and Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/jkmolczan/srv-search.git

2. **Prepare the data:**

    Place your `input.txt` file in the main directory.

3. **Configure the service:**

    Edit `config.yaml` to set the desired port and log level.

4. **Run the application::**
    ```bash
    make run
   ```
   or using docker:
   ```bash
   ./script/run
   ```
   
5. **Run the tests:**

    ```bash
    make test
   ```
   or using docker:
   ```bash
   ./script/test
   ```
## Usage

 * Endpoint: `GET /search/numbers/index/{number}` 
 * Example: `curl http://localhost:8080/search/numbers/index/100`
 * Response:
    ```json
    {
      "index": 100,
      "number": 100,
      "message": "Value found."
    }
    ```
## Logging

* Log Levels: Set in `config.yaml` (`debug`, `info`, `error`).
* Log Output: Printed to the console with timestamp, log level, file, line number, and message:
* Example:
    ```
    {"time":"2024-10-19T21:47:35.56407+02:00","level":"DEBUG","prefix":"echo","file":"search_handler.go","line":"57","message":"SearchNumberIndex: invalid number path param: one"}
    ```

## API documentation
API documentation is available under endpoint `/docs/api/swagger.yaml` after running the application or in the `swagger.yaml` file.

## Running the application

#### Makefile Commands

* `make run`: Build and run the application.
* `make test`: Run the unit tests.

#### Docker Commands
* `./script/run`: Build and run the application with Docker.
* `./script/test`: Run the unit tests with Docker.
* `./script/stop`: Stop the application container.