# httpetag
--
    import "vimagination.zapto.org/httpetag"

Package httpetag wraps a http.FileStstem to provide ETag headers

## Usage

#### func  New

```go
func New(root http.FileSystem) http.Handler
```
New creates a new Handler around an http.FileSystem which generates sha256
hash's of the file contents as ETags

#### func  NewWithHandler

```go
func NewWithHandler(root htto.FileSystem, handler http.Handler) http.Handler
```
NewWithHandler creates a new handler, as in New, but with a custom handler which
is passed to after setting the ETag header
