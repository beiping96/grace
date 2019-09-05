package grace

// Option define how goroutine works
type Option interface{ wrap(g Goroutine) Goroutine }

// OptionRestart declare restart option
// -1 means always restart
// func OptionRestart(times int) Option { return &optionRestart{times} }

// type optionRestart struct{ times int }

// func (o *optionRestart) wrap(g Goroutine) Goroutine {
// 	times := o.times
// 	return func(ctx context.Context) {
// 		for {

// 			times--
// 		}
// 	}
// }
