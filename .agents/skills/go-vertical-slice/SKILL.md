---
name: go-vertical-slice
description: "Diseñar, desarrollar, revisar o refactorizar APIs Go con Vertical Slice Architecture: organización por feature, DTOs de entrada, respuestas HTTP uniformes, composición explícita, PostgreSQL sin ORM y biblioteca estándar sin frameworks. Los tests solo se crean o ejecutan cuando el usuario los solicita explícitamente."
---

# Vertical Slice Architecture profesional para Go

Organizar el sistema por funcionalidades completas. Cada slice, como `user`, `role` o `book`, contiene su modelo, casos de uso, contrato de almacenamiento, implementación y transporte HTTP. Mantener juntas las piezas que cambian por la misma razón.

Esta arquitectura evita separar todo el proyecto en carpetas globales por tipo técnico. Una feature debe poder entenderse revisando principalmente su propio directorio.

## Principios obligatorios

1. Organizar por feature y usar nombres singulares en minúsculas.
2. Mantener el flujo `Handler → Service → Store`.
3. Inyectar dependencias mediante constructores; no usar estado global ni `init()` para conectar la aplicación.
4. Mantener `main.go` pequeño y concentrar la composición en `internal/app`.
5. Usar la biblioteca estándar para HTTP, JSON, contexto, logging y configuración.
6. No añadir frameworks HTTP, routers, ORM, validadores ni contenedores de dependencias.
7. Permitir únicamente dependencias externas justificadas por una capacidad que Go no incluye, como el driver de PostgreSQL o un algoritmo seguro de hash de contraseñas.
8. Mantener reglas de negocio fuera de handlers, SQL, configuración y composición.
9. Crear archivos y carpetas solo cuando tengan una responsabilidad real.
10. No crear, modificar ni ejecutar tests a menos que el usuario lo solicite explícitamente en la tarea actual.

## Alcance

Leer los archivos reales, `go.mod` y las instrucciones del repositorio. El entrypoint puede ser `main.go` o `cmd/server/main.go`: conservar su ubicación salvo que la tarea requiera moverlo. Respetar la versión de Go y los contratos públicos existentes.

Crear o editar esta skill solo modifica la skill. Añadir una feature no autoriza migrar todo el repositorio. Si instrucciones antiguas describen capas globales, una solicitud explícita del usuario de adoptar esta arquitectura prevalece para el alcance solicitado.

## Estructura de referencia

```text
main.go                     # o cmd/server/main.go si el proyecto ya lo usa
internal/
  app/
    app.go                  # composición de dependencias
    routes.go               # registro central de slices
    server.go               # configuración del servidor HTTP
  config/
    config.go               # lectura y validación del entorno
  platform/
    database/
      postgres.go           # apertura y comprobación de PostgreSQL
    logger/
      logger.go             # log estructurado compartido
    openapi/
      handler.go            # Swagger UI y especificación embebida
      openapi.yaml
  transport/
    http/
      response/
        response.go         # sobre JSON compartido
  middleware/               # middleware realmente compartido
  auth/
    dtos/
      request.go           # requests de registro y login
    model.go               # sesión y errores de autenticación
    password.go            # hash y comparación de contraseñas con bcrypt
    service.go             # AuthService, registro, login y validación JWT
    store.go               # CreateUser, GetUserByEmail y GetUserByID
    handler.go             # AuthHandler: register, login y me
  middleware/
    auth.go                # validación Bearer y usuario autenticado
  role/
    dtos/
      request.go           # entrada HTTP de la feature
    model.go
    service.go
    store.go
    handler.go
  user/
    dtos/
      request.go
    model.go
    service.go
    store.go
    handler.go
  book/                     # misma convención; puede quedar desconectada
    model.go
    service.go
    store.go
    handler.go
migrations/
  000001_create_roles.up.sql
  000001_create_roles.down.sql
storage/                    # datos locales si una feature lo requiere
Makefile                    # ejecución, build y migraciones
```

Los archivos principales de una slice comparten paquete. `dtos` es un subpaquete porque representa el contrato del transporte, no la entidad del dominio. Evitar carpetas globales `models`, `services`, `stores`, `repositories`, `controllers` o `handlers`.

| Archivo | Responsabilidad |
| --- | --- |
| `model.go` | Entidades, entradas del servicio y errores de dominio. Sin dependencias HTTP o SQL. |
| `service.go` | Casos de uso, validación de negocio y coordinación. Define la interfaz `Store` que consume, con métodos concretos necesarios. |
| `store.go` | Implementación del almacenamiento: `MemoryStore`, `FileStore` u otra autorizada. Sin respuestas HTTP ni reglas de negocio. |
| `dtos/request.go` | Contratos de entrada HTTP. Mantiene los detalles JSON fuera del modelo de dominio. |
| `handler.go` | Rutas, métodos HTTP, parámetros, decodificación de DTOs y traducción de errores a estados HTTP. Delega al servicio. |

