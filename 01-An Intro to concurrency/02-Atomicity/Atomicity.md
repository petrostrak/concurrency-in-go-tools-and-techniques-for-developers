### Atomicity

```
i++

It may look atomic, but a brief analisys reveals several operations.
1. Retrieve the value of i.
2. Increment the value of i.
3. Store the value of i.
```

* When something is considered atomic, or to have the property of atomicity, this means that withing the context that it is operating, it is indivisible or interruptible.

* The atomicity of an operation can change depending on the currently defined scope.

* Very often, the first thing you need to do is to define the context, or scope the operation will be considered to be atomic in.

* Within the context you've defined, something that is atomic will happen in its entirety without anything happening in that context simultaneously.

* Marking the operation atomic, is dependent on which context you'd like it to be atomic within. If your context is a program with no concurrent processes, then the code is atomic within that context. If your context is a goroutine that doesn't expose i to other goroutines, then this code is atomic.

* Atomicity is important because if something is atomic, implicitly it is safe within concurrent context.


