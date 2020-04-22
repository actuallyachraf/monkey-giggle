# monkey-giggle

After cloning the project.

## Build

```sh

user@box:$ go build cmd/

```

## Tour

Monkey is a functional language that support closures, conditionals and the usual
package.

- Bindings and Functions :

```javascript

let a = 5;
let b = 10;
let c = 120;

let add = fn(x,y,y){x + y + z};
```

- Conditionals

```javascript

let a = 10;
let b = -10;

let c = if((a+b) == 0){
    5;
} else {
    10;
}

```

Return statements are not needed the language is expression first.

The tests contain further code examples you can run.
