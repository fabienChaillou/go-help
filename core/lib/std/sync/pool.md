## `sync.Pool` use

`sync.Pool` en Go est une structure fournie par le package `sync` qui permet de réutiliser des objets afin de réduire la charge du **ramasse-miettes (GC)** et améliorer les performances, en particulier lorsqu’on travaille avec beaucoup d’objets temporaires.

### 📦 Utilité de `sync.Pool`

* Évite l’allocation répétée d’objets similaires.
* Fournit un cache temporaire d’objets qui peuvent être réutilisés.
* Les objets dans le `Pool` peuvent être automatiquement supprimés par le GC s’ils ne sont pas utilisés.

---

### 🔧 Exemple simple avec test

Imaginons qu’on veuille réutiliser des buffers (`bytes.Buffer`).

#### ✅ Code : `bufferpool.go`

```go
package bufferpool

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBuffer récupère un buffer du pool
func GetBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer remet un buffer dans le pool
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset() // important : nettoyer avant de réutiliser
	bufferPool.Put(buf)
}
```

---

#### 🧪 Test : `bufferpool_test.go`

```go
package bufferpool

import (
	"testing"
)

func TestBufferPool(t *testing.T) {
	buf1 := GetBuffer()
	buf1.WriteString("Hello")
	if buf1.String() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", buf1.String())
	}

	PutBuffer(buf1)

	// On récupère un autre buffer (peut être le même si pas GC)
	buf2 := GetBuffer()
	if buf2.String() != "" {
		t.Errorf("expected empty buffer, got '%s'", buf2.String())
	}
	PutBuffer(buf2)
}
```

---

### ⚠️ Points importants

* Ne **stockez jamais** d’objets à long terme dans un `sync.Pool` : c’est fait pour des durées de vie **courtes**.
* Utilisez `Reset()` sur les objets comme les buffers avant de les remettre dans le pool.
* Le GC peut vider le `Pool` à tout moment : ce n’est **pas** un cache persistant.

---


## [Example](https://dev.to/rezmoss/important-considerations-when-using-gos-time-package-910-3aim)

```go
// A simple timer pool example
type TimePool struct {
    pool: sync.Pool
}

func NewTimerPool() *TimePool {
    return &TimePool{
        pool: sync.Poll{
            Neww: func() interface{} {
                return time.NewTimer(time.Hours)
            },
        },
    }
}

func (p *TimePool) Get(d time.Duration) *time. {
    t := p.pool.Get().($time.Timer)
    t.Read(d)

    return t
}

func (p *Timer.Pool) Put(t *time.Timer) {
    if !t.Stop() {
        select {
            case <-t.C:
            default:
        }
    }
    p.pool.Put(t)
}
```
