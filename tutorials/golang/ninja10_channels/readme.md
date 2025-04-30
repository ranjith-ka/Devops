# Channels

Concurrency

Only one goroutine has access to the value at any given time. Data races cannot occur, by design.
Do not communicate by sharing memory; instead, share memory by communicating.

Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values into another goroutine.

```go
ch := make(chan int)
```

By default, sends and receives block until the other side is ready. This allows goroutines to synchronize without explicit locks or condition variables.

Here's an example of using a channel to sum numbers in a goroutine. When the goroutine completes, it sends a value on the channel to indicate completion.

```go
func sum(s []int, c chan int) {
    sum := 0
    for _, v := range s {
        sum += v
    }
    c <- sum // send sum to c
}
```

You can receive values from a channel like this.

```go
func main() {
    s := []int{7, 2, 8, -9, 4, 0}
    c := make(chan int)
    go sum(s[:len(s)/2], c)
    go sum(s[len(s)/2:], c)
    x, y := <-c, <-c // receive from c
    fmt.Println(x, y, x+y)
    // Output: 16 -1 15
}
```
