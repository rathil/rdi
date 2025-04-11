# rDI – Dependency Injection for Go 🥷

<img src="assets/mascot.png" align="right" width="100" alt="rDI Mascot">
rDI is a lightweight and intuitive dependency injection container for Go. It is designed to be flexible, testable, and unobtrusive, letting you wire up your dependencies cleanly and explicitly.

## Installation

```sh
$ go get -u github.com/rathil/rdi
```
Then import it:
```go
import "github.com/rathil/rdi"
```

## Quick start

### Custom functionality

```go
type IUser interface{}
type User struct{}

func NewIUser() (IUser, error) {
    return &User{}, nil
}
func NewUserPoint() (*User, error) {
    return &User{}, nil
}
func NewUser() User {
    return User{}
}

type IDevice interface{}
type Device struct{
    User *User
}

var sameErrorDevice = fmt.Errorf(`same device error`)

func NewDevicePoint(user *User, sameVar int) (*Device, error) {
    if sameVar > 5 {
        return nil, sameErrorDevice
    }
    return &Device{user}, nil
}
func NewDevice(user *User) (Device, error) {
    return Device{user}, nil
}
func NewIDevicePoint(user *User) (IDevice, error) {
    return &Device{user}, nil
}
```

### Initializing a container

```go
di := standard.New()
```

### Provide _(*User, User)_ dependency

```go
user := NewUserPoint()
// ...
if err := di.Provide(user); err != nil {
    panic(err)
}
```
Or using a constructor directly:
```go
if err := di.Provide(NewUserPoint); err != nil {
    panic(err)
}
```
Or with MustProvide:
```go
di.MustProvide(NewIUser)
di.MustProvide(NewUserPoint)
di.MustProvide(NewUser)
```

### Incorrect use of _IUser_ interface type

```go
user := NewIUser()
// ...
di.MustProvide(user) // Will register *User instead of IUser
```

### Correct way to provide _IUser_ interface
```go
di.MustProvide(NewIUser)
```

### Invoking _(IUser, *User, User)_ dependencies

```go
if err := di.Invoke(func(user *User) {
    // ...
}); err != nil {
    panic(err)
}
```
With error handling:
```go
var someError = fmt.Errorf(`some error`)
// ...
if err := di.Invoke(func(user *User) error {
    // ...
    if varName > 5 {
        return someError
    }
    return nil
}); err != nil {
    if !errors.Is(err, someError) {
        panic(err)
    }
}
```
Chaining multiple invocations:
```go
di.
	MustInvoke(func(user IUser) {
        // ...
    }).
	MustInvoke(func(user *User) {
        // ...
    }).
	MustInvoke(func(user User) {
        // ...
    })
```

### Provider Options

By default, providers are singleton and cached after first creation:

```go
di := standard.New().
    MustProvide(func() context.Context { // Will be called once
        return context.Background()
    })

di.MustInvoke(func(c context.Context) {})
di.MustInvoke(func(c context.Context) {})
```

To get a new instance on each request, use **WithTransient()**:

```go
di := standard.New().
    MustProvide(func() context.Context { // It will be called every time when the context is requested
        return context.Background()
    }, rdi.WithTransient())

di.MustInvoke(func(c context.Context) {})
di.MustInvoke(func(c context.Context) {})
```

### Container Copying and Overrides

You can create a child container to override or extend dependencies:

```go
di := standard.New().
    MustProvide(22).
    MustProvide(NewUserPoint)


di.MustInvoke(func(data int) { /* data == 22 */ })

diCopy := standard.NewWithParent(di).
    MustProvide(NewDevicePoint)

diCopy.MustInvoke(func(data int) { /* data == 22 */ })
err := diCopy.Invoke(func(device *Device) {
    // ...
})
// err == sameErrorDevice if input int is > 5


diCopyOfCopy := standard.NewWithParent(diCopy).
    MustProvide(3)

diCopyOfCopy.MustInvoke(func(data int) { /* data == 3 */ })
diCopyOfCopy.MustInvoke(func(device *Device) {
    // ...
})
```
Or you can use the integrated override functionality:
```go
standard.New().
    MustProvide(22).
    MustOverride(15).
    MustInvoke(func(data int) { /* data == 15 */ })
```
