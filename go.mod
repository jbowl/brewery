module github.com/jbowl/brewery

go 1.15

require (
	github.com/jbowl/apibrewery v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.33.2
)

//replace github.com/jbowl/apibrewery => /home/j/jsoft/github.com/jbowl/apibrewery
replace github.com/jbowl/apibrewery => ./apibrewery
