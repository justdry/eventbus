# Example

This example demonstrates how to use the [`github.com/justdry/eventbus`](https://github.com/justdry/eventbus) package for simple event-driven programming in Go.
It shows how to define events, subscribe handlers, and emit events with JSON payloads.

---

## Running the Example

The program expects **two arguments**:

1. **event name**
2. **JSON data**

Usage:

```bash
go run . <event-name> '<json-data>'
```

---

## Available Events

| Event Name | Expected JSON Structure                    | Description                                                  |
| ---------- | ------------------------------------------ | ------------------------------------------------------------ |
| `greeting` | `{"firstName": "John", "lastName": "Doe"}` | Greets a person with a message.                              |
| `weather`  | `{"weather": "sunny"}`                     | Responds with a weather-related message.                     |
| `error`    | _any data_                                 | Triggers an intentional error to demonstrate error handling. |

### Command Example

```bash
go run . greeting '{"firstName":"John","lastName":"Doe"}'
```

**Output:**

```
2025/10/19 14:21:05 Hi John Doe,
2025/10/19 14:21:05 How are you doing?
```
