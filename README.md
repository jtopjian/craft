# Craft

Craft is a library to help provision resources -- such as users, groups, and packages -- on a system.

While it's possible to use this library as-is, it's better suited as an underlying component to a more user-friendly tool.

This library is in very early stages and will likely have major changes.

## Development

This library is not vendored, so you will have to download dependent packages manually for now:

```shell
$ go get -u ./...
```

## Usage

If you really want to try using this, here's an example:

```go
package main

import (
        "fmt"

        "github.com/jtopjian/craft/client"
        "github.com/jtopjian/craft/resources/aptpkg"
        "github.com/sirupsen/logrus"
)

func main() {
        logger := logrus.New()
        logger.SetLevel(logrus.DebugLevel)

        c := client.Client{
                Logger: logger,
        }

        createOpts := aptpkg.CreateOpts{
                Name: "sl",
        }

        err := aptpkg.Create(c, createOpts)
        if err != nil {
                panic(err)
        }

        exists, err := aptpkg.Exists(c, "sl")
        if err != nil {
                panic(err)
        }

        fmt.Printf("Package sl exists: %t\n", exists)

        pkg, err := aptpkg.Read(c, "sl")
        if err != nil {
                panic(err)
        }

        fmt.Printf("Package sl: %#v\n", pkg)

        err = aptpkg.Delete(c, "sl")
        if err != nil {
                panic(err)
        }
}
```

Which will yield the following output:

```
DEBU[0000] Installing package
DEBU[0000] Package Create Options: aptpkg.CreateOpts{Name:"sl", Version:""}
DEBU[0005] Checking if package sl exists
DEBU[0005] Reading package sl
Package sl exists: true
DEBU[0005] Reading package sl
Package sl: aptpkg.AptPkg{Name:"sl", Version:"3.03-17build1", LatestVersion:"3.03-17build1"}
DEBU[0005] Deleting package sl
```
