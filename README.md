# quartz

Quartz enables you to call Go code from within your
Ruby code. This is accomplished by running a Go process
alongside your Ruby application in a thread and using
UNIX domain sockets for communicating between the two
processes.

## In Go

Import the Go package by adding this to your import path:

```go
import (
	"github.com/DavidHuie/quartz/go/quartz"
)
```

Quartz maintains most of Go's native RPC interface. Thus, export a struct
by calling this to your `main` function:

```go
myStruct := &myStructType{}
quartz.RegisterName("my_struct", myStruct)
```

Once all structs have been declared, start the RPC server by running this:

```go
quartz.Start()
```

## Defining exportable structs

Quartz requires that all arguments to exported struct methods be serializable
structs. Additionally, the arguments to a method should be (A, *B) where A and
B are any struct types. The method should also return an error.
Here's an example of an exportable struct and method:

```go
type Adder struct{}

type Args struct {
	A int
	B int
}

type Response struct {
	Sum int
}

func (t *Adder) Add(args Args, response *Response) error {
	*response = Response{args.A + args.B}
	return nil
}
```

## In Ruby

After you've found created a `go run`-able file, create a Go process wrapper:

```ruby
require 'quartz'

go_process = Quartz::GoProcess.new(file_path: 'spec/test.go')
```

Now you should create a client:

```ruby
client = Quartz::Client.new(go_process)
```

To call a method on a struct:

```ruby
client[:adder].call('Add', 'A' => 2, 'B' => 5)
# => { 'Sum' => 7 }
```

## Copyright

Copyright (c) 2014 David Huie. See LICENSE.txt for further details.
