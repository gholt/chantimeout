# Tests using Go channels and timeouts

I was experiencing bad behavior with Go's time.After and I was curious what the
deal might be. So I wrote some code that sends a ton of messages through a
channel with no timeout code whatsoever, then with time.After and other timeout
schemes. Here's the output I got:

```
noTimeout: 237ns/message, 2.375356888s elapsed
    0 sys, 19 mallocs
timeAfter: 1180ns/message, 11.805005762s elapsed, 0 timeouts
    487367864 sys, 30001846 mallocs
ticker:    405ns/message, 4.051711396s elapsed, 0 timeouts
    0 sys, 27 mallocs
timer:     714ns/message, 7.147521405s elapsed, 0 timeouts
    262144 sys, 27 mallocs
context1:  1156ns/message, 11.567305692s elapsed, 0 timeouts
    6818040 sys, 50007514 mallocs
context2:  833ns/message, 8.339954568s elapsed, 0 timeouts
    0 sys, 30 mallocs
```

Okay, speed first, time.After and context1 take about five times as long per
message as having no timeout code. The time.Ticker code is the fastest timeout
scheme, but also a bit wonky to work with and has limited precision.

But honestly, it's the secondary part that's even more troubling. time.After
causes about half a gigabyte of memory allocations over thirty million mallocs!
That's insane. I understand it's because we're creating a new *thing* for each
message; it's just that maybe there should be a note in the documentation
discouraging heavy use of this. "May cause destructive use of memory and thrash
the garbage collector; only use in extreme moderation."

The time.Ticker code isn't going to give you a perfect one second timeout;
it'll be somewhere between just over a second to just under two seconds. But it
could be modified to have a better resolution if needed.

The context1 code (using golang's context.WithTimeout) is slow, but it has a
lot more features and likely will become quite widely used. The problem I have
with it at the moment is all those mallocs, though the memory pressure itself
isn't too bad. What I need is a reusable context...

So I made a really basic reusable context; just the bare minimums to get things
to "work". That's context2 in the tests. It's slower than a straight
time.Timer, but faster than the default context1 and way better behaved on
mallocs and memory pressure. I pushed it up to github.com/gholt/context and
will probably flesh it out over time.

Have I completely missed something obvious? A good way to use time.After? Maybe
some other way of doing timeouts? If so, please, please let me know by popping
a GitHub issue on this repository. Thanks!
