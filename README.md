# go-strongparams
RubyOnRails inspired Golang implementation of Strong Parameters

## Strong Params
Strong parameters is a really great approach of whitelisting model properties. [`github.com/gorilla/schema`](github.com/gorilla/schema) is a really
great tool. But, it doesn't support dynamic whitelisting on the struct. Through struct tags it is only possible to
define which fields are required.

Let's take a look on the following use case:
```go
type OptionalParams struct {
    OptionalKey1 *string `params:"key1"`
    OptionalKey2 *string `params:"key2"`
}

optionalParams := OptionalParams{}
```
It is pretty complicated to retrieve the one or the other value from the request query string or post parameters.

This can easily be solved with `go-strongparams`. :smile:
```go
optionalParams := OptionalParams{}
values := url.Values{
    "key1": []string{"value"},
    "key2": []string{"ignoredValue"},
}

Params().Permit("key1").Values(values)(&optionalParams)

output, _ := json.Marshal(values)
fmt.Println(output) // > {"key1": "value", "key2": null}
```

### (*StrongParams) Require(requireKey string) *StrongParams
`Require` enables defining a single key which defines an entity that needs to be present.
```go
queryRequest := // ?entity[key1]=value1&entity[key2]=value2
Params().Require("entity").Query(queryRequest)(&optionalParams)

// OR

postFormRequest := // entity[key1]=value1
                   // entity[key2]=value2
Params().Require("entity").PostForm(postFormRequest)(&optionalParams)

// OR

values := url.Values{
    "entity[key1]": []string{"value1"},
    "entity[key2]": []string{"value2"},
}
Params().Permit("entity").Values(values)(&optionalParams)
```

### (*StrongParams) RequireOne(requireKey string) *StrongParamsRequireOne
`RequireOne` enables validating that a key is present and retrieving it.
```go
queryRequest := // ?entity[key1]=1&entity[key2]=2
i, err := Params().RequireOne("entity[key1]").Query(queryRequest)(strconv.Atoi)

// OR

postFormRequest := // entity[key1]=value1
                   // entity[key2]=value2
i, err := Params().RequireOne("entity[key1]").PostForm(postFormRequest)(strconv.Atoi)

// OR

values := url.Values{
    "entity[key1]": []string{"value1"},
    "entity[key2]": []string{"value2"},
}
i, err := Params().RequireOne("entity[key1]").Values(values)(strconv.Atoi)
```

### (*StrongParams) Permit(permitRule string, permitRules... string) *StrongParamsRequiredAndPermitted
`Permit` enables whitelisting keys.
```go
type EntityHolder struct {
    Entity string `params:"entity"`
}
entity := EntityHolder{}

queryRequest := // ?entity[key1]=value1&entity[field2]=value2
Params().Permit("entity:{key1, key2}").Query(queryRequest)(&entity)

// OR

postFormRequest := // entity[key1]=value1
                   // entity[key2]=value2
Params().Permit("entity:{key1, key2}").PostForm(postFormRequest)(&entity)

// OR

values := url.Values{
    "entity[key1]": []string{"value1"},
    "entity[key2]": []string{"value2"},
}
Params().Permit("entity:{key1, key2}").Values(values)(&entity)
```

### (*StrongParamsRequired) Permit(permitRule string, permitRules... string) *StrongParamsRequiredAndPermitted
`Permit` enables whitelisting keys.
```go
queryRequest := // ?entity[key1]=value1&entity[field2]=value2
Params().Require("entity").Permit("key1, key2").Query(queryRequest)(&optionalParams)

// OR

postFormRequest := // entity[key1]=value1
                   // entity[key2]=value2
Params().Require("entity").Permit("key1, key2").PostForm(postFormRequest)(&optionalParams)

// OR

values := url.Values{
    "entity[key1]": []string{"value1"},
    "entity[key2]": []string{"value2"},
}
Params().Require("entity").Permit("key1, key2").Values(values)(&optionalParams)
```

