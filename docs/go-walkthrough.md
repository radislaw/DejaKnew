# Go Walkthrough — FlashCards API

Разбор Go-кода бэкенда для разработчика, который знает Vue/TypeScript,
но Go видит впервые. Идём по файлам в том порядке, в котором они были
написаны, и на каждом объясняем не только «что», но и Go-концепции,
которые в нём впервые встречаются.

Читать сверху вниз: каждая часть опирается на предыдущую.

---

## Часть 0 — Карта различий Go и TS/JS

Прежде чем смотреть код, держи в голове главные отличия:

| Тема | TypeScript | Go |
|---|---|---|
| Видимость | `export` / `private` | Заглавная буква имени = публично, строчная = приватно |
| Классы | `class`, `extends` | Структуры (`struct`) + методы. Наследования нет, есть композиция |
| Интерфейсы | реализуются явно: `class X implements Y` | реализуются **неявно** — тип подходит под интерфейс, если у него есть нужные методы |
| Ошибки | `throw` / `try/catch` | ошибка — обычное значение, возвращается из функции |
| null | `null` / `undefined` | `nil` |
| Асинхронность | `async/await` | горутины и каналы (в этом проекте пока не нужны) |
| Менеджер пакетов | npm + package.json | go-модуль + go.mod |
| Типы значений | всё ссылочное, кроме примитивов | структуры копируются по значению; ссылка — через указатель (`*`) |

Последняя строка — самая непривычная. К ней вернёмся в Части 4.

---

## Часть 1 — Модуль и структура проекта

### go.mod — аналог package.json

```
module flashcards/api

go 1.25.0

require github.com/gin-gonic/gin v1.12.0
```

- `module flashcards/api` — имя модуля. Все импорты внутри проекта
  начинаются с него: `flashcards/api/internal/config`.
- `go 1.25.0` — минимальная версия Go.
- `require` — прямые зависимости. Блок `// indirect` ниже — транзитивные,
  Go дописывает их сам командой `go mod tidy`.
- `go.sum` — как `package-lock.json`: контрольные суммы для воспроизводимости.

### Пакеты, файлы, папки

В Go **единица кода — пакет (package)**, а не файл.

- Все `.go`-файлы в одной папке = один пакет. Первая строка файла
  объявляет имя: `package config`.
- Файлы одного пакета видят функции и типы друг друга **без импортов** —
  они в общем пространстве имён.
- Импортируется пакет целиком, по пути: `import "flashcards/api/internal/config"`,
  дальше обращение `config.Load()`.

### Папка `internal/` — приватность на уровне компилятора

Папка с именем `internal` — особая. Пакеты внутри неё может импортировать
**только** код, лежащий рядом с `internal` или ниже. Снаружи модуля
импорт `flashcards/api/internal/...` не скомпилируется.

Это бесплатная инкапсуляция: всё, что не должно быть публичным API
библиотеки, кладём в `internal/`. У нас весь бэкенд там, кроме точки
входа `cmd/server`.

---

## Часть 2 — config: внешние настройки

Файл `internal/config/config.go`.

```go
package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}
```

- `type Config struct { ... }` — объявление структуры. Структура — это
  именованный набор полей. Ближайший аналог в TS — `interface` для
  объекта, но `struct` это конкретный тип данных, а не контракт.
- Поля `Port`, `DatabaseURL`, `JWTSecret` — с **заглавной** буквы,
  значит публичные: их видно из других пакетов.

```go
func Load() Config {
	return Config{
		Port:      env("PORT", "8080"),
		...
	}
}
```

- `func Load() Config` — функция без аргументов, возвращает `Config`.
- `Config{ Port: ..., ... }` — создание значения структуры (literal).
  Возвращаем структуру **по значению** — вызывающий получает свою копию.

```go
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
```

- `env` — со **строчной** буквы → функция приватна для пакета `config`.
- `key, fallback string` — два аргумента одного типа, тип пишется один раз.
- `if v := os.Getenv(key); v != "" { }` — типичный для Go приём:
  объявить переменную (`v :=`) прямо в `if`, она живёт только внутри `if`.