## Dirección de dependencias

```text
main → app → handlers → services → interfaces Store
                   ↘ stores concretos → database/sql → PostgreSQL
```

- `main` solo inicia contexto, carga el arranque y maneja el error final.
- `app` conecta configuración, logger, base de datos, stores, servicios, handlers y servidor.
- El handler conoce al servicio, pero el servicio no conoce HTTP.
- El servicio define la interfaz `Store` mínima que necesita.
- `store.go` implementa esa interfaz y conoce `database/sql`.
- Una feature no importa `internal/app`.
- Evitar que una feature importe directamente otra feature. Cuando necesite colaboración, definir una interfaz pequeña en la feature consumidora e inyectarla desde `app`.
- El código compartido solo pertenece a `platform`, `transport`, `middleware` o `config` cuando dos o más slices realmente lo necesitan.

Usar constructores explícitos: `NewService(store Store) *Service`, `NewHandler(service *Service, logger *slog.Logger) *Handler`, `NewMemoryStore() *MemoryStore`. Ajustar los parámetros a las dependencias reales. Dentro del paquete usar `Service` y `Handler`; desde fuera se leen `book.Service` y `book.Handler`.

El flujo es `Handler → Service → Store`. Inyectar dependencias mediante constructores, incluido el logger donde se use. Evitar estado global e inicialización oculta con `init()`. No crear interfaces para cada struct: definirlas en el consumidor cuando exista una frontera útil.

## Mantener main.go pequeño

`main` carga el entorno, llama al arranque de la aplicación y maneja el error final y código de salida. No contiene reglas de negocio, SQL, handlers ni registro individual de todas las rutas. No imponer un límite artificial de líneas.

`internal/app` construye configuración, logger, stores, servicios, handlers y servidor. Cada handler puede exponer `RegisterRoutes(mux *http.ServeMux)`; la composición llama al método de cada feature. Usar un mux explícito y `http.Server` con tiempos de espera adecuados.

La apertura de recursos y la inicialización de persistencia pertenecen al ciclo de vida de la aplicación. Para una API sencilla, iniciar el servidor con `Server.ListenAndServe()` y permitir que el sistema operativo administre la terminación del proceso. Añadir señales, goroutines y `Server.Shutdown` únicamente cuando el usuario o los requisitos del despliegue pidan apagado ordenado. Devolver errores hasta `main`; evitar `log.Fatal` en capas internas.

Dividir la composición en `routes.go` o `server.go` cuando crezca. Las reglas siguen dentro de cada feature: no trasladar un main gigante a un app gigante. Las features no importan `app`. Para colaboración entre features, definir una interfaz pequeña en quien la consume e inyectar la implementación desde la composición, evitando ciclos.

## Configuración y variables de entorno

Guardar la configuración local en `.env` con formato `KEY=VALUE`. En producción, hacer que el entorno de ejecución inyecte las variables antes de iniciar el proceso. `internal/config` debe limitarse a leerlas con `os.Getenv` o `os.LookupEnv`, aplicar valores predeterminados no sensibles y validar las obligatorias.

No implementar un parser de archivos `.env` dentro del proyecto. Si el usuario requiere que `go run .` cargue `.env` directamente, usar una dependencia pequeña y específica como `github.com/joho/godotenv`, llamar `godotenv.Load()` una sola vez al inicio de `main` y conservar la configuración en `internal/config`. La ausencia de `.env` no debe impedir el arranque cuando el proceso ya recibe variables del entorno. `Load` no sobrescribe variables existentes, por lo que producción conserva prioridad.

No dejar URLs de conexión ni contraseñas como valores predeterminados dentro del código Go. Validar al iniciar las variables obligatorias, como `DATABASE_URL`, y devolver un error claro. Excluir `.env` del control de versiones; documentar nombres y ejemplos seguros en `.env.example` cuando el repositorio vaya a compartirse.

## Contrato de respuestas HTTP

Todas las respuestas JSON usan el mismo sobre:

```json
{
  "data": {},
  "success": true,
  "status": 200,
  "errors": []
}
```

- `data` contiene el recurso o la lista en respuestas exitosas y vale `null` cuando falla.
- `success` es `true` solo para estados exitosos.
- `status` repite el código HTTP real como número.
- `errors` siempre es un arreglo. En éxito está vacío; en error contiene objetos con `message` y, cuando corresponde a una entrada concreta, `field`.

