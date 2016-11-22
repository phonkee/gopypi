# Metadata

Metadata can inspect types and return JSON response for OPTIONS requests.
All api is chainable, but warning it changes data inplace.


Example:

```go
    type Product struct {
        Name string `json:"name"`
    }

    metadata := New("Product detail").Description("Basic information about product")
    metadata.Action(ACTION_CREATE).From(ProductNew{})
    metadata.Action(ACTION_DELETE)
    metadata.Action(ACTION_RETRIEVE).Field("result").From(Product{})
```

This will yield to this when json marshalled

```json
    {
        "name": "test endpoint",
        "description": "description",
        "actions": {
            "POST": {
                "type": "struct",
                "fields": {
                    "name": {
                        "type": "string",
                    }
                }
            }
        }
    }
```

You can even describe more complicated structures

```go
    type User struct {
        Username string `json:"username"`
    }

    md := New("Product detail").Description("Basic information about product")
    md.Action(ACTION_RETRIEVE).Field("result", "user").From(Product{})

```

Which then accepts structure like this:

```json
{
  "result": {
    "user": {
      "username": "phonkee"
    }
  }
}
```

Fields also support choices so you can do this

```go
    md := New()
    md.Action(ACTION_CREATE).Field("status").Choices().Add(1, "new").Add(2, "active").Add(3, "closed")
```

## TODO:
add support for custom validators such as min max value and other.


## Contribute:
Your contribution is welcome, feel free to send PR.