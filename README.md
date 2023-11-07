## Description
This CLI utility provides asynchronous processing of URLs from a file,   
display its content size and processing time. It can work in two modes.  
Fast and keep sequence mode. To enable pass -s flag  

In Fast mode, program reads data from file,  
and prints the result as soon as the answer is ready.   

In keep sequence mode, the program keeps the original order of urls,   
and provides the result after all  urls handled. Also, in this mode  
it is possible to save results in file, flag -o={file name}.  

By default, http request are done with timeout 5 sec.   
To change it, pass flag -t={number}.  
For details look at [Usage](#usage)

## Getting started
Download durl by using:
```sh
go install github.com/mrsubudei/url-cli/cmd/durl@latest
```

If utility is not available after downloading, copy bin file from folder where  
go installed to /usr/local/bin directory
 ```sh
sudo cp /home/$USER/go/bin/durl /usr/local/bin 
```

## Usage
```
durl -i="{file name}" [OPTION]...

	-i		input file name (required argument)
	-s		keep the original order of urls
	-o 		output file name (available only with -s flag)
	-t		custom timeout for http requests in sec (should be greater than 0)

Examples: 
durl -i="data.txt"						read data from file data.txt and print out the answer
durl -i="data.txt" 	-s -o="output.txt"	save the answer to output.txt
durl -i="data.txt" -t=10				set request timeout to 10 sec
```  