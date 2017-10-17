
Application Name: Salik Crawler 
Exe name: salik.exe
Author: Jiten Palaparthi
Version: v1.8
Last Modified: 21th March 2017

How to execute the file?

The application is developed on darwin 64 bit compiler and cross compiled for windows 64 bit architecture, hence executable file is available.

Execution steps

Copy the executable file in desired location on your pc.
if you are willing to execute it from any location add the path in system environment variable else browse to that path..
open run-cmd(administrator)- go to the path
type **salik.exe -u “username” -p “password” -sd “d-mm-yyyy” -ed “d-mm-yyyy” -a accountno -fn “filename” ** press enter..
-sd is start date
-ed is end date
-fn is file name for windows “c:/folder/folder/sample ” can be given for Mac or linux “/Users/Mac/Desktop/sample” like that to be given.
If -fn and file name is not given then out put is always output file

What is output structure?

The below is the sample output.json file structure

{
	"result": [{
		"TransactionId": "1232313",
		"TripDateTime": "21-Jan-2017 10:50:47 PM",
		"TransactionPostDate": "22-Jan-2017",
		"TollGateLocation": "Al Garhoud New Bridge",
		"TollGateDirection": "Abu Dhabi",
		"Amount": "4.00",
		"PlateSource": "Dubai",
		"PlateColor": "R",
		"PlateNumber": "323132",
		"TagNumber": 123232
	}, {
		"TransactionId": "32322",
		"TripDateTime": "21-Jan-2017 09:08:10 PM",
		"TransactionPostDate": "21-Jan-2017",
		"TollGateLocation": "Airport Tunnel",
		"TollGateDirection": "Sharjah",
		"Amount": "4.00",
		"PlateSource": "Dubai",
		"PlateColor": "R",
		"PlateNumber": "1231",
		"TagNumber": 34242342
	}]
}

Can this exe be extended?

Yes, since salik.exe is only an executable file, it can be executed in any program. At this point of time it is on demand execution , if required it can be extended for periodic execution from any other program. Make sure that the saluki server has robotics.txt file which has some guidelines for number of requests per second. Crawling is always based on the robotics.txt file guidelines.

What are changes in new Version?

1.v1.1 Chnages:User can  write the output in his desired file name with new -fn "filename" command line argument
2.v1.3:Changes:Zero salik records are also added to the system.
3.v1.4 Changes:log.txt file has been added upon failure server page.Usually the error is from the server side.
3.v1.5 Changes:Added simple watch dog implementation
4.v1.6 Changes:Added Excel export option and made few commond line changes
5.v1.7 Changes:Performance improvements.
6.v1.8 Changes:The traial version of CookieJas has been removed and reimplemented.

What bugs are resolved?

1.Few records are not fetch. The issue is identified as problem with total number of records devided by 10(which is a page counnt). Resolved by adding one more page count.
2.Even though total count is shown by the server, the actucal number of records it fetches is identified as different.When the result is null it stops fetching and write the previous data to the file.
3.Server pages are tried to fetch if not brought data then skipped to the next page.Usually
4.In version1.6 there are no bugs to resolve only changes in the requirement


Contact details
for any assistance do not hesitate to contact me 
@mobile: +91-9618558500
@email: jitenp@outlook.com
@skype:jpalaparthi