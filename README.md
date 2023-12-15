# freeport

Translated from [![ru](https://img.shields.io/badge/lang-ru-red.svg)](README.ru.md) with the help of an online translator.

Package for getting unused ports.

## Examples of use

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

    fmt.Println("Getted free port:", port)

    // or

    port = freeport.MustGet()
    fmt.Println("Getted free port:", port)
}
```

### 2. Request multiple ports at once (before using them)

The `freeport.Get()` and `freeport.MustGet()` functions use the global `freeport.Generator` which remembers which ports were issued from it and issues only the ports not previously requested from it (and not occupied at the time of the function call).

```go
package main

import (
    "fmt"
    "github.com/kiteggrad/freeport"
)

func main() {
    port1 := freeport.MustGet()

    // safe - freeport has memorized port1 and will not return it again
    port2 := freeport.MustGet()

    fmt.Println(port1 == port2) // false
}
```

### 3. Retry

Functions `freeport.Retry()`, `freeport.RetryBackoff()`, etc. use global `freeport.Generator` which remembers which ports were retrieved from it and retrieves only ports not requested from it earlier (and not occupied at the moment of function call).

Under the hood, [backoff](https://github.com/cenkalti/backoff) is used with its inherent behavior.

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

## Possible problems

### 1. Possible collisions when using freeport.Get() / freeport.MustGet()

#### 1.1. If you request very many ports at once

The number of unsuccessful attempts to take a random unoccupied port is limited.
So with each new request (before using the getted ports) the chance to get an error increases (the test for 1000 requests is always successful).
I start getting errors somewhere around 3500 ports requested.

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

#### 1.2. In parallel tests

It should be remembered that in different packages the `init` function and global variables are initialized anew for each imported package when executing tests. 
Consequently, global `freeport.Generator` in different packages may generate the same ports when running tests.

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

One possible solution to this problem is to use a shared instance of `freeport.Generator` in parallel running tests. 
Or use `freeport.Retry`.