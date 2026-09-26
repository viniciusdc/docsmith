# CLI reference

Generated from the live `--help` output of the binary by
`internal/tools/gendocs`. Run `make docs` to refresh it; edits between the
markers are overwritten.

<!-- gendocs:cli -->
hello is the sample command-line tool in the docsmith template. Replace it
with your own; the generators in internal/tools/ only need a cobra binary.

Global flags, accepted by every command: `-u/--upper`

| Command | Description |
|---------|-------------|
| [`count`](#hello-count) | Count lines, words and bytes, like wc |
| [`greet`](#hello-greet) | Print a greeting |

---

### `hello count`

Count lines, words and bytes, like wc

```sh
hello count <file|->...
```

```sh
hello count notes.txt
cat notes.txt | hello count -
```

---

### `hello greet`

Print a greeting

```sh
hello greet [name] [-g/--greeting GREETING]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-g/--greeting GREETING` | `Hello` | word to greet with |

```sh
hello greet
hello greet Ada --greeting Hi
```
<!-- /gendocs:cli -->
