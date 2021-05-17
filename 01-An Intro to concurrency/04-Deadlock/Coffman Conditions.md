## Coffman Conditions
There are a few conditions that must be present for deadlocks to arise:

### Mutual Exclusion
A concurrent process holds exclusive rights to a resource at any time.

### Wait for Condition
A concurrent process must simultaneously hold a resource and be waitiong for an additional resource.

### No Preemption
A resource held by a concurrent process can only be released by that process.

### Circular Wait
A concurrent process (P1) must be waiting on a chain of other concurrent processes (P2), which are in turn waiting on it (P1).