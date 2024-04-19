# Candystore
Lets find our top customers!

First, download the repo
Run:
```bash
$ go build
```

run the program
```bash
$ ./candystore
```

If you're really frisky pipe it to jq
```bash
$ ./candystore | jq
```

You should now see expected result 
```json
[
  {
    "name": "Jonas",
    "favouriteSnack": "Kexchoklad",
    "totalSnacks": 1982
  },
  {
    "name": "Annika",
    "favouriteSnack": "Center",
    "totalSnacks": 208
  },
  {
    "name": "Jane",
    "favouriteSnack": "Nötchoklad",
    "totalSnacks": 22
  },
  {
    "name": "Aadya",
    "favouriteSnack": "Center",
    "totalSnacks": 11
  }
]
```
