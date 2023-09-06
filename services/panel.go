package services

import (
	uuid "github.com/gofrs/uuid/v5"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

type PanelService struct{}

var _psPrimary _psToken
var _psSecondary _psToken
var _psLocker sync.RWMutex

type _psToken struct {
	Token    string
	ExpireAt time.Time
}

func (ps *PanelService) IsValid(token string) bool {
	_psLocker.RLock()
	defer _psLocker.RUnlock()

	if _psPrimary.Token == token {
		return _psPrimary.ExpireAt.After(time.Now())
	}
	if _psSecondary.Token == token {
		return _psSecondary.ExpireAt.After(time.Now())
	}

	return false
}

func (ps *PanelService) GetActiveToken() string {
	_psLocker.RLock()
	defer _psLocker.RUnlock()
	return _psPrimary.Token
}

func init() {
	_psIssueNewToken()
	go func() {
		//cycle our tokens every minute, this reduces attack vectors and replay attacks
		//while we could generate them per-request, this makes it easier
		ticker := time.NewTicker(time.Minute * 1)
		for {
			<-ticker.C
			_psIssueNewToken()
		}
	}()
}

func _psIssueNewToken() {
	_psLocker.Lock()
	defer _psLocker.Unlock()

	//move old key to secondary
	_psSecondary = _psPrimary

	//generate new primary
	_psPrimary = _psToken{
		Token:    _psGenerateToken(),
		ExpireAt: time.Now().Add(time.Second * 90),
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func _psGenerateToken() string {
	id, err := uuid.NewV4()
	if err != nil {
		//we could not generate an UUID, so we have to generate one
		//based on https://stackoverflow.com/a/31832326
		//we want speed because we call this... a lot
		n := 32
		src := rand.NewSource(time.Now().UnixNano())
		b := make([]byte, n)
		// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
		for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}

		return *(*string)(unsafe.Pointer(&b))
	}
	return id.String()
}