- `:=` — короткое объявление: тип выводится автоматически. Аналог `const`
  с выводом типа. Полная форма — `var v string = ...`.

---

## Часть 3 — domain: модели и ошибки

### user.go — структура-модель

```go
package domain

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}
```

Новое здесь — **струк-теги** (backtick-строка после поля):

```
`gorm:"type:uuid;primaryKey" json:"id"`
```

Это метаданные. Сам Go их не использует — их читают библиотеки через
рефлексию:
- `gorm:"..."` — как GORM отобразит поле в колонку БД.
- `json:"id"` — под каким именем поле попадёт в JSON.
- `json:"-"` у `PasswordHash` — **никогда не сериализовать в JSON**.
  Хэш пароля не утечёт в ответ API, даже если случайно отдать `User`
  целиком. Это защита на уровне типа.

### errors.go — sentinel-ошибки

```go
package domain

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotFound           = errors.New("not found")
)
```

- В Go **ошибка — это значение**, а не исключение. Нет `throw`.
- `var ( ... )` — блок объявления переменных уровня пакета.
- `errors.New("...")` создаёт ошибку. Присваиваем в переменную с
  префиксом `Err...` — это **sentinel-ошибка**: заранее объявленный
  «маяк», который слои сравнивают между собой.
- Идея: service возвращает `domain.ErrUserExists`, handler проверяет
  `errors.Is(err, domain.ErrUserExists)` и отдаёт HTTP 409. Так слои
  общаются об ошибках, не завися от текста сообщения.

---

## Часть 4 — auth: примитивы безопасности

### password.go — функции и множественный возврат

```go
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
```

- `(string, error)` — функция возвращает **два значения**. Это норма Go:
  «результат + ошибка».
- `[]byte(plain)` — приведение `string` к срезу байт (`[]byte`). bcrypt
  работает с байтами. `string(hash)` — обратное приведение.
- `hash, err := ...` — приняли оба возвращённых значения.
- `if err != nil { return "", err }` — **канонический паттерн Go**.
  После каждого вызова, который может упасть, проверяем `err`. Если
  ошибка — возвращаем «пустой» результат и саму ошибку наверх.
- `nil` — отсутствие значения. `return string(hash), nil` — «вот
  результат, ошибки нет».

```go
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
```

bcrypt-сравнение возвращает `error`: `nil` если пароль верный. Поэтому
`== nil` превращаем в `bool`.

### jwt.go — структура как «класс», методы, конструктор

Здесь впервые структура с поведением.

```go
type TokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: []byte(secret)}
}
```

- `TokenManager` хранит секрет. Поле `secret` со строчной буквы —
  приватное, снаружи пакета его не достать.
- `NewTokenManager` — **конструктор**. В Go нет ключевого слова
  `constructor`; есть соглашение: функция `NewТип`, которая создаёт и
  возвращает значение. Это просто функция.
- `*TokenManager` в возврате — **указатель** на структуру.
- `&TokenManager{...}` — оператор `&` берёт адрес созданного значения,
  то есть возвращаем указатель, а не копию.

**Указатели — ключевой момент.** В Go структура по умолчанию копируется
при передаче. Если бы `NewTokenManager` возвращал `TokenManager` (без
`*`), каждый, кто его получил, имел бы свою копию. Указатель `*T` — это
«ссылка на один и тот же объект в памяти». Аналогия с JS: там объекты
всегда передаются по ссылке; в Go ты выбираешь сам — `T` (копия) или
`*T` (ссылка).

```go
func (m *TokenManager) Issue(userID string) (string, error) {
	...
}
```

- `(m *TokenManager)` перед именем функции — **получатель (receiver)**.
  Это и делает функцию **методом** структуры. Читается: «метод `Issue`
  у `*TokenManager`».
- `m` — имя, под которым внутри метода доступна сама структура (как
  `this` в JS, но имя выбираешь ты, обычно короткое).
