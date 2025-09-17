# Blacklight
Version 2 of my old package manager. It is now a module and can be used programatically. It's also pure python, so no extra requirements are needed :)

- Keep in mind it will manipulate the current interpreter and will hang it when running apps, I recommend running a subprocess for now. 

## how to use
Installation
`pip install git+https://github.com/frijoles524/Blacklight.git#subdirectory=blacklight`
- main.py is an example script that uses blacklight to install and run an app that is defined in lists/raven.json. 
- venv_handler is the backend that handles things such as installing modules, loading environments and making sure pip is accessible.