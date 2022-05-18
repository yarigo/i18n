# i18n

i18n is a package that helps you translate Go programs into multiple languages.

## Features

- plural

## Installation

```sh
go install git.local/go/i18n@latest
```

## How to use

First you must create a json language file (see [`Translation file structure`](#file-structure) for details).

Load all translation files and set a fallback language. If you don't need a fallback language, you can set it as `language.Und`.

**i18n/en/main.json**

```json
[{ "id": "hello", "message": "Hello, world!" }]
```

**i18n/ru/main.json**

```json
[{ "id": "hello", "message": "Здравствуй, Мир!" }]
```

**main.go**

```go
package main

import (
  "fmt"
  "log"

  "github.com/yarigo/i18n"
  "golang.org/x/text/language"
)

func main() {
  // Create a new instance.
  t := i18n.New(&i18n.Config{Path: "./i18n", Fallback: language.Und})

  // Load all translation files.
  if err := t.Load(); err != nil {
    log.Fatalln(err.Error())
  }

  // Take printer for English language.
  en := t.Printer(language.English)
  // Take printer for Russian language.
  ru := t.Printer(language.Russian)

  // Print message.
  fmt.Println(en.Sprintf("hello"))
  // Print message.
  fmt.Println(ru.Sprintf("hello"))
}
```

```sh
$ go run main.go
Hello, world!
Здравствуй, Мир!
```

### plural

A selector matches an argument if:

- it is "other" or Other
- it matches the plural form of the argument: "zero", "one", "two", "few",
  or "many", or the equivalent Form
- it is of the form "=x" where x is an integer that matches the value of
  the argument.
- it is of the form "<x" where x is an integer that is larger than the
  argument.

For example:

```json
"rules": {
  "1": {
    "=0": "...",
    "one": "..."
  },
  "2": {
    ">10": "...",
    "other": "..."
  },
}
```

`1` and `2` its a index of variable.

## [Translation file structure](#file-structure)

One of the `message` or `rules` fields is required, but not both.

### Without plural

```json
[
  {
    "id": "unique message id",
    "message": "message text"
  }
]
```

### With plural

```json
[
  {
    "id": "unique message id",
    "rules": {
      "1": {
        "one": "first message id",
        "many": "message id %d"
      }
    }
  }
]
```

or

```json
[
  {
    "id": "unique message id",
    "rules": {
      "2": {
        "one": "message id %d and its %d text message",
        "many": "message id %d and message text %d"
      }
    }
  }
]
```

`1` and `2` are the number of the argument to use in the plural. If you only use the first argument, you can skip defining it like this:

```json
[
  {
    "id": "unique message id",
    "rules": {
      "one": "first message id",
      "many": "message id %d"
    }
  }
]
```
