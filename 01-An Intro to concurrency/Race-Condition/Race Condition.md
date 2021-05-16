### Race Condition

A race condition occurs when two or more operations must execute in the correct order, but the programm has not been written so that this order is guaranteed to be maintained.

Most of the time, this shows up in hat's called a data race; when one concurrent operation attempts to read a variable while at some undetemined time another concurrent operation is attempting to write to the same variable.