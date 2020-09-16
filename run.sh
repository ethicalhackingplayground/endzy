amass enum -passive -d $1 | httpx -silent | getJS -complete | tee -a domains
cat domains | ./endzy -domains domains | teleman