- `*TokenManager` (указатель) в получателе — метод работает с тем же
  экземпляром, а не с копией. Стандартный выбор для методов.
- Вызов: `tokens.Issue("user-1")`.

`Verify` устроен так же — разбирает токен, проверяет подпись, отдаёт
`(userID, error)`.

---

## Часть 5 — repository: доступ к БД

### db.go — подключение, миграции, оборачивание ошибок

```go
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	return db, nil
}
```

- `*gorm.DB` — указатель на структуру GORM, наш «хэндл» БД.
- `fmt.Errorf("connect to postgres: %w", err)` — создаёт **новую**
  ошибку с добавленным контекстом. Глагол `%w` («wrap») — особый:
  он не просто вставляет текст, а **вкладывает** исходную ошибку
  внутрь новой. Позже `errors.Is` сможет «развернуть» цепочку и найти
  оригинал. Без `%w` (например `%v`) контекст добавится, но связь с
  исходной ошибкой потеряется.

```go
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	return nil
}
```

- `AutoMigrate` смотрит на структуры и создаёт/правит таблицы под них.
- `&domain.User{}` — указатель на пустое значение `User`; GORM нужен
  тип, чтобы прочитать его поля и теги.
- **Каждую новую модель надо дописывать сюда** — это единственное место,
  где регистрируется схема.

### user_repo.go — методы репозитория и errors.Is

```go
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}
```

Тот же паттерн, что у `TokenManager`: структура хранит зависимость
(`*gorm.DB`), конструктор `NewUserRepo` её внедряет.

```go
func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

- `var user domain.User` — объявили переменную. Без присваивания она
  получает **нулевое значение**: все поля пустые (`""`, `0`, `nil`).
  В Go нет «неинициализированных» переменных.
- `r.db.Where(...).First(&user)` — GORM-запрос. `&user` передаём, чтобы
  GORM **записал результат по адресу** нашей переменной.
- `.Error` — GORM складывает ошибку запроса в это поле.
- `errors.Is(err, gorm.ErrRecordNotFound)` — проверяем: это ошибка
  «не найдено»? Если да — **переводим её в нашу доменную** `ErrNotFound`.
  Так слой repository прячет детали GORM: service не знает про GORM,
  он знает только `domain.ErrNotFound`.
- `return &user, nil` — отдаём указатель на найденного пользователя.

---

## Часть 6 — service: бизнес-логика и интерфейсы

Это самый важный для понимания Go файл — `service/auth_service.go`.

### Интерфейс, объявленный у потребителя

```go
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
}
```

- `interface` в Go — это **список методов**. Тип «удовлетворяет»
  интерфейсу, если у него есть все эти методы. **Никакого `implements`**.
- `repository.UserRepo` нигде не написал «я UserRepository». Но у него
  есть `Create`, `FindByEmail`, `FindByID` с нужными сигнатурами —
  значит он **автоматически** подходит. Это «структурная типизация»,
  она же «duck typing, проверяемый компилятором».
- Почему интерфейс объявлен **здесь**, в `service`, а не в `repository`?
  Принцип Go: **интерфейс принадлежит тому, кто его использует.**
  `service` говорит «мне нужен кто-то, кто умеет вот это». `repository`
  про этот интерфейс даже не знает. Это разворачивает зависимость:
  service не зависит от repository, оба зависят от абстракции.

### Почему это вообще нужно

```go
type AuthService struct {
	users  UserRepository   // <- интерфейс, не конкретный *UserRepo
	tokens *auth.TokenManager
}

