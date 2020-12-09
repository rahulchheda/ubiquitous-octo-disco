Task: Keep the Calender sync in all the devices.

Solution: Each device (localhost:12345-12350) has its own local Calender. If any of the devices gets a post request on its respective endpoint, it tries to ping the other devices endpoints to add that particular day's calender into it.


TODO:
One problem still left to solve: Find the IP from where it is getting the ping, and dont ping that client, or it will loop endlessly.