Centralizar este sobre en un paquete compartido de transporte, por ejemplo `internal/transport/http/response`. No duplicar structs ni codificadores en cada handler. Los errores de validación deben identificar el campo (`name`, `email`, `role_id`, `body` o `id`) para que el cliente pueda mostrarlos junto a la entrada correcta. No exponer errores internos ni SQL.

Los DTOs de request pertenecen a `internal/<feature>/dtos`; el handler los decodifica y los convierte explícitamente al modelo que recibe el servicio. El modelo de dominio no debe incorporar tags o campos solo por necesidades HTTP. Documentar en OpenAPI tanto el sobre común como los DTOs. Si una operación debe devolver el sobre JSON, no usar `204 No Content`; usar un estado que permita body, normalmente `200 OK`.

## Biblioteca estándar y almacenamiento

Usar `net/http`, `encoding/json`, `context`, `errors`, `log/slog`, `sync` y `os` según corresponda. No añadir frameworks, routers externos, ORM, contenedores de dependencias ni validadores.

Sin dependencias externas, elegir memoria o archivos según la durabilidad requerida. Proteger mapas compartidos con mutex y evitar exponer referencias mutables internas. Para archivos, manejar fallos de escritura y consistencia; no presentar almacenamiento de un solo proceso como seguro para múltiples procesos.

`database/sql` es estándar, pero no incluye drivers SQLite, PostgreSQL o MySQL. No prometer conexión funcional sin driver. Si se exige una base de datos y se prohíben dependencias, explicar la incompatibilidad y resolver esa elección con el usuario. Conservar drivers y datos existentes fuera del alcance solicitado; no eliminarlos ni sustituir persistencia por memoria silenciosamente. Una excepción explícita permite el driver necesario sin habilitar otros paquetes externos.

Para PostgreSQL, escribir SQL explícito con placeholders `$1`, `$2` y usar `QueryContext`, `QueryRowContext` o `ExecContext`. No crear tablas durante el arranque. Gestionar el esquema mediante migraciones versionadas `.up.sql` y `.down.sql`. Traducir errores concretos del driver a errores de dominio antes de llegar al handler.

En autenticación, una columna llamada `password` almacena siempre el hash y nunca la contraseña en texto plano. Usar un algoritmo específico de hash de contraseñas, como `golang.org/x/crypto/bcrypt`, y comparar con la función de la misma biblioteca. No incluir `password` en las respuestas JSON. Para JWT, usar una biblioteca mantenida como `github.com/golang-jwt/jwt/v5`, firmar con `JWT_SECRET`, exigir el algoritmo esperado y validar expiración, firma e issuer. Incluir el ID del usuario en `sub`. Enviar el JWT como `Authorization: Bearer <token>`. Colocar la validación Bearer y el ID autenticado en `internal/middleware` para que otras rutas puedan reutilizarlos.

## Política de testing

No crear archivos `*_test.go`, no modificar tests existentes y no ejecutar `go test`, cobertura, benchmarks ni pruebas manuales de endpoints salvo que el usuario pida testing explícitamente en la solicitud actual.

Una mención genérica a calidad, validación o buenas prácticas no autoriza tests. `gofmt`, `go build` y `go vet` son verificaciones de código y sí deben ejecutarse después de cambios de implementación.

Cuando el usuario solicite tests:

- Limitar los casos al alcance pedido.
- Colocarlos junto a la slice correspondiente.
- Usar `testing` y `net/http/httptest` de la biblioteca estándar.
- Preferir dobles pequeños que implementen la interfaz `Store`.
- No introducir frameworks de assertions o mocking.
- Informar exactamente qué comando se ejecutó y su resultado.

## Contratos y verificación

Pasar `context.Context` desde `r.Context()` a operaciones que puedan bloquear. Separar errores de dominio de detalles de almacenamiento y clasificarlos con `errors.Is`. El handler traduce entrada inválida, ausencia y conflictos a respuestas apropiadas, y registra fallos internos sin exponer detalles al cliente. Limitar el body y rechazar JSON inválido o documentos extra. Elegir las rutas según la versión de Go del proyecto y comprobar métodos no permitidos.

Para cambios de código, aplicar `gofmt` a los archivos modificados y ejecutar `go build ./...` y `go vet ./...`. Reportar bloqueos de toolchain o checks fallidos sin afirmar que pasaron. Para cambios exclusivos de la skill, validar sus metadatos y coherencia sin modificar ni arrancar la aplicación.