func NewAuthService(users UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}
```

`AuthService` хранит `users` **как интерфейс**. В `main.go` туда передадут
настоящий `*repository.UserRepo` (с живой БД). А в тестах — поддельный
`fakeUserRepo` (с map в памяти). `AuthService` разницы не замечает —
ему достаточно, что у объекта есть нужные методы.

Это и есть **внедрение зависимостей (DI)**. В Go нет DI-фреймворка:
зависимости — просто аргументы конструктора. Вся «магия» — в том, что
тип аргумента это интерфейс.

### Логика Register

```go
func (s *AuthService) Register(email, password string) (string, *domain.User, error) {
	existing, err := s.users.FindByEmail(email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return "", nil, err
	}
	if existing != nil {
		return "", nil, domain.ErrUserExists
	}
	...
}
```

- Возврат из **трёх** значений: `(token, user, error)`.
- Проверка тонкая: `FindByEmail` вернёт `ErrNotFound`, если почты нет —
  и это **хорошо** (значит можно регистрировать). Поэтому: «ошибка
  есть И это НЕ ErrNotFound» → настоящая проблема, выходим. Если же
  пользователь нашёлся (`existing != nil`) → почта занята.
- Дальше: хэшируем пароль, генерируем `uuid.NewString()` для ID,
  `s.users.Create(user)`, выдаём токен. Каждый шаг — с проверкой `err`.

`Login` симметричен: найти юзера, `auth.CheckPassword`, выдать токен.
Неверная почта и неверный пароль оба возвращают `ErrInvalidCredentials`
— намеренно не подсказываем атакующему, что именно не так.

---

## Часть 7 — handler и middleware: HTTP-слой

### auth_handler.go — Gin-обработчики

```go
type credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
```

- Структура только для разбора тела запроса. Строчная буква (`credentials`)
  — она приватна для пакета handler.
- Тег `binding:"required,email"` читает Gin: поле обязательное и должно
  быть валидной почтой. `min=8` — пароль не короче 8 символов. Валидация
  декларативная, без ручных `if`.

```go
func (h *AuthHandler) Register(c *gin.Context) {
	var body credentials
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := h.auth.Register(body.Email, body.Password)
	if errors.Is(err, domain.ErrUserExists) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	...
}
```

- `c *gin.Context` — объект запроса/ответа Gin (запрос, параметры,
  методы записи ответа).
- `c.ShouldBindJSON(&body)` — распарсить JSON тела в структуру `body`
  (по адресу `&body`) и заодно прогнать `binding`-валидацию.
- `gin.H{"error": ...}` — `gin.H` это просто `map[string]any`, удобный
  тип для JSON-ответа.
- Обработчик **тонкий**: распарсил → позвал service → разложил ошибки
  по HTTP-статусам (`errors.Is` → 409 / 401 / 500). Бизнес-логики нет —
  она вся в service. Это и есть «слои не перепрыгивать».

### middleware/auth.go — функция, возвращающая функцию

```go
func RequireAuth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "..."})
			return
		}
		userID, err := tokens.Verify(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "..."})
			return
		}
		c.Set(ContextUserID, userID)
		c.Next()
	}
}
```

- `gin.HandlerFunc` — это тип `func(*gin.Context)`. В Go функции —
  обычные значения, у них есть типы, их можно возвращать.
- `RequireAuth` **возвращает функцию** — это **замыкание (closure)**:
  внутренняя функция «помнит» `tokens` из внешней. Так мы передаём
  middleware зависимость (`TokenManager`), хотя сам Gin вызовет её
  просто как `func(c)`.
- `strings.CutPrefix` отрезает `"Bearer "`. Возвращает `(остаток, ok)`,
  где `ok bool` — был ли префикс. Множественный возврат и тут.
- `c.Set(ContextUserID, userID)` — кладём ID юзера в контекст запроса.
- `c.Next()` — пропустить дальше по цепочке. `c.AbortWithStatusJSON` —
  оборвать цепочку и ответить ошибкой.

Дальше `Me` достаёт ID: `c.GetString(middleware.ContextUserID)`.

---

## Часть 8 — router и main: склейка

### router.go — группировка маршрутов

```go
func Router(tokens *auth.TokenManager, authHandler *AuthHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { ... })

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(tokens))
		{
			protected.GET("/me", authHandler.Me)
		}
	}
	return r
}
```

- `r.Group("/api/v1")` — общий префикс пути.
- `protected.Use(middleware.RequireAuth(tokens))` — на эту подгруппу
  навешан middleware: все её маршруты сначала пройдут проверку токена.
  `/register` и `/login` остаются открытыми.
- Фигурные скобки `{ }` тут — просто визуальная группировка (новый
  блок видимости), Gin их не требует. Чисто читаемость.

### main.go — ручное внедрение зависимостей

```go
func main() {
	cfg := config.Load()

	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("migration: %v", err)
	}

	tokens := auth.NewTokenManager(cfg.JWTSecret)
	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthService(userRepo, tokens)
	authHandler := handler.NewAuthHandler(authService)

	r := handler.Router(tokens, authHandler)
	r.Run(":" + cfg.Port)
}
```

`main` — точка входа (пакет `main`, функция `main`). Здесь видна вся
**цепочка сборки снизу вверх**:

```
db ──► userRepo ──► authService ──► authHandler ──► Router
                tokens ──┘                  └── tokens ──┘
