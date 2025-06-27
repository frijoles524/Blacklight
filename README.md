# Blacklight
This project aims to facilitate management, downloading, running and optimizing of Ravendevteam's free software. I'm aware of the plans to make their own software manager and optimizations, but as it's still in development I've decided to make this project.
# How to use
- First, download the executable from releases or clone this repository and build from source (cgo must be enabled)
- Now, place the executable in your desired location
- Add the executable to path
- Run the program as administrator once to ensure correct permissions while updating and fetching runtime (not really nessesary, but can prevent an error)
- You can now use blacklight. Run blacklight.exe in a terminal for details
# More information
- Why the need for a runtime? Blacklight comes with a lightweight python runtime, required to run raven software from source code. By default their software is compiled with nuitka, which really slows down the startup time and sets off antiviruses because it has to write to the disk. Blacklight can run raven software more efficiently through the use of this runtime.
- Upcoming features. These include a more native experience for running apps (with the help of start menu items and shortcuts), more functionality and other cool features.