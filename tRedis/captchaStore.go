package tRedis


type captchaStore struct {

}

var captchaStoreInstance captchaStore
func GetInstanceByCaptchaStore () *captchaStore{
	return &captchaStoreInstance
}

// Set sets the digits for the captchas id.
func (store *captchaStore) Set(id string, value string){
	SetSingle(id,value,300)
}

// Get returns stored digits for the captchas id. Clear indicates
// whether the captchas must be deleted from the store.
func (store *captchaStore) Get(id string, clear bool) string {
	v := GetString(id)
	if clear {
		DelKey(id)
	}
	return v
}

//	Verify 检验
//	id string
//	answer string
//	clear bool
func (store *captchaStore) Verify(id, answer string, clear bool) bool{
	v := GetString(id)
	if "" == v {
		return false
	}
	flag := v == answer
	if clear {
		DelKey(id)
	}
	return flag
}
