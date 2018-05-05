# GoRoutineTest
with dispatcher:
round 1 round 2 round 3 round 4 round 5 round 6 round 7 round 8 round 9 round 10 
without dispatcher:
round 1 round 2 round 3 round 4 round 5 round 6 round 7 round 8 round 9 round 10 

Averages:

1 go routines: dispatcher=822.1µs, no dispatcher=64.5048ms

10 go routines: dispatcher=905.7µs, no dispatcher=65.2275ms

100 go routines: dispatcher=1.2671ms, no dispatcher=57.5885ms

1000 go routines: dispatcher=5.1713ms, no dispatcher=68.2719ms

10000 go routines: dispatcher=27.949ms, no dispatcher=329.4079ms

100000 go routines: dispatcher=461.3224ms, no dispatcher=965.2058ms

1000000 go routines: dispatcher=2.3125777s, no dispatcher=7.7980054s
