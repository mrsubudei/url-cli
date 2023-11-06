## Getting started
Download durl by using:
```sh
go install github.com/mrsubudei/url-cli/cmd/durl@latest
```

If utility is not available after downloading, copy bin file to /usr/local/bin directory
 ```sh
sudo cp /home/$USER/go/bin/durl /usr/local/bin 
```

## Usage
```
durl -f="{file name}" [OPTION]...

	-f		input file name (required argument)
	-s		keep the original order of urls
	-o 		output file name (available only with -s flag)
	-t		custom timeout for http requests in sec (should be >= 1)

Examples: 
durl -f="data.txt"			            read data from file data.txt and print out the answer
durl -f="data.txt" -s -o="output.txt"	save the answer to output.txt
durl -f="data.txt" -t=10		        set request timeout to 10 sec
```  