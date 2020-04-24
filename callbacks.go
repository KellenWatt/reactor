package reactor

// ReadCallback <- func(interface{})
// WriteCallback <- func(interface{}, interface{})
// BindingFunc <- func(interface{}) interface{}


func (r ReadCallback) Async() ReadCallback {
	return func(v interface{}) {
		go r(v)
	}
}

func (r ReadCallback) Concurrent() ReadCallback {
	conReadLock.Lock()
		if conRead == nil {
			conRead = make(chan conReadState, 100)
			go runConcurrentRead()
		}
	conReadLock.Unlock()
	return func(v interface{}) {
		conRead <- conReadState{v, r}
	}
}

func (r ReadCallback) Conditional(f func(interface{}) bool) ReadCallback {
	return func(v interface{}) {
		if f(v) {
			r(v)
		}
	}
}


func (w WriteCallback) Async() WriteCallback {
	return func(prev, v interface{}) {
		go w(prev, v)
	}
}

func (w WriteCallback) Concurrent() WriteCallback {
	conWriteLock.Lock()
		if conWrite == nil {
			conWrite = make(chan conWriteState, 100)
			go runConcurrentWrite()
		}
	conWriteLock.Unlock()
	return func(prev, v interface{}) {
		conWrite <- conWriteState{prev, v, w}
	}
}

func (w WriteCallback) Conditional(f func(interface{},interface{}) bool) WriteCallback {
	return func(prev, v interface{}) {
		if f(prev, v) {
			w(prev, v)
		}
	}
}


// var AssignmentBinding = func(v interface{}) interface{} {return v}
// 
// func (f BindingFunc) Concurrent(b Binder) BindingFunc {
//     conBindLock.Lock()
//         if conBind == nil {
//             conBind = make(chan conBindState, 100)
//             runConcurrentBind()
//         }
//     conBindLock.Unlock()
//     return func(v interface{}) interface{} {
//         conBind <- conBindState{v, b, f}
//     }
// }
// 
// func (f BindingFunc)
// 
// 