### [`github.com/gorilla/schema`](github.com/gorilla/schema) dot notation
[`github.com/gorilla/schema`](github.com/gorilla/schema) uses a dot notation (eg. `entity.0.key`) instead of brackets notation (eg.
`entity[0][key]`). `go-strongparams` helps to overcome this downside. Before passing the `url.Values` to the
`schema.Decoder` it transforms the bracket notation keys to dot notation keys.
```go
queryRequest := // ?entity[key1]=value1&entity[field2]=value2
Params().Query(queryRequest)(&entity)

// OR

postFormRequest := // entity[key1]=value1
                   // entity[key2]=value2
Params().PostForm(postFormRequest)(&entity)

// OR

values := url.Values{
    "entity[key1]": []string{"value1"},
    "entity[key2]": []string{"value2"},
}
Params().Values(values)(&entity)
```

More examples in [./test/StrongParams_test.go](./test/StrongParams_test.go)

### `schema.Decoder`
By default `go-strongparams` uses the following `schema.Decoder` configuration.
```go
decoder := schema.NewDecoder()
decoder.SetAliasTag("params")
```
Tag alias **`params`** is used instead of `schema.Decoder` default **`schema`**.

To override the default decoder two methods can be used.
- Overrides given StrongParams struct's decoder:
  ```go
  Params().WithDecoder(newDecoder)
  ```
- Every new StrongParams struct will have the following decoder:
  ```go
  WithDecoder(newDecoder).Params()
  ```
  This way there is no need to define your explicit decoder for every `StrongParam` separately.

## Permitter
The `github.com/vellotis/go-strongparams/permitter` package holds a simple rule engine that is capable of validating if
a specific key is permitted or not.

The rules can have the following elements:
- Key (KeyLiteral) defines a whitelisted key:
  - string literal of ASCII numbers and letters, eg. `someKey`
  - single quote `'` wrapped string literal of ASCII numbers, letters and spaces, eg. `'some key'`

- Object (ObjectLiteral) defines a whitelisted object containing any nested literals:
  - `{ KeyLiteral, KeyLiteral }` eg.<br/>
    `{ key1, key2 }` matches query `key1=value1&key2=value2`

  - `{ KeyLiteral:Array, KeyLiteral:ObjectLiteral, KeyLiteral }` eg.<br/>
    `{ key1:[], key2:{objKey}, key3}` matches query `key1[]=value&key2[objKey]=objValue&key3=keyValue`

  **NOTE!** Rule `{}` doesn't match anything


- Array (ArrayLiteral) defines a whitelisted array containing any nested literals:
  - `[]` matches array of string values, eg.<br/>
    `[]=value1&[]value2` or `[0]=value1&[1]value2`
  - `[ KeyLiteral ]` eg.<br/>
    `[ key ]` matches query `[0][key]=value`
  - `[ KeyLiteral:Array, KeyLiteral:ObjectLiteral, KeyLiteral ]` eg.<br/>
    `{ key1:[], key2:{objKey}, key3 }` matches query
    `[0][key1][]=value1&[0][key1][]=value2&[0][key2][objKey]=objValue&[0][key3]=keyValue`

Some examples:

```go
ParsePermitted(/*Rule*/).IsPermitted(/*Permitted key*/)
```

| Rule                                | Permitted keys                                                        |
|-------------------------------------|-----------------------------------------------------------------------|
| `key`                               | `key`                                                                 |
| `obj:{key}`                         | `obj[key]`                                                            |
| `key1,key2,key3`                    | `key1`<br>`key2`<br>`key3`                                            |
| `key1:[sub1,sub2,sub3],key2:{sub4}` | `key1[0][sub1]`<br>`key1[0][sub2]`<br>`key1[0][sub3]`<br>`key2[sub4]` |
| `key1:{'sub x1':{sub x2:[key]}}`    | `key[sub x1][sub x2][0][key]`                                         |
| `key:[]`                            | `key[]`<br>`key[0]`<br>`key[1]`<br>`key[n]` (n = Number >= 0)         |
| `key:{}`                            | Matches noting                                                        |

More examples in [./permittable/test/Permittable_test.go](./permittable/test/Permittable_test.go)

## License

BSD licensed. See the LICENSE file for details.