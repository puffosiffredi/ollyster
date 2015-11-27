
Ollyster AIMS to be a clone of Friendi.ca, GNU Social , Diaspora, written in go.

It is based on the original work of  Gou(合)  Gou is a Japanese BBS software, used to create BBS and federate them using a P2P scheme.
Gou is Pretty popular in Japan, even is not that known in the western world.

You can find the original Gou software [here](https://github.com/shingetsu-gou/shingetsu-gou)

Or you can download executable binaries from [here](https://github.com/shingetsu-gou/shingetsu-gou/releases).


# What is now Ollyster and how it differs from Gou?

1. In Gou all users are anonymous. Ollyster will add user identity. 
2. Gou is not a Social Network, so there is no User Profile. Ollyster will add one.
3. Gou has no personal stream, Ollyster will add one as a local "BBS"
4. Gou represents data as a list of BBS, Ollyster will represent it as a personal, social stream.

# How Ollyster is different from GNU Social, Friendi.ca &co

1. No dependencies. No Mysql, no apache, no curl,  no libraries, just go, golang and an executable.
2. No PHP: get rid of ~400 security issues included with PHP.
3. No "Ostatus". This protocol sucks, end of story. 
4. Just unzip on disk , compile and run. Compile is what we do in go: "go get; go build" . That's it.







