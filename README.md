# Tests using Go channels and timeouts

I was experiencing bad behavior with Go's time.After and I was curious what the
deal might be. So I wrote some code that sends a ton of messages through a
channel with no timeout code whatsoever, then with time.After, and finally with
a _fuzzy_ time.Ticker. Here's the output I got:

```
noTimeout: 239ns/message, 2.394148382s elapsed
    0 sys, 21 mallocs
timeAfter: 1166ns/message, 11.66043857s elapsed, 0 timeouts
    464394448 sys, 30000574 mallocs
ticker:    404ns/message, 4.046922356s elapsed, 0 timeouts
    0 sys, 12 mallocs
timer:     714ns/message, 7.14432028s elapsed, 0 timeouts
    0 sys, 12 mallocs
```

Okay, speed first, time.After takes about five times as long per message
whereas the time.Ticker takes about twice and time.Timer about three times as
long when comparing them to the no timeout code.

But honestly, it's the secondary part that's even more troubling. time.After
causes about half a gigabyte of memory allocations over thirty million mallocs!
That's insane. I understand it's because we're creating a new *thing* for each
message; it's just that maybe there should be a note in the documentation
discouraging heavy use of this. "May cause destructive use of memory and thrash
the garbage collector; only use in extreme moderation."

The time.Ticker code isn't going to give you a perfect one second timeout;
it'll be somewhere between just over a second to just under two seconds. But it
could be modified to have a better resolution if needed.

I think the time.Timer solution is the "best". It has more time overhead per
message than time.Ticker, but it's more precise and, I think, easier to read.

Have I completely missed something obvious? A good way to use time.After? Maybe
some other way of doing timeouts? If so, please, please let me know by popping
a GitHub issue on this repository. Thanks!
