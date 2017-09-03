2017 09 02

First design issue: concurrent file access control.  

Pass 1 was a wrapper struct around an \*os.File that used a sync.RWMutex to control Read and Write calls.  This failed when I realized that after reading from the file, the next pageload would not start over at the beginning.  I looked at adjusting the ReadAt mark or using Seek to move the cursor, but this got complicated when considering that I did _not_ want to adjust Write calls.

Questions this raised:

* Is it a bad idea for a long-running server to hold an open file for an extended time?
* How does the system handle multiple simultaneous "file" objects on the same filesystem object?

Pass 2 was a separate manager in charge of opening/closing file objects for every request, guarded again my a sync.RWMutex.  I miss rust's guard syntax for this kind of problem.  Pass 2A looked like:

```go
func (m *SafeManager) GetR() (*os.File, err) { 
     m.Guard.RLock()
     return os.Open(m.FileName)
}
func (m *SafeManager) DoneR(f *os.File) {
     f.Close()
     m.RUnlock()
}

func Example() {
     f, err := m.GetR()
     if err != nil { ... }
     defer m.DoneR(f)
     ...
}
```

But I ended up with my current Pass 2B system, of having a wrapper around \*os.File _and_ a manager: calls to the manager wait on the mutex, then open and wrap a \*os.File in a SafeFile, which wraps the Close() method to also unlock the appropriate guard in the manager.  This means when you use a SafeFile you don't need to worry about how it's safety is implemented beyond normal File closing.

Not sure what to do if there is an error on the file close call: I just ignore these for now. 

So, the current system can safely read/write to a file based on web calls.  Next up, creating a data structure to represent tab entries, and implementing JSON parsing for read/writes.  I have some concern about processing needed once the list gets long, but at present I believe having a yearly file rotation will keep things fine.

