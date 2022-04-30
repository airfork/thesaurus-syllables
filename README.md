# Thesaurus Syllables

A simple endpoint to search a word and get related words as well the number of syllables each word has.

## Running

Install go 1.18, it will work on lower versions just change the value in [go.mod](go.mod).

Run `go run .` to run without building.

Run `go build .` to build an executable

## Searching

To get synonyms, visit `localhost:8080/api/?search=fancy`. This will use [Datamuse's](http://www.datamuse.com/)
"Means like search". Adding `rel` as a query param, `localhost:8080/api/?search=fancy&rel`, will use
the related word search which can provide different results. 

Sample response (truncated)

```json
[
    {
        "word": "chic",
        "score": 72095,
        "numSyllables": 1
    },
    {
        "word": "cool",
        "score": 68919,
        "numSyllables": 1
    },
    {
        "word": "cute",
        "score": 73124,
        "numSyllables": 1
    },
    {
        "word": "fad",
        "score": 70658,
        "numSyllables": 1
    }
]
```

## Changing Values

Visit [main.go](main.go) and modify the `PORT` and `BasePath` variables to modify the port and base path of the
endpoint respectively 
