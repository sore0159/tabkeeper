### 2017 09 04

CSS has proved to be the bulk of this project design, as usability is the primary focus of the project.  Besides the normal "making the page look good" tasks, there was particular difficulty configuring css to highlight checkbox/radio button labels based on the status of their boxes.

I'm happy with html5's form validations for the Description/amount fields.  Though since the step for amount is $0.01 I'm disappointed it is not simple to hide the field's step-wheel.  I've set up defaults so that for the most part, using the page should require minimal typing.

Instead of working on 'saved' charges that would require an entire extra subsystem, I am allowing charges to be marked as 'repeatable'; this allows the UI for repeating charges to be built into the listing of all charges in a natural manner.  I default to non-repeatable simply to minimize display clutter.

Color design will require user feedback: I think having different colors for each direction of charge will be helpful, but it could be too visually noisy.

Next step is to get it running on my local system and configure apache to redirect properly.  This may require some reconfiguration of the app's routing, but otherwise the project is ready for use!


### 2017 09 03

It turns out the specific needs of my app were not suited to the more abstract concurrency guards I had built around my files.  Safe read and safe write methods are nice, but my only two interactions are really "read all" or "add this item to whatever you already have".  So I overhauled the concurrency guards to just specifically handle those two operations and abandoned the file-level wrappers.

With regards to security, my first pass of cut&paste IP checking failed, because when I worked on it before I never ran from localhost.  It's not really a great solution anyway so I shelved it for now, but I need to figure out something before this is done.

Also, I realized logging is a serious concern.  There are several places where things could go wrong with file read/writes, and as this program is meant to run in the background on my server hardware I need to setup a logfile for any potential errors.  Nothing too complicated, but another gear that needs to be fitted into the machine.  For now I'm using the golang STD log package, works fine.

Starting to tinker around with the HTML interface for the app, which means working with go's template package.  Another area where I've heard there are some really nice fancy alternatives, but checking out new tools is the opposite of getting something done.

I'm going to have to start thinking about how to implement "saved" commands: if I want a group of monthly reoccurring changes to be available with a single button click, but also configurable through the app, that means another separate subsystem with its own UI and persistence needs.


### 2017 09 02

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

