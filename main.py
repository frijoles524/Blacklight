from blacklight import load_apps, install_software, run_software
from blacklight.venv_handler.piputil import install_pip
from blacklight.venv_handler.installer import install_dependency_global

# Normal startup sequence. install_pip only runs when pip is not detected
store = load_apps("lists")
install_pip()
# Requiring global dependencies is also possible. 
# some modules like QScintilla install pyqt5 or other requirements, which inflates the size of each app by upwards of 100mb
install_dependency_global("pyqt5")
install_dependency_global("QScintilla")

app_name = "scratchpad"
latest_version = store.get_latest_version(app_name)
app = store.get_app(app_name, latest_version)

install_software(app)

run_software(app.name, app.version, store)
