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

- Builin Functions

```javascript

giggle>> len("Hello World")
11
giggle>> len([1,2,3,4,5])
5
giggle>> tail([1,2,3,4,5])
[2, 3, 4, 5]
giggle>> head([1,2,3,4,5])
1
giggle>> append([1,2,3,4,5],6)
[1, 2, 3, 4, 5, 6]

```

- Types

```javascript

giggle>> "Hello World"
Hello World
giggle>> 2555
2555
giggle>> true
true
giggle>> {"one":1,"two":2,"three":3}
{one: 1, two: 2, three: 3}
giggle>> [1,2,3,4,5]
[1, 2, 3, 4, 5]
giggle>> let map = {"one":1,"two":2,"three":3}
giggle>> map["one"]
1
```

- Functions

```javascript

giggle>> let map = fn(arr,f){
    let iter = fn(arr,acc){
         if (len(arr) == 0){
              acc
        } else {
            iter(tail(arr),append(acc,f(head(arr))));
        }
    };
    iter(arr,[]);
    };
giggle>> map([1,2,3],square)
[1, 4, 9]
giggle>> let cube = fn(x){ x*x*x}
giggle>> cube(3)
27
giggle>> map([1,2,3],cube)
[1, 8, 27]

```

- REPL

```javascript
Make me giggle !
giggle>> let add_mod = fn(x,y,z){ (x + y) % z};
giggle>> add_mod(15,16,3)
1
giggle>> add_mod(252,2343,13)
8
giggle>> exit
Ohhh you are leaving already !⏎

````

Return statements are not needed the language is expression oriented.

The tests contain further code examples you can run.
