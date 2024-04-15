# Car catalog
Сервис для хранения информации о машинах. По переданным регистрационным номерам запрашивает в сторонней базе полную информацию. 

Предполагаемое стороннее хранилище должно соответствовать сваггер-описанию:

    ```yaml
    openapi: 3.0.3
    info:
    title: Car info
    version: 0.0.1
    paths:
    /info:
    get:
    parameters:
    - name: regNum
    in: query
    required: true
    schema:
    type: string
    responses:
    '200':
    description: Ok
    content:
    application/json:
    schema:
    $ref: '#/components/schemas/Car'
    '400':
    description: Bad request
    '500':
    description: Internal server error
    components:
    schemas:
    Car:
    required:
    - regNum
    - mark
    - model
    - owner
    type: object
    properties:
    regNum:
    type: string
    example: X123XX150
    mark:
    type: string
    example: Lada
    model:
    type: string
    example: Vesta
    year:
    type: integer
    example: 2002
    owner:
    $ref: '#/components/schemas/People'
    People:
    required:
    - name
    - surname
    type: object
    properties:
    name:
    type: string
    surname:
    type: string
    patronymic:
    type: string
    ```


В качестве локального хранилища используется PostgreSQL.

## Объекты

Машина:

    RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  Person `json:"owner"`

Владелец:

    Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`

Страница-результат запроса get:

    Cars        []Car `json:"cars"`
	PageNo      int   `json:"page_number"`
	Limit       int   `json:"limit"`
	PagesAmount int   `json:"pages_amount"`

## API
Сервис работает с форматом JSON.

Доступные методы:

    GET /cars - возвращает страницу машин, отфильтрованных согласно параметрам запроса
	POST /cars - запрашивает у стороннего хранилища информацию по машинам с переданными регистрационными номерами и, при её наличии, добавляет машины в локальную базу данных
	PATCH /cars/:regNum - редактирует информацию о машине согласно запросу
	DELETE /cars/:regNum - удаляет машину с переданным регистрационным номером

## Переменные окружения

Сервис умеет считывать переменные из файла .env в директории исполняемого файла (в корне проекта).

В примерах указаны дефолтные значения. Если программа не сможет считать пользовательские env, то возьмет их.
Однако без заполненных дополнительно переменных Postgres и перемнной стороннего хранилища сервис корректно работать не будет.

Переменные сервера:

    SERVER_LISTEN=:8088
    SERVER_READ_TIMEOUT=5s
    SERVER_WRITE_TIMEOUT=5s
    SERVER_IDLE_TIMEOUT=30s

Переменные Postgres:

    PG_USER=
	PG_PASSWORD=
	PG_HOST=localhost
	PG_PORT=5432
	PG_DATABASE=

Переменные предполагаемого стороннего хранилища:

    REC_URL=

Переменные логгера:

    LOG_LEVEL=info