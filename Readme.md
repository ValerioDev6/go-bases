# CRUD de usuarios con Go y SQLite

Aplicación web para crear, listar, editar y eliminar usuarios. Las rutas HTTP usan únicamente `net/http`; las páginas y estilos están en `web/templates` y `web/static`.

## Requisitos

- Go compatible con la versión indicada en `go.mod`.
- SQLite 3 y su biblioteca de desarrollo (`sqlite3.h` y `libsqlite3`), porque el acceso a la base de datos usa cgo.

## Ejecutar

Desde la raíz del repositorio:

```bash
go run -buildvcs=false ./cmd/server
```

Abre `http://localhost:8080/users`. La base de datos se crea automáticamente en `data/users.db` y la tabla `users` se inicializa al arrancar.

Variables de entorno opcionales:

- `PORT`: puerto HTTP (por defecto `8080`).
- `ENV`: `dev` o `prod` para el formato del log (por defecto `dev`).
- `DB_PATH`: ruta del archivo SQLite (por defecto `data/users.db`).

La ruta `GET /` redirige al listado. `GET /health` responde `ok`.

## Verificar

```bash
go test ./...
go vet ./...
go build -buildvcs=false ./...
```
go get github.com/mattn/go-sqlite3


estructura recoemdano 

/internal
/book
    service.go
    store.go
    model.go
