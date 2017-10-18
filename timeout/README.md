# Kill another goroutine with endless loop

open one channel in goroutine to receive "end" signal, and use `return`, instead of `break` for this nested `select` statement.

```go
			case <-quit:
				log.Print("time is up, quit.")
				return
```