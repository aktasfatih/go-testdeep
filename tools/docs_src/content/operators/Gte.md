---
title: "Gte"
weight: 10
---

```go
func Gte(minExpectedValue interface{}) TestDeep
```

[`Gte`]({{< ref "Gte" >}}) operator checks that data is greater or equal than
*minExpectedValue*. *minExpectedValue* can be any numeric or
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable) value. *minExpectedValue* must be the
same kind as the compared value if numeric, and the same type if
[`time.Time`](https://golang.org/pkg/time/#Time) (or assignable).

[TypeBehind]({{< ref "operators#typebehind-method" >}}) method returns the [`reflect.Type`](https://golang.org/pkg/reflect/#Type) of *minExpectedValue*.


### Examples

{{%expand "Int example" %}}	t := &testing.T{}

	got := 156

	ok := Cmp(t, got, Gte(156), "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte(155), "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte(157), "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
{{%expand "String example" %}}	t := &testing.T{}

	got := "abc"

	ok := Cmp(t, got, Gte("abc"), `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte("abb"), `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = Cmp(t, got, Gte("abd"), `checks "%v" is ≥ "abd"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
## CmpGte shortcut

```go
func CmpGte(t TestingT, got interface{}, minExpectedValue interface{}, args ...interface{}) bool
```

CmpGte is a shortcut for:

```go
Cmp(t, got, Gte(minExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Int example" %}}	t := &testing.T{}

	got := 156

	ok := CmpGte(t, got, 156, "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = CmpGte(t, got, 155, "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = CmpGte(t, got, 157, "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
{{%expand "String example" %}}	t := &testing.T{}

	got := "abc"

	ok := CmpGte(t, got, "abc", `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = CmpGte(t, got, "abb", `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = CmpGte(t, got, "abd", `checks "%v" is ≥ "abd"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
## T.Gte shortcut

```go
func (t *T) Gte(got interface{}, minExpectedValue interface{}, args ...interface{}) bool
```

[`Gte`]({{< ref "Gte" >}}) is a shortcut for:

```go
t.Cmp(got, Gte(minExpectedValue), args...)
```

See above for details.

Returns true if the test is OK, false if it fails.

*args...* are optional and allow to name the test. This name is
used in case of failure to qualify the test. If `len(args) > 1` and
the first item of *args* is a `string` and contains a '%' `rune` then
[`fmt.Fprintf`](https://golang.org/pkg/fmt/#Fprintf) is used to compose the name, else *args* are passed to
[`fmt.Fprint`](https://golang.org/pkg/fmt/#Fprint). Do not forget it is the name of the test, not the
reason of a potential failure.


### Examples

{{%expand "Int example" %}}	t := NewT(&testing.T{})

	got := 156

	ok := t.Gte(got, 156, "checks %v is ≥ 156", got)
	fmt.Println(ok)

	ok = t.Gte(got, 155, "checks %v is ≥ 155", got)
	fmt.Println(ok)

	ok = t.Gte(got, 157, "checks %v is ≥ 157", got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}
{{%expand "String example" %}}	t := NewT(&testing.T{})

	got := "abc"

	ok := t.Gte(got, "abc", `checks "%v" is ≥ "abc"`, got)
	fmt.Println(ok)

	ok = t.Gte(got, "abb", `checks "%v" is ≥ "abb"`, got)
	fmt.Println(ok)

	ok = t.Gte(got, "abd", `checks "%v" is ≥ "abd"`, got)
	fmt.Println(ok)

	// Output:
	// true
	// true
	// false
{{% /expand%}}