# Tests using Go channels and timeouts

I was experiencing bad behavior with Go's time.After and I was curious what the
deal might be. So I wrote some code that sends a ton of messages through a
channel with no timeout code whatsoever, then with time.After and other timeout
schemes. Here's the output I got:

```
noTimeout: 238ns/message, 2.388013862s elapsed
    262144 sys, 20 mallocs
timeAfter: 1179ns/message, 11.79265836s elapsed, 0 timeouts
    472585648 sys, 30000690 mallocs
ticker:    412ns/message, 4.123238732s elapsed, 0 timeouts
    0 sys, 26 mallocs
timer:     715ns/message, 7.153991546s elapsed, 0 timeouts
    262144 sys, 30 mallocs
context1:  1157ns/message, 11.578018106s elapsed, 0 timeouts
    7342328 sys, 50007630 mallocs
```

Okay, speed first, time.After and context take about five times as long per
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

The context code is slow, but it has a lot more features and likely will become
quite widely used. The problem I have with it at the moment is all those
mallocs, though the memory pressure itself isn't too bad. What I need is a
resettable context; I may work on that next.

Have I completely missed something obvious? A good way to use time.After? Maybe
some other way of doing timeouts? If so, please, please let me know by popping
a GitHub issue on this repository. Thanks!
