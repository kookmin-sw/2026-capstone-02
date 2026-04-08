# Imp

Imp acts as the input language of traceinspector. It's a simple imperative language that supports the following:

- integer, boolean, and dynamic array values
- if-else statements
- while loops
- function calls (Imp is pass-by-value, but arrays are passed by reference)

Imp also supports the following builtin functions:
- `Scanf(fmt_string string, ...vars) -> None`: Writes values read from stdin specified by `fmt_string` into locations `vars`. `fmt_string` is `"%t"` or `"%d"`, and `vars` should be assignable expressions (variable names or array indexes).
- `Print(...vals) -> None`: Prints to stdout variadic arguments `vals`. Note that newline is not automatically added.
- `make_array(size int, default val) -> array[val_ty]`: Returns an array of length `size` with values set as `default`
- `len(arr array[var]) -> int`: Returns the length of `arr`

## go2imp

- You can just call `Print`, or use `fmt.Print`
- The same applies for `Scanf`. **You do not need to add the address operator(`&`) for variables**
- Do not use `var` declarations. Instead use short assignments `:=` and typed array compositelits
- To declare an array, use `var := []int{...}`