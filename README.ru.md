# freeport

[![en](https://img.shields.io/badge/lang-en-green.svg)](README.md)

Пакет для получения незанятых портов.

## Примеры использования

### 1. Get / MustGet
```go
package main

import (
    "fmt"
    "github.com/kiteggrad/freeport"
)

func main() {
    port, err := freeport.Get()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Obtained free port:", port)

    // or

    port = freeport.MustGet()
    fmt.Println("Obtained free port:", port)
}
```

### 2. Запросить несколько портов сразу (до их использования)

Функции `freeport.Get()` и `freeport.MustGet()` используют глобальный `freeport.Generator` который запоминает какие порты были из него выданы и выдаёт только не запрошенные из него ранее (и не занятые на момент вызова функции).

```go
package main

import (
    "fmt"
    "github.com/kiteggrad/freeport"
)

func main() {
    port1 := freeport.MustGet()

    // safe - freeport запомнил port1 и не вернёт его снова
    port2 := freeport.MustGet()

    fmt.Println(port1 == port2) // false
}
```

### 3. Retry

Функции `freeport.Retry()`, `freeport.RetryBackoff()` и тп. используют глобальный `freeport.Generator` который запоминает какие порты были из него выданы и выдаёт только не запрошенные из него ранее (и не занятые на момент вызова функции).

Под капотом используется [backoff](https://github.com/cenkalti/backoff) со свойственным ему поведением.

```go
package main

import (
    "fmt"
    "github.com/kiteggrad/freeport"
)

func main() {
    port, err := freeport.Retry(func(port int) error {
        if port % 2 == 0 {
            return fmt.Error("some possible error, port", port) // some possible error, port 8081
        }

        fmt.Println("connect to port", port) // connect to port 8081

        return nil
    })

    fmt.Println("succesfully connected to port", port) // succesfully connected to port 8081
}
```

## Возможные проблемы

### 1. Возможны коллизии при использовании freeport.Get() / freeport.MustGet()

#### 1.1. Если запросить очень много портов сразу

Колл-во неудачных попыток занять случайный не занятый порт ограничено.
Так что с каждым новым запросом (до использования полученных портов) шанс получить ошибку увеличивается (тест на 1000 запросов всегда проходит удачно).
У меня ошибки начинают возникать где-то примерно с 3500 запрошенных портов.

```go
package some

import (
    "github.com/kiteggrad/freeport"
)

func main() {
	for i := 0; i < 4000; i++ {
		_ = freeport.MustGet() // panic
	}
}
```

#### 1.2. В параллельно выполняющихся тестах

Следует помнить что в разных пакетах при выполнении тестов для каждого импортируемого пакета функция `init` и глобальные переменные инициализируются заново. 
Следовательно глобальные `freeport.Generator` в разных пакетах при выполнении тестов могут выдавать одинаковые порты.

- internal/pkg1/pkg1_test.go
    ```go
    package pkg1

    import (
	    "testing"
        "fmt"
        "github.com/kiteggrad/freeport"
    )

    func Test(t *testing.T) {
        port := generator.MustGet()
        fmt.Println(port) // 8080
    }
    ```
- internal/pkg1/pkg2_test.go
    ```go
    package pkg2

    import (
	    "testing"
        "fmt"
        "github.com/kiteggrad/freeport"
    )

    func Test(t *testing.T) {
        port := generator.MustGet()
        fmt.Println(port) // 8080
    }
    ```

Одно из возможных решений этой проблемы - использовать в параллельно выполняющихся тестах общий экземпляр `freeport.Generator`. 
Или использовать `freeport.Retry`.