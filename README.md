# bindiff
--
    import "github.com/BenLubar/bindiff"

Package bindiff provides a bidirectional binary patch for pairs of []byte.

## Usage

```go
var ErrCorrupt = errors.New("bindiff: corrupt patch")
```
ErrCorrupt is the only possible error from functions in this package.

#### func  Diff

```go
func Diff(old, new []byte, granularity int) (patch []byte)
```
Diff computes the difference between old and new. A granularity of 1 or more
combines changes with no greater than that many bytes between them.

#### func  Forward

```go
func Forward(old, patch []byte) (new []byte, err error)
```
Forward retrieves the second argument to Diff given the first argument and its
output.

#### func  Reverse

```go
func Reverse(new, patch []byte) (old []byte, err error)
```
Reverse retrieves the first argument to Diff given the second argument and its
output.
