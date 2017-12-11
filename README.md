# text api

## Spam

Check out the readme in the `smspam` project.

Grabbed this: UCI Irvine Machine Learning Repo data set for spam:
[SMS Spam Collection Data Set](http://archive.ics.uci.edu/ml/datasets/SMS+Spam+Collection#)

go to `smspam` and run `go run script/build_script.go`(there is a readme there too)... also need to change a file path in

```go
func init() {
    //path from the root of text api...
    //environ = ParseEnv("production", "config.toml")
    environ = ParseEnv("production", "smspam/config.toml")
}
```

from `"smspam/config.toml"` to `"config.toml"` in order to get it to work.. until i fix stupid config file path issues.

copied enron data [http://www2.aueb.gr/users/ion/data/enron-spam/](http://www2.aueb.gr/users/ion/data/enron-spam/) to spam/not_spam:

```sh
cp spam/* $mywork/text_api/smspam/build_data/training/spam
cp ham/* $mywork/text_api/smspam/build_data/training/not_spam
```

then i put in in my x/training data. must ignore this data for github