```

Каждый уровень получает свои зависимости в конструкторе. Никакого
DI-контейнера — просто по очереди создаём объекты и передаём предыдущие
в следующие. Это вся «архитектура».

- `log.Fatalf` — напечатать и **завершить процесс**. Если БД недоступна
  или миграция упала — сервер не должен стартовать молча. Падаем громко.
- Миграция вызвана **до** `r.Run` — таблицы готовы раньше, чем примем
  первый запрос.

---

## Часть 9 — тесты

Файл `service/auth_service_test.go`.

- Файлы с тестами заканчиваются на `_test.go`. Они компилируются только
  при `go test`, в обычную сборку не попадают.
- Тест — функция `func TestXxx(t *testing.T)`.
- `t.Fatalf(...)` — пометить тест упавшим и остановить его.

### Поддельный репозиторий

```go
type fakeUserRepo struct {
	byEmail map[string]*domain.User
	byID    map[string]*domain.User
}

func (f *fakeUserRepo) Create(u *domain.User) error { ... }
func (f *fakeUserRepo) FindByEmail(email string) (*domain.User, error) { ... }
func (f *fakeUserRepo) FindByID(id string) (*domain.User, error) { ... }
```

`fakeUserRepo` хранит юзеров в `map` (аналог JS-объекта/Map) вместо БД.
У него есть **те же три метода**, что в интерфейсе `UserRepository` —
значит он этому интерфейсу удовлетворяет. `NewAuthService` спокойно
принимает его вместо настоящего репозитория.

Вот зачем нужен был интерфейс из Части 6: тесты `AuthService` гоняют
**настоящую бизнес-логику** (хэширование, проверку дублей, выдачу
токена), но **без БД** — быстро и изолированно. Реальную работу с
PostgreSQL проверяем отдельно, интеграционными тестами на repository.

---

## Что освоено к этому моменту

Пройденные Go-концепции:

- пакеты, модуль, `internal/`, видимость по регистру буквы
- структуры, методы и получатели (`func (r *T) ...`)
- указатели: `*T` (тип-указатель), `&x` (взять адрес)
- множественный возврат, `error` как значение, паттерн `if err != nil`
- sentinel-ошибки, `errors.Is`, оборачивание через `%w`
- **интерфейсы и неявная реализация** — главное отличие от TS
- внедрение зависимостей через конструкторы
- струк-теги (`gorm`, `json`, `binding`) и кто их читает
- замыкания (middleware)
- нулевые значения, `:=` vs `var`
- тесты `_test.go`, поддельные реализации интерфейса

## Что дальше (следующие под-шаги 5)

- 5b — `Deck` CRUD: те же слои (domain → repository → service → handler),
  скопированные с auth как с шаблона.
- Позже — Go-порт SRS-алгоритма, study-эндпоинты.

Когда будешь читать новый код проекта — ищи в нём эти же паттерны:
структура + конструктор, интерфейс у потребителя, тонкий handler,
`if err != nil`. Они повторяются везде намеренно.
