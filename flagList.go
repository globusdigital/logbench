package main

type flagList struct {
	name string
	vals []string
}

func (l *flagList) String() string { return l.name }

func (l *flagList) Set(value string) error {
	l.vals = append(l.vals, value)
	return nil
}

func (l *flagList) RemoveDuplicates() {
	nw := make([]string, 0, len(l.vals))
	reg := make(map[string]struct{})
	for _, loggerName := range l.vals {
		if _, ok := reg[loggerName]; !ok {
			nw = append(nw, loggerName)
			reg[loggerName] = struct{}{}
		}
	}
	l.vals = nw
}
