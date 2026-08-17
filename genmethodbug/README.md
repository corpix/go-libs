# go1.27rc3 compiler deadlock on valid generic-method code (fatal error: all goroutines are asleep - deadlock!)

## go version

    go version go1.27rc3 linux/amd64

## What triggers it

`lib/lib.go` declares a method with its own type parameters -- Go 1.27
added generic methods as a language feature, so this is valid code:

    func (p *impl[T]) M[V any]() Alias[V] {
        return nil
    }

where `Alias[T any] = *impl[T]` is a generic type alias.

(Confirmed this is intentional, not an accident: the standard library
itself uses the same shape at this exact pinned go1.27rc3 tag, see
`src/math/rand/v2/rand.go:213`:
`func (r *Rand) N[Int intType](n Int) Int`.)

1. `go build ./lib/` (type-checking straight from source) compiles this
   with no error, correctly, since the declaration is valid Go 1.27.

2. Building any package that *imports* `lib` (forcing the compiler to read
   `lib`'s compiled export data back in via the unified-IR importer instead
   of type-checking from source) deadlocks the compiler process instead of
   compiling successfully:

       # genmethodbug
       fatal error: all goroutines are asleep - deadlock!

       goroutine 1 [sync.Mutex.Lock]:
       ...
       cmd/compile/internal/types2.(*Named).unpack(...)
           cmd/compile/internal/types2/named.go:226
       ...
       cmd/compile/internal/importer.(*reader).signature(...)
           cmd/compile/internal/importer/ureader.go:352
       cmd/compile/internal/importer.(*pkgReader).objIdx.func1.1(...)
           cmd/compile/internal/importer/ureader.go:483
       cmd/compile/internal/types2.(*Named).unpack(...)
           cmd/compile/internal/types2/named.go:266
       ...
       (repeats)

   Goroutine 1 re-enters `(*Named).unpack` for the same `Named` type while
   it is already mid-unpack (the mutex at named.go:226 is not re-entrant),
   so the second `Lock()` call blocks forever. Since it's the only
   goroutine, the runtime reports a deadlock instead of the compiler
   crashing more directly.

   The recursion is driven by resolving the generic alias `Alias[V]` used
   as the method's return type: `newAliasInstance` -> `subst` ->
   `subster.typ` -> `Named.TypeParams` -> `Named.unpack`, which in this
   case re-enters the unified-IR reader (`objIdx.func1.1` ->
   `signature` -> `param` -> `typ` -> `Instantiate` -> `instance` ->
   `newAliasInstance` -> `subst` -> ...) for the very same `Named` type
   that is still being unpacked.

## Reproduce

    cd genmethodbug
    go build ./lib/       # succeeds, correctly -- this is valid Go 1.27
    go build ./...        # deadlocks the compiler -- this is the bug

## Minimal ingredients (confirmed by trimming)

- A generic type alias to a pointer-to-generic-struct:
  `type Alias[T any] = *impl[T]`.
- A generic method (Go 1.27 feature) on that struct, whose result type is
  the generic alias instantiated with the method's own type parameter:
  `func (p *impl[T]) M[V any]() Alias[V]`.
- A second package that imports the first and references the alias
  (e.g. `var _ lib.Alias[int]`), forcing export-data re-import rather than
  source type-checking.

Removing any of these (turning `M` into a plain function with no receiver,
dropping the type param on `M`, using `*impl[V]` directly instead of the
alias, or only blank-importing `lib` without referencing `Alias`) makes the
deadlock go away.

## Expected behavior

Both `go build ./lib/` and `go build ./...` should succeed -- the code is
valid Go 1.27. Instead, the second one deadlocks the compiler as soon as
another package needs to re-import the generic method's signature from
export data.

## Root cause / suggested fix

Non-reentrant lock around a callback that can reenter it, in
`src/cmd/compile/internal/types2/named.go`, `(*Named).unpack()`:

    n.mu.Lock()
    defer n.mu.Unlock()
    ...
    tparams, underlying, methods, delayed := n.loader(n)   // ~line 262

`n.loader(n)` is the unified-IR importer's callback
(`cmd/compile/internal/importer/ureader.go`), and while `n.mu` is still
held it can recurse back into `n.TypeParams()` -> `unpack()` for the same
`n` (exactly what the crash's stack trace shows) -- `sync.Mutex` is not
reentrant, so the second `Lock()` blocks forever. There is already a TODO
directly above the lock acknowledging the risk: "loader... would need to
support reentrant calls though" -- reentrancy was anticipated but never
actually guarded against. Generic methods returning an instantiation of a
generic alias appear to be the first case that actually exercises this
reentrant path through the unified-IR importer.

Suggested fix: track "currently unpacking" state on `Named` and either
short-circuit a reentrant call (return the not-yet-fully-populated `Named`
instead of re-locking), or fail with a clean internal-compiler error
("internal error: cycle unpacking type during import") instead of a silent
mutex hang. There is a related TODO just above the lock noting that
locking may be avoidable entirely when `n.check != nil`, since
type-checking is not concurrent per-`Checker` -- a cheap recursion-guard
flag fits that same observation and would turn any future instance of
this class of bug into a diagnosable crash instead of an opaque
"goroutines asleep" runtime message.
