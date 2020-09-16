subfinder -d $1 -all -silent | httpx -silent | getJS -complete | tee -a domains
cat domains | ./endzy -domains domains | slackcat --channel $2 --token $3 --stream
