Маленькая CLI-утилита на Go. По ссылке на Git-репозиторий она скачивает
файл `go.mod`, печатает имя модуля и
версию Go, а затем сверяет каждую зависимость с публичным
`Go module proxy` и сообщает, какие из них можно
обновить.

Опционально умеет интерактивно поднимать версии и сохранять получившийся
`go.mod` в локальный файл.

## Возможности

- Берёт `go.mod` напрямую через `raw.githubusercontent.com`,
  клонировать репозиторий не нужно.
- Запрашивает `https://proxy.golang.org/<module>/@latest` для каждой
  зависимости параллельно.
- В режиме `--update` интерактивно поднимает версии и пишет новый `go.mod`.

## Использование

```text
Использование: ./mws-project [флаги] <git-repo-url>

Флаги:
  -branch string
        Ветка, тег или коммит, из которого берётся go.mod (по умолчанию "HEAD")
  -o string
        Путь для сохранения обновлённого go.mod (используется с --update)
        (по умолчанию "go.mod.updated")
  -update
        Интерактивно обновить устаревшие require и записать go.mod локально
```
URL репозитория можно передавать в любой из форм:

- `https://github.com/spf13/cobra`
- `https://github.com/spf13/cobra.git`
- `github.com/spf13/cobra`

Сейчас поддерживается только `github.com`.

## Примеры

### 1. Посмотреть устаревшие зависимости

```sh
$ ./mws-project https://github.com/spf13/cobra
Загружаю https://raw.githubusercontent.com/spf13/cobra/HEAD/go.mod

Модуль:        github.com/spf13/cobra
Версия Go:     1.15
Зависимостей:  4

МОДУЛЬ                               ТЕКУЩАЯ  ПОСЛЕДНЯЯ
-----------------------------------------------------------
github.com/cpuguy83/go-md2man/v2     v2.0.3   v2.0.4   (есть обновление)
github.com/inconshreveable/mousetrap v1.1.0   v1.1.0
github.com/spf13/pflag               v1.0.5   v1.0.5
gopkg.in/yaml.v3                     v3.0.1   v3.0.1

Можно обновить: 1 модул(я/ей).
```

### 2. Анализ конкретной ветки или тега

```sh
./mws-project --branch v1.8.0 github.com/spf13/cobra
```

### 3. Интерактивное обновление

```sh
$ ./mws-project --update -o ./go.mod.new https://github.com/spf13/cobra
...
Можно обновить: 1 модул(я/ей).
Обновить github.com/cpuguy83/go-md2man/v2: v2.0.3 -> v2.0.4 ? [y/N]: y
Обновлённый go.mod (1 изменени(й)) записан в ./go.mod.new
```

Оригинальный `go.mod` upstream-проекта не меняется, обновлённый манифест
сохраняется локально, чтобы его можно было сравнить перед применением.

## Структура проекта

```
.
├── main.go                       парсинг флагов и точка входа
├── internal/outdated/
│   ├── gitrepo.go                URL Git-репозитория -> raw-URL go.mod
│   ├── proxy.go                  обёртка над net/http и клиент proxy.golang.org
│   ├── check.go                  скачивание go.mod и параллельная проверка
│   ├── report.go                 печать заголовка и таблицы
│   └── update.go                 интерактивное обновление go.mod
├── go.mod
└── README.md
```
