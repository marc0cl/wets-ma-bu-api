# API Restaurant - Documentación

Esta es una API para la gestión de usuarios y restaurantes, con autenticación JWT y control de acceso basado en roles.

## Documentación Swagger

La documentación de la API está disponible en formato Swagger. Puedes acceder a ella de las siguientes maneras:

1. **Swagger UI**: Accede a `/swagger` para ver la interfaz interactiva de la API
2. **Swagger JSON**: Accede a `/swagger.json` para obtener la definición de la API en formato JSON
3. **Swagger YAML**: Accede a `/swagger.yaml` para obtener la definición de la API en formato YAML

## Endpoints principales

### Autenticación

- **POST /api/v1/auth/register**: Registrar un nuevo usuario
- **POST /api/v1/auth/login**: Iniciar sesión y obtener un token JWT

### Usuarios

- **GET /api/v1/users/{id}**: Obtener información de un usuario
- **PUT /api/v1/users/{id}**: Actualizar información de un usuario
- **DELETE /api/v1/users/{id}**: Eliminar un usuario

### Restaurantes

- **POST /api/v1/restaurants**: Crear un nuevo restaurante
- **PUT /api/v1/restaurants/{id}**: Actualizar un restaurante
- **DELETE /api/v1/restaurants/{id}**: Eliminar un restaurante
- **GET /api/v1/users/{userId}/restaurants**: Obtener todos los restaurantes de un usuario
- **GET /api/v1/users/{userId}/restaurants/{id}**: Obtener un restaurante específico de un usuario

## Autenticación

Para acceder a las rutas protegidas, es necesario incluir un token JWT en el encabezado de autorización:

```
Authorization: Bearer [token]
```

El token se obtiene al iniciar sesión con la ruta `/api/v1/auth/login`.

## Control de acceso

- Los usuarios regulares solo pueden acceder y modificar sus propios recursos
- Los administradores pueden acceder y modificar todos los recursos

## Modelos

### Usuario
```json
{
  "id": 1,
  "name": "Nombre Usuario",
  "email": "usuario@ejemplo.com",
  "role": "user"
}
```

### Restaurante
```json
{
  "id": 1,
  "name": "Nombre Restaurante",
  "address": "Dirección",
  "phone": "123456789",
  "user_id": 1
}
```

## Ejemplo de uso

### Registro de usuario

```bash
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ejemplo Usuario",
    "email": "usuario@ejemplo.com",
    "password": "contraseña123",
    "role": "user"
  }'
```

### Iniciar sesión

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "contraseña123"
  }'
```

### Crear un restaurante

```bash
curl -X POST http://localhost:8000/api/v1/restaurants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer [token]" \
  -d '{
    "name": "Restaurante Ejemplo",
    "address": "Calle Principal 123",
    "phone": "123456789"
  }'
```

### Obtener restaurantes de un usuario

```bash
curl -X GET http://localhost:8000/api/v1/users/1/restaurants \
  -H "Authorization: Bearer [token]"
```
