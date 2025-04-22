package lottus

type InMemorySessionStorage struct {
	Sessions map[string]Session
}

func (ss InMemorySessionStorage) Get(id string) (error, Session) {
	return nil, ss.Sessions[id]
}

func (ss InMemorySessionStorage) Add(s Session) error {
	ss.Sessions[s.Msisdn] = s
	return nil
}

func (ss InMemorySessionStorage) Delete(s Session) error {
	return nil
}

func (ss InMemorySessionStorage) Update(s Session) error {
	ss.Sessions[s.Msisdn] = s
	return nil
}
