# API de usuarios, roles y autenticación con Go y PostgreSQL

API para administrar roles y usuarios y autenticar sesiones con PostgreSQL mediante `database/sql` y el driver `pgx`, sin ORM. El módulo `book` se conserva en el repositorio, pero no se conecta ni registra en la aplicación.

```text
internal/
├── app/                 # composición, rutas y servidor HTTP
├── auth/                # registro, login, logout y me
├── book/                # módulo conservado, no registrado
├── role/                # CRUD de roles; incluye dtos/request.go
├── user/                # CRUD de usuarios; incluye dtos/request.go
├── config/              # variables de entorno y valores predeterminados
├── transport/http/      # formato compartido de respuestas JSON
├── middleware/          # autenticación Bearer y contexto del usuario
└── platform/
    ├── database/        # conexión a la base de datos
    ├── logger/          # logging estructurado
    └── openapi/         # especificación OpenAPI y Swagger UI
migrations/              # migraciones SQL versionadas
storage/                 # archivos locales de ejecución
```

Cada feature mantiene su propio `model.go`, `service.go`, `store.go` y `handler.go`. El almacenamiento puede cambiar sin modificar las reglas de negocio de la feature.

## Requisitos

- Go compatible con la versión indicada en `go.mod`.
- PostgreSQL en ejecución.
- CLI `migrate` con soporte para PostgreSQL.

## Ejecutar

Desde la raíz del repositorio:

```bash
go run .
```

La API queda disponible en `http://localhost:8080` después de crear la base y aplicar las migraciones.

Swagger UI queda disponible en `http://localhost:8080/docs/`. La especificación OpenAPI se sirve en `http://localhost:8080/docs/openapi.yaml`.

Al iniciar, `godotenv.Load()` carga el archivo local `.env`; después `internal/config` lee las variables mediante `os.Getenv`. Si `.env` no existe, la aplicación usa las variables que ya tenga el proceso, como ocurre normalmente en producción. `DATABASE_URL` es obligatoria.

- `PORT`: puerto HTTP (por defecto `8080`).
- `DATABASE_URL`: conexión PostgreSQL.
- `JWT_SECRET`: clave de al menos 32 caracteres para firmar y validar JWT.
- `ENV`: `dev` usa logs legibles; `prod` usa logs JSON.

Todas las respuestas usan este formato:

```json
{
  "data": {},
  "success": true,
  "status": 200,
  "errors": []
}
```

Cuando hay un error, `data` es `null`, `success` es `false` y `errors` contiene objetos con `message` y, para validaciones, `field`.

Rutas disponibles:

- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/me` — requiere `Authorization: Bearer <token>`
- `GET/POST /roles`
- `GET/PUT/DELETE /roles/{id}`
- `GET/POST /users`
- `GET/PUT/DELETE /users/{id}`

`register` asigna el rol predeterminado `user`, cifra la contraseña con bcrypt y devuelve el usuario creado. `login` devuelve un JWT Bearer con 24 horas de duración. Su payload contiene `sub`, `name`, `email`, `iat`, `exp` e `iss`. Los JWT no se almacenan.

## PostgreSQL y migraciones

Crear la base de datos y aplicar las migraciones:

```bash
make db-create
make migrate-up
```

Consultar la versión aplicada:

```bash
make migrate-version
```

Revertir una migración:

```bash
make migrate-down
```

Crear una migración nueva:

```bash
make migrate-create name=add_example
```

## Verificar

```bash
go vet ./...
go build ./...
```


[//]: # ()
[//]: # ({)

[//]: # ("email": "vquispealarcon@gmail.com",)

[//]: # ("password": "acasio2003")

[//]: # (})
