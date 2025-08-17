# Использование ролей в JWT токенах

## Что было добавлено

1. **Поле `role` в JWT токен** - теперь при генерации access token включается роль пользователя
2. **Новый метод `ValidateTokenAndExtractClaims`** - извлекает все данные из токена, включая роль
3. **Обновленные middleware** - автоматически извлекают роль и добавляют в контекст Gin

## Структура TokenClaims

```go
type TokenClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}
```

## Данные в контексте Gin

После аутентификации в контексте доступны:
- `c.Get("userID")` - ID пользователя (int64)
- `c.Get("username")` - имя пользователя (string)
- `c.Get("email")` - email пользователя (string)
- `c.Get("role")` - роль пользователя (string)

## Пример использования в обработчике

```go
func ExampleHandler(c *gin.Context) {
    // Получаем роль из контекста
    role, exists := c.Get("role")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
        return
    }

    userRole := role.(string)
    
    // Проверяем роль
    switch userRole {
    case "admin":
        // Логика для администратора
        c.JSON(http.StatusOK, gin.H{
            "message": "Admin access granted",
            "data": "sensitive admin data"
        })
    case "user":
        // Логика для обычного пользователя
        c.JSON(http.StatusOK, gin.H{
            "message": "User access granted",
            "data": "regular user data"
        })
    default:
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Unknown role"
        })
    }
}
```

## Пример middleware для конкретной роли

```go
func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
            c.Abort()
            return
        }

        if role.(string) != requiredRole {
            c.JSON(http.StatusForbidden, gin.H{
                "error": fmt.Sprintf("Required role: %s, got: %s", requiredRole, role)
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Использование в роутере
router.GET("/admin-only", RequireRole("admin"), AdminOnlyHandler)
```

## Пример JWT токена (декодированный)

**Header:**
```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

**Payload:**
```json
{
  "sub": 123,
  "username": "john_doe",
  "email": "john@example.com",
  "role": "admin",
  "exp": 1642723200,
  "iat": 1642636800
}
```

## Обновленные middleware

### AuthMiddleware
- Извлекает все данные пользователя из токена
- Добавляет userID, username, email, role в контекст

### AdminMiddleware  
- Проверяет роль прямо из токена (без запроса к БД)
- Разрешает доступ только пользователям с ролью "admin"

## Преимущества

1. **Производительность** - роль берется из токена, не нужны запросы к БД
2. **Безопасность** - роль подписана в токене, нельзя подделать
3. **Простота** - легко проверить роль в любом обработчике
4. **Централизованность** - роли управляются через JWT

## Использование в job handler

```go
func (h *JobHandler) ExampleMethod(c *gin.Context) {
    userID, _ := c.Get("userID")
    role, _ := c.Get("role")
    
    // Проверяем права доступа
    if role.(string) == "admin" {
        // Администратор может видеть все работы
    } else {
        // Обычный пользователь видит только свои работы
        // Используем userID для фильтрации
    }
}
```